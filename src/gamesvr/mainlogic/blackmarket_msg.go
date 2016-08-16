package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"time"
)

//! 玩家请求黑市信息
func Hand_GetBlackMarketInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetBlackMarket_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBlackMarketInfo Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetBlackMarket_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.BlackMarketModule.CheckReset()

	if player.GetVipLevel() < gamedata.EnterVipLevel && player.BlackMarketModule.OpenEndTime < time.Now().Unix() {
		gamelog.Error("BlackMarket not open: vipLevel %d  openEndTime: %d", player.GetVipLevel(), player.BlackMarketModule.OpenEndTime)
		response.RetCode = msg.RE_BLACK_MARKET_NOT_OPEN
		return
	}

	for _, v := range player.BlackMarketModule.GoodsLst {
		goodsData := gamedata.GetBlackMarketGoodsInfo(v.ID)

		var goods msg.MSG_BlackMarketGoods
		goods.ID = v.ID
		goods.ItemID = goodsData.ItemID
		goods.ItemNum = goodsData.ItemNum
		goods.IsBuy = v.IsBuy
		goods.CostMoneyID = goodsData.CostMoneyID
		goods.CostMoneyNum = goodsData.CostMoneyNum
		goods.Recommend = goodsData.Recommend
		response.GoodsLst = append(response.GoodsLst, goods)
	}

	response.OpenEndTime = player.BlackMarketModule.OpenEndTime
	response.RefreshTime = player.BlackMarketModule.RefreshTime
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买黑市道具
func Hand_BuyBlackMarketGoods(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyBlackMarket_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_BuyBlackMarketGoods Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_BuyBlackMarket_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//player.BlackMarketModule.CheckReset()

	if player.GetVipLevel() <= gamedata.EnterVipLevel && player.BlackMarketModule.OpenEndTime < time.Now().Unix() {
		response.RetCode = msg.RE_BLACK_MARKET_NOT_OPEN
		return
	}

	//! 检测商品是否能够购买
	isExist := false
	for _, v := range player.BlackMarketModule.GoodsLst {
		if v.ID == req.ID {
			isExist = true

			if v.IsBuy == true {
				response.RetCode = msg.RE_ALEADY_BUY
				return
			}
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测商品金钱是否足够
	goodsData := gamedata.GetBlackMarketGoodsInfo(req.ID)
	if player.RoleMoudle.CheckMoneyEnough(goodsData.CostMoneyID, goodsData.CostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(goodsData.CostMoneyID, goodsData.CostMoneyNum)

	//! 给予物品
	player.BagMoudle.AddAwardItem(goodsData.ItemID, goodsData.ItemNum)

	//! 设置状态
	for i, v := range player.BlackMarketModule.GoodsLst {
		if v.ID == req.ID {
			player.BlackMarketModule.GoodsLst[i].IsBuy = true
		}
	}

	go player.BlackMarketModule.DB_UpdateBuyMark(req.ID)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家查询黑市状态
func Hand_GetBlackMarketStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetBlackMarketStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBlackMarketStatus Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetBlackMarketStatus_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.OpenEndTime = player.BlackMarketModule.OpenEndTime
	response.RetCode = msg.RE_SUCCESS
}