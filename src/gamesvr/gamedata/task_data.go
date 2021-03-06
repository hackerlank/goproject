package gamedata

import (
	"gamelog"
	"strings"
)

const (
	TASK_MAINCOPY_CHALLENGE      = 1  //! 挑战主线副本次数
	TASK_MAINCOPY_STAR           = 2  //! 主线副本星数
	TASK_ELITECOPY_CHALLENGE     = 3  //! 精英副本挑战次数
	TASK_ELITECOPY_STAR          = 4  //! 精英副本星数
	TASK_DAILYCOPY_CHALLENGE     = 5  //! 日常副本挑战次数
	TASK_FAMOUSCOPY_CHALLENGE    = 6  //! 名将副本挑战次数
	TASK_EQUI_STRENGTHEN         = 7  //! 装备强化次数
	TASK_EQUI_REFINED            = 8  //! 装备精炼次数
	TASK_GEM_STRENGTHEN          = 9  //! 宝物强化
	TASK_HERO_CULTURE            = 10 //! 英雄培养
	TASK_ARENA_CHALLENGE         = 11 //! 竞技场挑战
	TASK_USER_LOGIN              = 12 //! 玩家登陆
	TASK_RECHARGE                = 13 //! 玩家充值
	TASK_PASS_MAIN_COPY_CHAPTER  = 14 //! 通过主线副本章节
	TASK_LEVEL_UP                = 15 //! 玩家升级
	TASK_HERO_EQUI_STRENGTH      = 16 //! 上阵所有英雄六件装备强化等级
	TASK_HERO_EQUI_QUALITY       = 17 //! 上阵所有英雄六件装备品质
	TASK_ARENA_RANK              = 18 //! 竞技场排名
	TASK_COMPOSITION             = 19 //! 合成宝物
	TASK_COMPOSITION_PURPLE      = 20 //! 合成紫色宝物
	TASK_COMPOSITION_ORANGE      = 21 //! 合成橙色宝物
	TASK_SGWS_RESET              = 22 //! 三国无双重置次数
	TASK_SGWS_RANK               = 23 //! 三国无双最高排名
	TASK_HERO_EQUI_REFINED       = 24 //! 上阵所有英雄精炼等级
	TASK_HERO_EQUI_REFINED_MAX   = 25 //! 最高精炼等级
	TASK_HERO_DESTINY_LEVEL      = 26 //! 上阵所有英雄天命等级
	TASK_HERO_DESTINY_LEVEL_MAX  = 27 //! 最高天命等级
	TASK_BUY_ZHENGTAOLING        = 28 //! 购买征讨令
	TASK_ATTACK_REBEL_DAMAGE     = 29 //! 攻击叛军伤害
	TASK_REBEL_EXPLOIT           = 30 //! 围剿叛军功勋累积
	TASK_PASS_EPIC_COPY          = 31 //! 通关史诗战役
	TASK_HERO_STORE_REFRESH      = 32 //! 神将商店刷新
	TASK_HERO_STORE_BUY          = 33 //! 神将商店购买商品
	TASK_SGWS_STAR               = 34 //! 三国无双星数
	TASK_HERO_GEM_REFINED        = 35 //! 英雄宝物精炼
	TASK_HERO_GEM_REFINED_MAX    = 36 //! 英雄宝物精炼最高等级
	TASK_FIGHT_VALUE             = 37 //! 战斗力
	TASK_CARD_MASTER_SCORE       = 38 //! 获取卡牌大师积分
	TASK_ROB_TIMES               = 39 //! 今日夺宝次数
	TASK_KILL_REBEL              = 40 //! 今日击杀叛军个数
	TASK_SPENT_MONEY             = 41 //! 消费任意金额
	TASK_COMPLETE_ALL_TASK       = 42 //! 完成所有任务
	TASK_GET_HUNT_SCORE          = 43 //! 获得巡游积分
	TASK_SEND_ACTION             = 44 //! 赠送精力数
	TASK_CAMP_BATTLE_KILL        = 45 //! 阵营战击杀
	TASK_TERRITORY_HUNT          = 46 //! 领地巡逻次数
	TASK_SENIOR_SUMMON           = 47 //! 高级抽将次数
	TASK_SINGLE_RECHARGE         = 48 //! 单笔充值元数
	TASK_AWAKE_STORE_REFRESH     = 49 //! 觉醒商店刷新
	TASK_AWAKE_STORE_BUY         = 50 //! 觉醒商店购买
	TASK_BUY_ACTION_STRENGTH     = 51 //! 购买体力道具
	TASK_BUY_ACTION_ENERGY       = 52 //! 购买精力道具
	TASK_HERO_QUALITY            = 53 //! 玩家品质
	TASK_CUR_HERO_BREAK          = 54 //! 上阵六名英雄等级
	TASK_DIAOWEN_QUALITY         = 55 //! 雕文品质
	TASK_HERO_WAKE               = 56 //! 上阵六名英雄觉醒等级
	TASK_GOD_TIMES               = 57 //! 点亮火炬
	TASK_EQUI_STAR_THREE         = 58 //! N件装备升至3星
	TASK_EQUI_STAR_FIVE          = 59 //! N件装备升至5星
	TASK_FASHION_COMPOSE         = 60 //! N品质时装合成
	TASK_FASHION_STRENGTH        = 61 //! 时装强化
	TASK_FASHION_RECAST          = 62 //! 时装熔炼
	TASK_PASS_ELITE_COPY_CHAPTER = 63 //! 通关精英副本章节
	TASK_GUAJI_HOURS             = 64 //! 挂机小时数
	TASK_CAMP_BATTLE_GROUP_KILL  = 65 //! 阵营战团灭对手次数
	TASK_PASS_MINING             = 66 //! 通关矿区
	TASK_TERRITORY_PATROLTIME    = 67 //! 领地征讨巡逻时间
	TASK_SCORE_RANK              = 68 //! 积分赛段位
	TASK_FIGHT_WINING            = 69 //! 争霸赛外围胜利次数
	TASK_WANDER_SWEEP            = 70 //! 云游扫荡
	TASK_FRIEND_NUM              = 71 //! 好友数量
	TASK_GUILD_LEVEL             = 72 //! 所在公会等级
	TASK_CAMP_HERO_FULL_1        = 73 //! 指定阵营英雄收集
	TASK_CAMP_HERO_FULL_2        = 74 //! 指定阵营英雄收集
	TASK_CAMP_HERO_FULL_3        = 75 //! 指定阵营英雄收集
	TASK_CAMP_HERO_FULL_4        = 76 //! 指定阵营英雄收集
	TASK_PET_QUALITY             = 77 //! 获取宠物品质
	TASK_PET_LEVEL               = 78 //! 宠物升级
	TASK_PET_GOD                 = 79 //! 宠物神练
	TASK_PET_STAR                = 80 //! 宠物升星
)

type ST_TaskInfo struct {
	TaskID       int //! 任务唯一标识
	Type         int //! 任务类型
	Count        int //! 次数
	AwardItem    int //! 奖励物品
	Score        int //! 获得积分
	NeedMinLevel int //! 需求最小等级 当种类为七日活动时,则代表第几天到第几天的活动
	NeedMaxLevel int //! 需求最大等级
}

type ST_AchievementInfo struct {
	TaskID    int //! 任务唯一标识
	Type      int //! 任务类型
	Count     int //! 次数
	AwardItem int //! 奖励物品
	FrontID   int //! 前置成就
	NeedLevel int //! 需求最小等级
}

type ST_TaskSevenActivityInfo struct {
	AwardType   int //! 奖励模板
	TaskID      int //! 任务唯一标识
	TaskType    int //! 任务类型
	Count       int //! 次数
	AwardItem   int //! 奖励物品
	OpenDay     int //! 开放时间
	IsSelectOne int //! 是否为全部发放或玩家选择
}

//! 七日活动-半价限购
type ST_TaskSevenActivityStore struct {
	AwardType int //! 奖励模板
	OpenDay   int
	ItemID    int
	ItemNum   int
	MoneyID   int
	MoneyNum  int
	Limit     int
}

var GT_TaskType_Lst [][]int = nil
var GT_Task_List []ST_TaskInfo = nil
var GT_Achievement_Lst []ST_AchievementInfo = nil
var GT_SevenActivity_Lst []ST_TaskSevenActivityInfo = nil
var GT_SevenActivityStore_Lst []ST_TaskSevenActivityStore = nil

func InitTaskTypeParser(total int) bool {
	GT_TaskType_Lst = make([][]int, total+1)
	return true
}

func ParseTaskTypeRecord(rs *RecordSet) {
	taskType := CheckAtoi(rs.Values[0], 0)

	if len(rs.Values[1]) <= 0 {
		panic("field: 1          Is Empty")
	}

	pv := strings.Split(rs.Values[1], "|")
	for _, v := range pv {
		GT_TaskType_Lst[taskType] = append(GT_TaskType_Lst[taskType], CheckAtoi(v, 1))
	}
}

func InitTaskParser(total int) bool {
	GT_Task_List = make([]ST_TaskInfo, total+1)

	return true
}

func GetTaskSubType(taskType int) []int {
	if taskType > len(GT_TaskType_Lst)-1 {
		gamelog.Error("GetTaskSubType fail. invalid taskType: %d", taskType)
		return []int{}
	}
	return GT_TaskType_Lst[taskType]
}

func ParseTaskRecord(rs *RecordSet) {
	taskID := CheckAtoi(rs.Values[0], 0)
	GT_Task_List[taskID].TaskID = taskID
	GT_Task_List[taskID].Type = rs.GetFieldInt("type")
	GT_Task_List[taskID].Count = rs.GetFieldInt("count")
	GT_Task_List[taskID].AwardItem = rs.GetFieldInt("award")
	GT_Task_List[taskID].Score = rs.GetFieldInt("score")
	GT_Task_List[taskID].NeedMinLevel = rs.GetFieldInt("minlevel")
	GT_Task_List[taskID].NeedMaxLevel = rs.GetFieldInt("maxlevel")
}

func InitAchievementParser(total int) bool {
	GT_Achievement_Lst = make([]ST_AchievementInfo, total+1)

	return true
}

func ParseAchievementRecord(rs *RecordSet) {
	taskID := CheckAtoi(rs.Values[0], 0)
	GT_Achievement_Lst[taskID].TaskID = taskID
	GT_Achievement_Lst[taskID].Type = rs.GetFieldInt("type")
	GT_Achievement_Lst[taskID].Count = rs.GetFieldInt("count")
	GT_Achievement_Lst[taskID].AwardItem = rs.GetFieldInt("award")
	GT_Achievement_Lst[taskID].FrontID = rs.GetFieldInt("front")
	GT_Achievement_Lst[taskID].NeedLevel = rs.GetFieldInt("needlevel")
}

func InitSevenActivityParser(total int) bool {
	GT_SevenActivity_Lst = make([]ST_TaskSevenActivityInfo, total+1)
	return true
}

func ParseSevenActivityRecord(rs *RecordSet) {
	taskID := CheckAtoi(rs.Values[0], 0)
	GT_SevenActivity_Lst[taskID].TaskID = taskID
	GT_SevenActivity_Lst[taskID].AwardType = rs.GetFieldInt("award_type")
	GT_SevenActivity_Lst[taskID].TaskType = rs.GetFieldInt("tasktype")
	GT_SevenActivity_Lst[taskID].Count = rs.GetFieldInt("count")
	GT_SevenActivity_Lst[taskID].AwardItem = rs.GetFieldInt("award")
	GT_SevenActivity_Lst[taskID].OpenDay = rs.GetFieldInt("openday")
	GT_SevenActivity_Lst[taskID].IsSelectOne = rs.GetFieldInt("is_select_one")
}

func InitSevenActivityStoreRecord(total int) bool {
	GT_SevenActivityStore_Lst = make([]ST_TaskSevenActivityStore, total+1)
	return true
}

func ParseSevenActivityStoreRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_SevenActivityStore_Lst[id].OpenDay = id
	GT_SevenActivityStore_Lst[id].ItemID = rs.GetFieldInt("itemid")
	GT_SevenActivityStore_Lst[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_SevenActivityStore_Lst[id].MoneyID = rs.GetFieldInt("moneyid")
	GT_SevenActivityStore_Lst[id].MoneyNum = rs.GetFieldInt("moneynum")
	GT_SevenActivityStore_Lst[id].Limit = rs.GetFieldInt("limit")
}

func GetTaskInfo(taskid int) *ST_TaskInfo {
	if taskid >= len(GT_Task_List) || taskid <= 0 {
		gamelog.Error("GetTaskInfo Error: invalid taskid %d", taskid)
		return nil
	}

	return &GT_Task_List[taskid]
}

func GetAchievementInfo(taskid int) *ST_AchievementInfo {
	if taskid >= len(GT_Achievement_Lst) || taskid <= 0 {
		gamelog.Error("GetAchievementTaskInfo Error: invalid taskid %d", taskid)
		return nil
	}

	return &GT_Achievement_Lst[taskid]
}

func GetSevenTaskInfoFromAwardType(awardType int) []ST_TaskSevenActivityInfo {
	taskLst := []ST_TaskSevenActivityInfo{}
	for i := 0; i < len(GT_SevenActivity_Lst); i++ {
		if GT_SevenActivity_Lst[i].AwardType == awardType {
			taskLst = append(taskLst, GT_SevenActivity_Lst[i])
		}
	}

	return taskLst
}

func GetSevenTaskInfo(taskid int) *ST_TaskSevenActivityInfo {
	if taskid >= len(GT_SevenActivity_Lst) || taskid <= 0 {
		gamelog.Error("GetAchievementTaskInfo Error: invalid taskid %d", taskid)
		return nil
	}

	return &GT_SevenActivity_Lst[taskid]
}

//! 获取七日活动任务
func GetSevenDayTask() []ST_TaskSevenActivityInfo {
	return GT_SevenActivity_Lst
}

//! 获取完成前置的成就任务
func GetNextAchievement(taskID int) *ST_AchievementInfo {
	for i := 0; i < len(GT_Achievement_Lst); i++ {
		if GT_Achievement_Lst[i].FrontID == taskID {
			return &GT_Achievement_Lst[i]
		}
	}

	return nil
}

//! 获取物品信息
func GetSevnActivityItemInfo(openDay int, activityType int) *ST_TaskSevenActivityStore {
	for i := 0; i < len(GT_SevenActivityStore_Lst); i++ {
		if GT_SevenActivityStore_Lst[i].AwardType == activityType && GT_SevenActivityStore_Lst[i].OpenDay == openDay {
			return &GT_SevenActivityStore_Lst[i]
		}
	}
	return nil
}
