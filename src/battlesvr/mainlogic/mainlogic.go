package mainlogic

import (
	"battlesvr/gamedata"
	"gamelog"
	"msg"
	"runtime"
	"runtime/debug"
	"tcpserver"
)

func Init() bool {
	//初始化连接管理器
	InitConMgr()

	//初始化房间管理器
	InitRoomMgr()

	//加载配制文件
	gamedata.LoadConfig()

	return true
}

func BatSvrMsgDispatcher(pTcpConn *tcpserver.TCPConn, msgid int16, extra int16, pdata []byte) {
	if pTcpConn == nil {
		gamelog.Error("BatSvrMsgDispatcher Error: pTcpConn == nil")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				gamelog.Error("MsgID %d Error  %s", msgid, debug.Stack())
			}
		}
	}()

	switch msgid {
	case msg.MSG_ENTER_ROOM_REQ:
		G_RoomMgr.Hand_EnterRoom(pTcpConn, pdata)
	case msg.MSG_LOAD_CAMPBAT_ACK:
		G_RoomMgr.Hand_LoadCampBatAck(pTcpConn, pdata)
	case msg.MSG_CHECK_IN_REQ:
		G_RoomMgr.Hand_CheckInReq(pTcpConn, pdata)
	case msg.MSG_DISCONNECT:
		G_RoomMgr.Hand_Disconnect(pTcpConn, pdata)
	case msg.MSG_HEART_BEAT:
		G_RoomMgr.Hand_HeartBeat(pTcpConn, pdata)
	default:
		{
			var msgno int32
			var pRoom *TBattleRoom = nil
			if pTcpConn.ConnID < 10000 {
				pRoom = G_RoomMgr.GetRoomByID(extra)
			} else {
				msgno = int32(pdata[3])<<24 | int32(pdata[2])<<16 | int32(pdata[1])<<8 | int32(pdata[0])
				if msgno != pTcpConn.PackNo {
					gamelog.Error("BatSvrMsgDispatcher Error: PackNo Not Right msgno:%d, packNo:%d,msgid:%d", msgno, pTcpConn.PackNo, msgid)
					return
				}

				pTcpConn.PackNo = msgno + 1
				pRoom = G_RoomMgr.GetRoomByID(int16(pTcpConn.Extra))
			}

			if pRoom == nil {
				gamelog.Error("BatSvrMsgDispatcher Error: Invalid RoomID:%d, PlayerID:%d", pTcpConn.Extra, pTcpConn.ConnID)
				return
			}

			var msg TMessage
			msg.MsgID = msgid
			msg.MsgData = pdata
			pRoom.MsgList <- msg
		}

	}

}
