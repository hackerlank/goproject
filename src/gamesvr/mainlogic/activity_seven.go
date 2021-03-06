package mainlogic

import (
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 七日活动表结构
type TActivitySevenDay struct {
	ActivityID int32       //! 活动ID
	TaskList   []TTaskInfo //! 任务列表
	BuyLst     IntLst      //! 已购买限购商品列表

	VersionCode    int32            //! 版本号
	ResetCode      int32            //! 迭代号
	activityModule *TActivityModule //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivitySevenDay) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivitySevenDay) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self

	self.TaskList = []TTaskInfo{}
	self.BuyLst = []int{}

	self.VersionCode = vercode
	self.ResetCode = resetcode

	awardType := G_GlobalVariables.GetActivityAwardType(activityID)
	taskLst := gamedata.GetSevenTaskInfoFromAwardType(awardType)
	for _, v := range taskLst {
		var info TTaskInfo
		if v.TaskID == 0 {
			continue
		}
		info.ID = v.TaskID
		info.Status = 0
		info.Count = 0
		info.Type = v.TaskType
		self.TaskList = append(self.TaskList, info)
	}

}

//! 刷新数据
func (self *TActivitySevenDay) Refresh(versionCode int32) {
	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivitySevenDay) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.TaskList = []TTaskInfo{}
	self.BuyLst = []int{}

	self.DB_Reset()
}

func (self *TActivitySevenDay) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivitySevenDay) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivitySevenDay) RedTip() bool {
	return false
}

func (self *TActivitySevenDay) DB_Refresh() {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_Refresh fail. ActivityID: %d", self.ActivityID)
		return
	}

	filedName := fmt.Sprintf("sevenday.%d.versioncode", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.VersionCode}})
}

func (self *TActivitySevenDay) DB_Reset() {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_Reset fail. ActivityID: %d", self.ActivityID)
		return
	}

	filedName1 := fmt.Sprintf("sevenday.%d.tasklist", index)
	filedName2 := fmt.Sprintf("sevenday.%d.buylst", index)
	filedName3 := fmt.Sprintf("sevenday.%d.resetcode", index)
	filedName4 := fmt.Sprintf("sevenday.%d.versioncode", index)

	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName1: self.TaskList,
		filedName2: self.BuyLst,
		filedName3: self.ResetCode,
		filedName4: self.VersionCode}})
}

//! 设置玩家任务进度
func (self *TActivitySevenDay) DB_UpdatePlayerSevenTask(taskID int, count int, status int) {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_UpdatePlayerSevenTask fail: Not find activityID: %d ", self.ActivityID)
		return
	}

	indexTask := -1
	for i, v := range self.TaskList {
		if v.ID == taskID {
			indexTask = i
			break
		}
	}

	if indexTask < 0 {
		gamelog.Error("Sevenday DB_UpdatePlayerSevenTaskStatus fail: Not find activityID: %d  taskID: %d", self.ActivityID, taskID)
		return
	}

	filedName := fmt.Sprintf("sevenday.%d.tasklist.%d.taskstatus", index, indexTask)
	filedName2 := fmt.Sprintf("sevenday.%d.tasklist.%d.taskcount", index, indexTask)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName:  status,
		filedName2: count}})
}

//! 设置玩家七日活动限购购买标记
func (self *TActivitySevenDay) DB_AddPlayerSevenTaskMark(ID int) {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_AddPlayerSevenTaskMark fail")
		return
	}

	filedName1 := fmt.Sprintf("sevenday.%d.buylst", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$push": bson.M{filedName1: ID}})
}
