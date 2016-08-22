package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 设置玩家任务进度
func (taskmodule *TTaskMoudle) UpdatePlayerTask(taskID int, count int, status int) bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTask", bson.M{"_id": taskmodule.PlayerID, "tasklist.taskid": taskID}, bson.M{"$set": bson.M{
		"tasklist.$.taskcount":  count,
		"tasklist.$.taskstatus": status}})
}

//! 增加玩家成就达成列表
func (taskmodule *TTaskMoudle) AddAchievementCompleteLst(achievementID int) bool {
	return mongodb.AddToArray(appconfig.GameDbName, "PlayerTask", bson.M{"_id": taskmodule.PlayerID}, "achievedlist", achievementID)
}

//! 设置玩家成就进度
func (taskmodule *TTaskMoudle) UpdatePlayerAchievement(taskID int, count int, status int) bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTask", bson.M{"_id": taskmodule.PlayerID, "achievementlist.id": taskID}, bson.M{"$set": bson.M{
		"achievementlist.$.taskcount":  count,
		"achievementlist.$.taskstatus": status}})
}

//! 设置玩家任务积分
func (taskmodule *TTaskMoudle) UpdatePlayerTaskScore(score int) bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTask", bson.M{"_id": taskmodule.PlayerID}, bson.M{"$set": bson.M{
		"taskscore": score}})
}

//! 设置玩家任务积分宝箱领取状态
func (taskmodule *TTaskMoudle) UpdatePlayerTaskScoreAwardStatus() bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTask", bson.M{"_id": taskmodule.PlayerID}, bson.M{"$set": bson.M{
		"scoreawardstatus": taskmodule.ScoreAwardStatus}})
}

//! 日常任务信息存储数据库
func (self *TTaskMoudle) UpdateDailyTaskInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTask", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"taskscore":        self.TaskScore,
		"scoreawardstatus": self.ScoreAwardStatus,
		"scoreawardid":     self.ScoreAwardID,
		"tasklist":         self.TaskList}})
}

//! 更新成就信息
func (self *TTaskMoudle) UpdateAchievement(info *TAchievementInfo, findID int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTask",
		bson.M{"_id": self.PlayerID, "achievementlist.id": findID}, bson.M{"$set": bson.M{
			"achievementlist.$.taskstatus": info.TaskStatus,
			"achievementlist.$.id":         info.ID,
			"achievementlist.$.type":       info.Type}})
}

//! 更新重置时间
func (self *TTaskMoudle) UpdateResetTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTask", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"resetday": self.ResetDay}})
}
