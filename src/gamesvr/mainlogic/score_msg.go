package mainlogic

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"time"
)

func SelectScoreTarget(pPlayer *TPlayer, value int) bool {
	if pPlayer.ScoreMoudle.Score < value {
		//	return false
	}

	return true
}

//请求积分赛目标信息
func Hand_GetScoreTarget(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreTarget_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreTarget Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreTarget_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if gamedata.IsFuncOpen(gamedata.FUNC_SCORE_SYSTEM, player.GetLevel(), player.GetVipLevel()) == false {
		gamelog.Error("Hand_GetScoreTarget Error: Score system func not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.Score = player.ScoreMoudle.Score
	response.Targets, response.Rank = player.ScoreMoudle.GetScoreTargets()
	response.RetCode = msg.RE_SUCCESS
}

func Hand_GetScoreBattleCheck(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreBattleCheck_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreBattleCheck Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreBattleCheck_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if req.TargetIndex < 0 || req.TargetIndex >= 3 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_GetScoreBattleCheck Error Invalid TargetIndex:%d", req.TargetIndex)
		return
	}

	var GetFightTargetReq msg.MSG_GetFightTarget_Req
	GetFightTargetReq.PlayerID = player.ScoreMoudle.ScoreEnemy[req.TargetIndex].PlayerID
	GetFightTargetReq.SvrID = player.ScoreMoudle.ScoreEnemy[req.TargetIndex].SvrID

	if GetFightTargetReq.PlayerID == 0 || GetFightTargetReq.SvrID == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_GetScoreBattleCheck Error Invalid PlayerID:%d, and SvrID:%d", GetFightTargetReq.PlayerID, GetFightTargetReq.SvrID)
		return
	}

	buffer, _ = json.Marshal(GetFightTargetReq)
	http.DefaultClient.Timeout = 3 * time.Second
	httpret, err := http.Post(appconfig.CrossGetFightTarget, "text/HTML", bytes.NewReader(buffer))
	if err != nil || httpret == nil {
		gamelog.Error("Hand_GetScoreBattleCheck failed, err : %s !!!!", err.Error())
		return
	}

	buffer = make([]byte, httpret.ContentLength)
	httpret.Body.Read(buffer)
	httpret.Body.Close()
	var GetFightTargetAck msg.MSG_GetFightTarget_Ack
	err = json.Unmarshal(buffer, &GetFightTargetAck)
	if err != nil {
		gamelog.Error("Hand_GetScoreBattleCheck  Unmarshal fail, Error: %s", err.Error())
		return
	}

	response.PlayerData = GetFightTargetAck.PlayerData
	response.RetCode = GetFightTargetAck.RetCode
}

//玩家提交战斗结果
func Hand_SetScoreBattleResult(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SetScoreBattleResult_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetScoreBattleResult Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SetScoreBattleResult_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//如果打赢了
	if req.WinBattle == 1 {
		player.ScoreMoudle.Score += gamedata.OneTimeFightScore
		player.ScoreMoudle.WinTime += 1
	} else {
		player.ScoreMoudle.Score -= gamedata.OneTimeFightScore
		player.ScoreMoudle.WinTime = 0
	}
	if player.ScoreMoudle.Score < 0 {
		player.ScoreMoudle.Score = 0
	}

	player.ScoreMoudle.FightTime += 1
	player.ScoreMoudle.DB_SaveScoreAndFightTime()

	player.RoleMoudle.CostAction(1, 1)
	response.Targets, response.Rank = player.ScoreMoudle.GetScoreTargets()
	response.RetCode = msg.RE_SUCCESS
	return
}

func (score *TScoreMoudle) GetScoreTargets() ([]msg.MSG_Target, int) {
	var ScoreTargetReq msg.MSG_CrossQueryScoreTarget_Req
	ScoreTargetReq.PlayerID = score.PlayerID
	ScoreTargetReq.Score = score.Score
	ScoreTargetReq.SvrID = GetCurServerID()
	ScoreTargetReq.SvrName = GetCurServerName()
	ScoreTargetReq.HeroID = score.ownplayer.HeroMoudle.CurHeros[0].HeroID
	ScoreTargetReq.FightValue = 0
	ScoreTargetReq.PlayerName = score.ownplayer.RoleMoudle.Name

	b, _ := json.Marshal(ScoreTargetReq)
	http.DefaultClient.Timeout = 3 * time.Second
	httpret, err := http.Post(appconfig.CrossQueryScoreTarget, "text/HTML", bytes.NewReader(b))
	if err != nil || httpret == nil {
		gamelog.Error("GetScoreTargets failed, err : %s !!!!", err.Error())
		return nil, -1
	}

	buffer := make([]byte, httpret.ContentLength)
	var ScoreTargetAck msg.MSG_CrossQueryScoreTarget_Ack
	httpret.Body.Read(buffer)
	httpret.Body.Close()

	err = json.Unmarshal(buffer, &ScoreTargetAck)
	if err != nil {
		gamelog.Error("GetScoreTargets  Unmarshal fail, Error: %s", err.Error())
		return nil, -1
	}

	score.rank = ScoreTargetAck.NewRank
	for i := 0; i < len(ScoreTargetAck.TargetList); i++ {
		score.ScoreEnemy[i].FightValue = ScoreTargetAck.TargetList[i].FightValue
		score.ScoreEnemy[i].HeroID = ScoreTargetAck.TargetList[i].HeroID
		score.ScoreEnemy[i].Level = ScoreTargetAck.TargetList[i].Level
		score.ScoreEnemy[i].Name = ScoreTargetAck.TargetList[i].Name
		score.ScoreEnemy[i].PlayerID = ScoreTargetAck.TargetList[i].PlayerID
		score.ScoreEnemy[i].SvrName = ScoreTargetAck.TargetList[i].SvrName
		score.ScoreEnemy[i].SvrID = ScoreTargetAck.TargetList[i].SvrID
	}

	return ScoreTargetAck.TargetList[0:3], ScoreTargetAck.NewRank
}

//请求积分赛排行榜信息
func Hand_GetScoreRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreRank Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreRank_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.RetCode = msg.RE_SUCCESS
	response.ScoreRankList = []msg.MSG_ScoreRankInfo{}
	response.MyRank = player.ScoreMoudle.rank
	response.MyScore = player.ScoreMoudle.Score

	var ScoreRankReq msg.MSG_CrossQueryScoreRank_Req
	b, _ := json.Marshal(ScoreRankReq)
	http.DefaultClient.Timeout = 3 * time.Second
	httpret, err := http.Post(appconfig.CrossQueryScoreRank, "text/HTML", bytes.NewReader(b))
	if err == nil && httpret != nil {
		buffer = make([]byte, httpret.ContentLength)
		var ScoreRankAck msg.MSG_CrossQueryScoreRank_Ack
		httpret.Body.Read(buffer)
		httpret.Body.Close()
		err = json.Unmarshal(buffer, &ScoreRankAck)
		if err != nil {
			gamelog.Error("Hand_GetScoreRank Query Cross Rank Unmarshal fail, Error: %s", err.Error())
			return
		}
		response.ScoreRankList = ScoreRankAck.ScoreRankList
	} else {
		gamelog.Error("Hand_GetScoreRank failed, err : %s !!!!", err.Error())
	}

	return
}

//请求积分赛战斗次数奖励
func Hand_RcvScoreTimeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_RcvScoreTimeAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RcvScoreTimeAward Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_RcvScoreTimeAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	pTimeAwardInfo := gamedata.GetScoreTimeAward(req.TimeAwardID)
	if pTimeAwardInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_RcvScoreTimeAward Invalid Time Award ID: %d", req.TimeAwardID)
		return
	}
	for _, v := range player.ScoreMoudle.RecvAward {
		if v == req.TimeAwardID {
			response.RetCode = msg.RE_ALREADY_RECEIVED
			gamelog.Error("Hand_RcvScoreTimeAward Already received award: %d", req.TimeAwardID)
			return
		}
	}

	if player.ScoreMoudle.FightTime < pTimeAwardInfo.Times {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_RcvScoreTimeAward Not Enough Time: %d", player.ScoreMoudle.FightTime)
		return
	}

	dropItem := gamedata.GetItemsFromAwardID(pTimeAwardInfo.AwardID)
	for _, v := range dropItem {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	player.BagMoudle.AddAwardItems(dropItem)
	player.ScoreMoudle.RecvAward = append(player.ScoreMoudle.RecvAward, req.TimeAwardID)
	player.ScoreMoudle.DB_UpdateRecvAward()
	response.RetCode = msg.RE_SUCCESS
	return
}

//请求积分赛战斗次数奖励信息
func Hand_GetScoreTimeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreTimeAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreTimeAward Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreTimeAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.Awards = player.ScoreMoudle.RecvAward
	response.FightTime = player.ScoreMoudle.FightTime
	response.RetCode = msg.RE_SUCCESS
	return
}

//! 玩家请求积分商店的状态
//! 消息: /get_score_store_state
func Hand_GetScoreStoreState(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreStoreState_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyScoreFightTime Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreStoreState_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.AwardIndex = append(response.AwardIndex, player.ScoreMoudle.AwardStoreIndex...)

	for _, v := range player.ScoreMoudle.StoreBuyRecord {
		var itemInfo msg.MSG_StoreBuyData
		itemInfo.ID = v.ID
		itemInfo.Times = v.Times
		response.ItemLst = append(response.ItemLst, itemInfo)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买积分商店道具
//! 消息: /buy_score_store_item
func Hand_BuyScoreStoreItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BuyScoreStoreItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyScoreStoreItem Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuyScoreStoreItem_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 获取购买物品信息
	itemInfo := gamedata.GetScoreStoreItem(req.StoreItemID)
	if itemInfo == nil {
		gamelog.Error("Hand_BuyScoreStoreItem Error: GetScoreStoreItem nil ID: %d ", req.StoreItemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断购买等级
	if player.GetLevel() < itemInfo.NeedLevel {
		gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough level")
		response.RetCode = msg.RE_NOT_ENOUGH_LEVEL
		return
	}

	//! 根据类型判断积分
	if itemInfo.Type == 2 && itemInfo.NeedScore > player.ScoreMoudle.Score {
		gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough Score")
		response.RetCode = msg.RE_NOT_ENOUGH_SCORE
		return
	}

	//! 判断货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum) == false {
		gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough money")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 判断道具是否足够
	if itemInfo.CostItemID != 0 {
		if player.BagMoudle.IsItemEnough(itemInfo.CostItemID, itemInfo.CostItemNum*req.BuyNum) == false {
			gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough Item")
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}
	}

	//! 检测购买次数是否足够
	if itemInfo.Type == 1 {
		//! 普通商品
		isExist := false
		for i, v := range player.ScoreMoudle.StoreBuyRecord {
			if v.ID == req.StoreItemID {
				isExist = true
				if v.Times+req.BuyNum > itemInfo.MaxBuyTime {
					gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
					response.RetCode = msg.RE_NOT_ENOUGH_TIMES
					return
				}

				player.ScoreMoudle.StoreBuyRecord[i].Times += req.BuyNum
				go player.ScoreMoudle.DB_UpdateStoreItemBuyTimes(i, player.ScoreMoudle.StoreBuyRecord[i].Times)
			}
		}

		if isExist == false {
			//! 首次购买
			if req.BuyNum > itemInfo.MaxBuyTime {
				gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
				response.RetCode = msg.RE_NOT_ENOUGH_TIMES
				return
			}

			var itemData TStoreBuyData
			itemData.ID = req.StoreItemID
			itemData.Times = req.BuyNum
			player.ScoreMoudle.StoreBuyRecord = append(player.ScoreMoudle.StoreBuyRecord, itemData)
			go player.ScoreMoudle.DB_AddStoreItemBuyInfo(itemData)
		}

		//! 扣除货币
		player.RoleMoudle.CostMoney(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum)

		if itemInfo.CostItemID != 0 {
			player.BagMoudle.RemoveNormalItem(itemInfo.CostItemID, itemInfo.CostItemNum*req.BuyNum)
		}

		//! 发放物品
		player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum*req.BuyNum)

	} else if itemInfo.Type == 2 {
		//! 奖励
		if player.ScoreMoudle.AwardStoreIndex.IsExist(req.StoreItemID) >= 0 {
			gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
			response.RetCode = msg.RE_NOT_ENOUGH_TIMES
			return
		}

		player.RoleMoudle.CostMoney(itemInfo.CostMoneyID, itemInfo.CostMoneyNum)
		player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum)

		player.ScoreMoudle.AwardStoreIndex.Add(req.StoreItemID)
		go player.ScoreMoudle.DB_AddStoreAwardInfo(req.StoreItemID)
	}

	response.RetCode = msg.RE_SUCCESS
}

//购买积分赛战斗次数
func Hand_BuyScoreFightTime(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BuyScoreTime_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyScoreFightTime Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuyScoreTime_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	maxTime := gamedata.GetFuncVipValue(gamedata.FUNC_SCORE_FIGHT_TIME, player.GetVipLevel())
	if player.ScoreMoudle.BuyFightTime >= maxTime {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_BuyScoreFightTime Not Enough Time")
		return
	}

	response.RetCode = msg.RE_SUCCESS

	cost := gamedata.GetFuncTimeCost(gamedata.FUNC_SCORE_FIGHT_TIME, player.ScoreMoudle.FightTime)
	pCopyInfo := gamedata.GetCopyBaseInfo(gamedata.ScoreCopyID)
	if pCopyInfo == nil {
		response.RetCode = msg.RE_INVALID_COPY_ID
		gamelog.Error("Hand_BuyScoreFightTime Invalid CopyID :%d", gamedata.ScoreCopyID)
		return
	}

	player.RoleMoudle.CostMoney(2, cost)
	player.RoleMoudle.AddAction(pCopyInfo.ActionType, 1)
	player.ScoreMoudle.BuyFightTime += 1
	player.ScoreMoudle.DB_SaveBuyFightTime()
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(pCopyInfo.ActionType)

	return
}