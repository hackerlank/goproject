package main

import (
	// "appconfig"
	"gamelog"
	"gamesvr/tcpclient"
	"msg"
	"utility"
)

var (
	accountip string = "http://121.40.207.78:8081"
)

func main() {
	gamelog.InitLogger("httptest", 10)
	RegTcpMsgHandler()

	InitPlayerMgr()
	for i := 1; i < 200; i++ {
		CreatePlayer(i)
	}

	StartTest()
	utility.StartConsoleWait()
}

func RegTcpMsgHandler() {
	tcpclient.HandleFunc(msg.MSG_DISCONNECT, Hand_DisConnect)
	tcpclient.HandleFunc(msg.MSG_CONNECT, Hand_Connect)
	tcpclient.HandleFunc(msg.MSG_ENTER_ROOM_ACK, Hand_EnterRoomAck)
	tcpclient.HandleFunc(msg.MSG_ENTER_ROOM_NTY, Hand_NoneFunction)
	tcpclient.HandleFunc(msg.MSG_MOVE_STATE, Hand_NoneFunction)

}

func Hand_NoneFunction(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
}

func Hand_Connect(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_Connect")

	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		gamelog.Info("Hand_Connect Error: pClient == nil")
		return
	}

	player := pClient.ExtraData.(*TPlayer)
	if player == nil {
		gamelog.Info("Hand_Connect Error: player == nil")
		return
	}

	if pClient.ConType == tcpclient.CON_TYPE_BATSVR {
		player.userEnterRoom()
	} else {
		player.userCheckIn()
	}

	return
}

func Hand_DisConnect(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_DisConnect")
	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	return
}

func Hand_EnterRoomAck(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_EnterRoomAck")
	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	player := pClient.ExtraData.(*TPlayer)
	if player == nil {
		gamelog.Info("Hand_EnterRoomAck Error: player == nil")
		return
	}

	var req msg.MSG_EnterRoom_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_EnterRoomAck : Message Reader Error!!!!")
		return
	}

	player.Heros = req.Heros

	player.IsEnter = true

	player.PackNo = req.BeginMsgNo

	return
}
