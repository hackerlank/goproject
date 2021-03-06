package msg

//! GM请求更新配制文件
//! 消息: /update_gamedata
type MSG_UpdateGameData_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	TbName     string //表名
}

type MSG_UpdateGameData_Ack struct {
	RetCode int //返回码
}

//! GM增发全服奖励
//! 消息: /add_svr_award
type MSG_SvrAward_Add_Req struct {
	SessionID  string         //GM SessionID
	SessionKey string         //GM SessionKey
	Value      string         //! 参数
	ItemLst    []MSG_ItemData //! 奖励内容
}

type MSG_SvrAward_Add_Ack struct {
	RetCode int
}

//! GM删除全服奖励
//! 消息: /del_svr_award
type MSG_SvrAward_Del_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	ID         int
}

type MSG_SvrAward_Del_Ack struct {
	RetCode int
}

//! GM发个人奖励
//! 消息: /send_award_to_player
type MSG_Send_Award_Player_Req struct {
	SessionID  string         //GM SessionID
	SessionKey string         //GM SessionKey
	TargetID   int32          //目标玩家
	Value      string         //参数
	ItemLst    []MSG_ItemData //奖励内容
}

type MSG_Send_Award_Player_Ack struct {
	RetCode int
}

//! 查看当前服务器状态
//! 消息: /get_server_info
type MSG_GetServerInfo_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey

}

type MSG_GetServerInfo_Ack struct {
	RetCode      int
	SvrID        int32  //当前的服务器ID
	SvrName      string //当前服务器名字
	OnlineCnt    int    //在线人数
	MaxOnlineCnt int    //总人数
	RegisterCnt  int    //总注册人数
}

//验证玩家登录请求
//消息:/gm_set_svrstate
type MSG_SetGameSvrState_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	SvrID      int32  //服务器ID
	SvrState   uint32 //服务器标记
	SvrDefault uint32 //是否默认
}

type MSG_SetGameSvrState_Ack struct {
	RetCode int //返回码 0:成功 1: 失
}

//请求服务器列表
//消息:/get_server_list
type MSG_GetServerList_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
}

type MSG_GetServerList_Ack struct {
	RetCode int
	SvrList []ServerNode //服务器结点表
}

//gm用户登录
//消息:/gm_login
type MSG_GmLogin_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
}

type MSG_GmLogin_Ack struct {
	RetCode int
}

//消息:/gm_enable_account
type MSG_GmEnableAccount_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	PlayerID   int32  //角色ID也是账号ID
	SvrID      int32  //分区ID
	RoleName   string //角色名字
	Enable     int32  //0:表示禁用, 1:表示启用
}

type MSG_GmEnableAccount_Ack struct {
	RetCode int
}

//消息:/gm_add_giftaward
type MSG_AddGiftAward_Req struct {
	SessionID  string
	SessionKey string
	ItemID     []int //物品ID
	ItemNum    []int //物品数量
}

type MSG_AddGiftAward_Ack struct {
	RetCode int
	AwardID int
}

//消息:/gm_make_giftcode
type MSG_MakeGiftCode_Req struct {
	SessionID   string  //GM SessionID
	SessionKey  string  //GM SessionKey
	Platform    []int32 //平台ID
	SvrID       []int32 //服务器ID
	EndTime     int32   //结束时间
	GiftAwardID int32   //奖励ID
	GiftCodeNum int     //激活码数量
	IsAll       bool    //是否为全服发放
}

type MSG_MakeGiftCode_Ack struct {
	RetCode   int
	GiftCodes []string //激活码
}

//消息:/gamesvr_giftcode
type MSG_GameSvrGiftCode_Req struct {
	ID        string //礼包ID
	SvrID     int32  //服务器ID
	AccountID int32  //玩家ID
}

type MSG_GameSvrGiftCode_Ack struct {
	RetCode int
	ItemID  []int //物品ID
	ItemNum []int //物品数量
}

//! 查询账号服务器ID
//! 消息: /query_account_id
type MSG_QueryAccountID_Req struct {
	Name string
}

type MSG_QueryAccountID_Ack struct {
	AccountID int32
}

//!	查询玩家信息
//! 消息: /query_account_info
type MSG_QueryAccountInfo_Req struct {
	AccountID int32
}

type MSG_QueryAccountInfo_Ack struct {
	RetCode       int
	AccountName   string //! 账号
	AccountPwd    string //! 密码
	CreateTime    int32  //! 创建时间
	LastLoginTime int32  //! 上次登录时间
	Platform      int32  //! 平台
	Enable        int32  //! 封号状态 0: 表示禁用  1: 表示启用
}

//! 查询玩家信息-GameSvr
//! 消息: /query_player_info
type MSG_QueryPlayerInfo_Req struct {
	PlayerID   int32
	PlayerName string
}

type MSG_QueryPlayerInfo_Ack struct {
	RetCode        int
	PlayerID       int32   //! ID
	PlayerName     string  //! 昵称
	Sex            int     //! 性别
	Phone          string  //! 手机
	Mac            string  //! Mac地址
	Charge         int32   //! 充值额
	ChargeGetMoney int32   //! 充值所获钻石
	ChargeTimes    int32   //! 充值次数
	Level          int     //! 等级
	VIPLevel       int     //! VIP等级
	Money          [14]int //! 货币
	Strength       int     //! 体力
	Action         int     //! 精力
	AttackTimes    int     //! 净化次数
	FightValue     int32   //! 战力
	System         string  //! 手机系统
	LastLogoffTime int32   //! 上次登出时间
	IsOnline       bool    //! 是否在线
	LastLoginIP    string  //! 上次登录IP
}

//! 剔除竞技场排行榜
//! 消息: /kick_cheat_ranker
type MSG_KickCheatRanker_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	PlayerID   int32
}

type MSG_KickCheatRanker_Ack struct {
	RetCode int
}

type MSG_ActivityUpdate struct {
	ID        int32
	Name      string
	Desc      string
	Ad        string
	CycleType int
	TimeType  int
	BeginTime int
	EndTime   int
	AwardTime int
	Type      int
	AwardType int
	Status    int
	Icon      int
	Inside    int
	Days      int
	Change    int //! 修改标记  0->无修改   1->修改   2->添加
}

//! 更新活动表
//! 消息: /update_activity_list
type MSG_UpdateActivityList_Req struct {
	SessioniD   string
	SessionKey  string
	ActivityLst []MSG_ActivityUpdate
}

type MSG_UpdateActivityList_Ack struct {
	RetCode int
}

//! 获取服务器IP黑白名单
//! 消息: /get_net_list
type MSG_GetNetList_Req struct {
	SessionID  string
	SessionKey string
	SvrID      int32
}

type MSG_GetNetList_Ack struct {
	RetCode   int
	WhiteList []string //! 白名单
	BlackList []string //! 黑名单
}

//! 增加服务器IP名单
//! 消息: /add_net_list
type MSG_AddNetList_Req struct {
	SessionID  string
	SessionKey string
	SvrID      int32
	ListType   int    //! 1->白名单 2->黑名单
	IP         string //! IP
}

type MSG_AddNetList_Ack struct {
	RetCode int
}

//! 删除服务器IP名单
//! 消息: /del_net_list
type MSG_DelNetList_Req struct {
	SessionID  string
	SessionKey string
	SvrID      int32
	ListType   int    //! 1->白名单 2->黑名单
	IP         string //! IP
}

type MSG_DelNetList_Ack struct {
	RetCode int
}

//! 查询游戏服ip
//! 消息: /gm_query_svrip
type MSG_QuerySvrIp_Req struct {
	SessionID  string
	SessionKey string
	SvrID      int32
}

type MSG_QuerySvrIp_Ack struct {
	RetCode int
	SvrIp   string
}
