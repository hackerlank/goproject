package tcpserver

import (
	"bufio"
	"encoding/binary"
	"gamelog"
	"io"
	"msg"
	"net"
	"time"
)

type TCPConn struct {
	conn       net.Conn
	reader     *bufio.Reader //包装conn减少conn.Read的io次数
	writeChan  chan []byte
	closeFlag  bool //此标记仅在ResetConn、Close中设置，其它地方只读
	Cleaned    bool
	OnNetClose func()

	//应用级的数据
	ConnID int32 //一般是用户ID或者服务器ID为int32,这个是必须的。
	Extra  int32 //用户ID之外的数据，在战场服作roomid，在聊天聊作guildid。
	PackNo int32 //包序号
}

func newTCPConn(conn net.Conn, pendingWriteNum int) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.ResetConn(conn)
	tcpConn.writeChan = make(chan []byte, pendingWriteNum)
	return tcpConn
}

//closeFlag仅在ResetConn、Close中设置，其它地方只读
func (tcpConn *TCPConn) ResetConn(conn net.Conn) {
	tcpConn.conn = conn
	tcpConn.reader = bufio.NewReader(conn)
	tcpConn.closeFlag = false
}
func (tcpConn *TCPConn) Close() {
	if tcpConn.closeFlag {
		return
	}
	tcpConn.conn.Close()
	tcpConn.doWrite(nil) //触发writeRoutine结束
	tcpConn.closeFlag = true
}

func (tcpConn *TCPConn) doWrite(b []byte) {
	if tcpConn.closeFlag {
		gamelog.Error("TCPConn.doWrite Error: tcpConn.closeFlag:%v", tcpConn.closeFlag)
		return
	}

	select {
	case tcpConn.writeChan <- b: //chan满后再写即阻塞，select进入default分支报错
	default:
		gamelog.Error("doWrite: channel full")
		tcpConn.conn.(*net.TCPConn).SetLinger(0)
		tcpConn.conn.Close()
	}
}

// b must not be modified by other goroutines
func (tcpConn *TCPConn) write(b []byte) {
	if tcpConn.closeFlag || b == nil {
		gamelog.Error("TCPConn.Write Error: b == nil or closeFlag:%v", tcpConn.closeFlag)
		return
	}

	tcpConn.doWrite(b)
}

func (tcpConn *TCPConn) WriteMsg(msgID int16, extra int16, msgdata []byte) bool {
	if tcpConn.closeFlag {
		return false
	}

	msgLen := len(msgdata)
	msgbuffer := make([]byte, 8, 8+msgLen)
	binary.LittleEndian.PutUint32(msgbuffer, uint32(msgLen))
	binary.LittleEndian.PutUint16(msgbuffer[4:], uint16(extra))
	binary.LittleEndian.PutUint16(msgbuffer[6:], uint16(msgID))
	msgbuffer = append(msgbuffer[:8], msgdata...)
	tcpConn.write(msgbuffer)
	return true
}

func (tcpConn *TCPConn) WriteMsgData(msgdata []byte) bool {
	tcpConn.write(msgdata)
	return true
}

func (tcpConn *TCPConn) ReadProcess() error {
	var err error
	var msgHeader = make([]byte, 8)
	var msgLen int32
	var msgID int16
	var Extra int16
	var firstTime bool = true

	//循环结束，会在ReadRoutine中紧接着关闭tcpConn
	for {
		if tcpConn.closeFlag {
			break
		}

		if firstTime == true {
			tcpConn.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			firstTime = false
		} else {
			tcpConn.conn.SetReadDeadline(time.Time{})
		}

		_, err = io.ReadAtLeast(tcpConn.reader, msgHeader, 8)
		if err != nil {
			gamelog.Error("ReadProcess error: Read Header Error : Disconnect from client")
			return err
		}

		// parse len
		msgLen = int32(binary.LittleEndian.Uint16(msgHeader[:4]))
		Extra = int16(binary.LittleEndian.Uint16(msgHeader[4:6]))
		msgID = int16(binary.LittleEndian.Uint16(msgHeader[6:]))
		if msgLen <= 0 || msgLen > 10240 {
			gamelog.Error("ReadProcess error: Invalid msgLen :%d", msgLen)
			break
		}

		if msgID <= msg.MSG_BEGIN || msgID >= msg.MSG_END {
			gamelog.Error("ReadProcess error: Invalid msgID :%d", msgID)
			break
		}

		// data
		msgData := make([]byte, msgLen)
		_, err = io.ReadAtLeast(tcpConn.reader, msgData, int(msgLen))
		if err != nil {
			gamelog.Error("ReadProcess error: Read Data Error :%s", err.Error())
			return err
		}

		MsgDispatcher(tcpConn, msgID, Extra, msgData)
	}

	return nil
}

//连接的写协程
func (tcpConn *TCPConn) WriteRoutine() {
	for b := range tcpConn.writeChan {
		if b == nil {
			break
		}

		if tcpConn.closeFlag {
			break
		}

		_, err := tcpConn.conn.Write(b)
		if err != nil {
			gamelog.Error("WriteRoutine error: %s", err.Error())
			break
		}
	}
	tcpConn.conn.Close()
}

//连接的读协程
func (tcpConn *TCPConn) ReadRoutine() {
	tcpConn.ReadProcess()
	tcpConn.Close()
	tcpConn.OnNetClose()
	//通知业务层net断开
	MsgDispatcher(tcpConn, msg.MSG_DISCONNECT, 0, nil)
}
