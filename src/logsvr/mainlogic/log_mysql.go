package mainlogic

import (
	"database/sql"
	"gamelog"
	"msg"
	_ "mysql"
)

type TMysqlLog struct {
	db       *sql.DB
	tx       *sql.Tx
	stmt     *sql.Stmt
	query    string
	writeCnt int
	flushCnt int
	datasrc  string
	svrid    int32
}

func (self *TMysqlLog) CreateLogTable(svrid int32) bool {
	sql := `CREATE TABLE if not exists gamelog(
			eventid int not null,
			srcid int not null,
			chnlid int not null,
			svrid int not null,
			playerid int not null,
			level int not null,
			viplvl int not null,
			time int not null,
			param1 int not null,
			param2 int not null);`
	_, err := self.db.Exec(sql)
	if err != nil {
		gamelog.Error("CreateLogTable Error : %s", err.Error())
		return false
	}
	return true
}

func (self *TMysqlLog) Start(filename string, svrid int32) bool {
	self.datasrc = filename
	self.svrid = svrid
	self.db = nil
	self.tx = nil

	var err error = nil
	self.db, err = sql.Open("mysql", self.datasrc)
	if err != nil {
		gamelog.Error("TMysqlLog Open file Error : %s", err.Error())
		return false
	}

	err = self.db.Ping()
	if err != nil {
		gamelog.Error("TMysqlLog ping Error : %s", err.Error())
		return false
	}

	if false == self.CreateLogTable(svrid) {
		return false
	}

	self.tx, err = self.db.Begin()
	if err != nil {
		gamelog.Error("TMysqlLog Start Error : %s", err.Error())
		return false
	}
	self.query = `INSERT INTO gamelog (eventid,	srcid,svrid,chnlid,playerid,level,viplvl,time,param1,param2)VALUES(?,?,?,?,?,?,?,?,?,?);`
	self.stmt, err = self.tx.Prepare(self.query)
	if err != nil {
		gamelog.Error("Start Error : self.tx.Prepare: %s", err.Error())
		return false
	}
	return true
}

func (self *TMysqlLog) WriteLog(pdata []byte) {
	if pdata == nil {
		self.Flush()
		return
	}
	var req msg.MSG_SvrLogData
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("MysqlLog::WriteLog : Message Reader Error!!!!")
		return
	}

	_, err := self.stmt.Exec(req.EventID, req.SrcID, req.SvrID, req.ChnlID, req.PlayerID, req.Level, req.VipLvl, req.Time, req.Param[0], req.Param[1])
	if err != nil {
		gamelog.Error("MysqlLog::WriteLog Error :%s !!!!", err.Error())
	}
	self.writeCnt++
	if self.writeCnt >= self.flushCnt {
		self.Flush()
	}
}

func (self *TMysqlLog) Close() {
	self.tx.Commit()
	self.db.Close()
	self.writeCnt = 0
}

func (self *TMysqlLog) Flush() {
	self.tx.Commit()
	self.stmt.Close()
	var err error
	self.tx, err = self.db.Begin()
	if err != nil {
		gamelog.Error("Start Error : self.tx.Prepare: %s", err.Error())
		return
	}

	self.stmt, err = self.tx.Prepare(self.query)
	if err != nil {
		gamelog.Error("Start Error : self.tx.Prepare: %s", err.Error())
		return
	}
}

func (self *TMysqlLog) SetFlushCnt(cnt int) {
	self.flushCnt = cnt
	if cnt <= 0 {
		self.flushCnt = 100
	}
}
