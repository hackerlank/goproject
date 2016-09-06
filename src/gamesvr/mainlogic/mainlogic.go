package mainlogic

func Init() bool {
	G_Players = make(map[int32]*TPlayer, 1)
	G_SelectPlayers = make([]*TPlayer, 0, 10000)

	//初始化数据库处理器
	InitDbProcesser()

	G_SimpleMgr.Init()

	//! 初始化竞技场排行榜数据
	InitArenaMgr()

	//! 初始化开服基金购买人数
	InitBuyOpenFundNum()

	//初始化工会系统
	InitGuildMgr()

	//初始化排行榜系统
	InitRankMgr()

	//初始化全局变量
	G_GlobalVariables.Init()

	//预加载角色
	PreLoadPlayers()

	//初始化定时器管理器
	G_Timer.Init()

	return true

}
