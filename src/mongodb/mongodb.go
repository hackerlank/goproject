package mongodb

import (
	"gamelog"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	G_db_Addr       string
	G_db_Connection *mgo.Session = nil
)

func Init(addr string) bool {

	G_db_Addr = addr
	G_db_Connection = nil

	var err error
	G_db_Connection, err = mgo.Dial(G_db_Addr)
	if err != nil {
		gamelog.Error(err.Error())
		panic("Mongodb Init Failed " + err.Error())
	}

	G_db_Connection.SetPoolLimit(20)

	return true
}

func InitWithAuth(addr string, username string, password string) bool {

	G_db_Addr = addr
	G_db_Connection = nil

	mgoDialInfo := mgo.DialInfo{
		Addrs:     []string{addr},
		Timeout:   5 * time.Second,
		Username:  username,
		Password:  password,
		PoolLimit: 20,
	}

	var err error
	G_db_Connection, err = mgo.DialWithInfo(&mgoDialInfo)
	if err != nil {
		gamelog.Error(err.Error())
		panic("Mongodb Init Failed " + err.Error())
	}

	return true
}

//获取MongoDB连接
func GetDBSession() *mgo.Session {

	if G_db_Connection == nil {
		gamelog.Error("GetDBSession Failed, G_db_Connection is nil!!")
		panic("db connections is nil!!")
	}

	return G_db_Connection.Clone()
}

//更新多条记录
func UpdateToDBAll(dbname string, collection string, search bson.M, stuff bson.M) bool {
	s := GetDBSession()
	defer s.Close()
	coll := s.DB(dbname).C(collection)
	_, err := coll.UpdateAll(search, stuff)
	if err != nil {
		gamelog.Error3("UpdateToDB Failed: DB:[%s] Collection:[%s] search:[%v], stuff:[%v], Error:%v", dbname, collection, search, stuff, err.Error())
		return false
	}

	return true
}

//更新一条记录
func UpdateToDB(dbname string, collection string, search bson.M, stuff bson.M) bool {
	s := GetDBSession()
	defer s.Close()
	coll := s.DB(dbname).C(collection)
	err := coll.Update(search, stuff)
	if err != nil {
		gamelog.Error3("UpdateToDB Failed: DB:[%s] Collection:[%s] search:[%v], stuff:[%v], Error:%v", dbname, collection, search, stuff, err.Error())
		return false
	}

	return true
}

//更新一条记录  某个字段的值
//增加了数组修改器的修改，修改数组中某个索引的值
//err := coll.UpdateToDB("PlayerDB", "PlayerHero", bson.M{"herolist.heroid": 3}, bson.M{"$set": bson.M{"herolist.$.skilllist": skill}})
//err := UpdateArrayField("PlayerDB", "PlayerHero", bson.M{"herolist.heroid": 3}, "herolist.$.skilllist",  skill)
//一次更新全部的技能，但不用更新所有的集合，而是更新一条记录，数据量要小得多。
//mongodb.UpdateToDB("PlayerDB", "PlayerHero", bson.M{"herolist.heroid": 3}, bson.M{"$set": bson.M{"herolist.$.skilllist": skill}})

//增加一个字段值
func IncFieldValue(dbname string, collection string, search bson.M, field string, value int) bool {
	s := GetDBSession()
	defer s.Close()
	coll := s.DB(dbname).C(collection)
	err := coll.Update(search, bson.M{"$inc": bson.M{field: value}})
	if err != nil {
		gamelog.Error3("IncFieldValue Failed: DB:[%s] Collection:[%s] Error:[%s]", dbname, collection, err.Error())
		return false
	}

	return true
}

//插入一条记录
func InsertToDB(dbname string, collection string, data interface{}) bool {
	s := GetDBSession()
	defer s.Close()
	coll := s.DB(dbname).C(collection)
	err := coll.Insert(&data)
	if err != nil {
		if !mgo.IsDup(err) {
			gamelog.Error("InsertToDB Failed: DB:[%s] Collection:[%s] Error:[%s]", dbname, collection, err.Error())
		} else {
			gamelog.Warn("InsertToDB Failed: DB:[%s] Collection:[%s] Error:[%s]", dbname, collection, err.Error())
		}
		return false
	}

	return true
}

//删掉指定的一条记录
func RemoveFromDB(dbname string, collection string, search bson.M) error {
	s := GetDBSession()
	defer s.Close()

	coll := s.DB(dbname).C(collection)

	return coll.Remove(search)
}

//查询指定的记录是否存在
func IsRecordExist(dbname string, collection string, search bson.M) bool {
	s := GetDBSession()
	defer s.Close()
	coll := s.DB(dbname).C(collection)
	nCount, err := coll.Find(search).Count()
	if err == mgo.ErrNotFound {
		return false
	}

	if err == nil {
		return nCount > 0
	}

	panic(err.Error())

	return false
}

//删掉指定的一条记录
func AddToArray(dbname string, collection string, search bson.M, fieldname string, data interface{}) bool {
	return UpdateToDB(dbname, collection, search, bson.M{"$push": bson.M{fieldname: data}})
}

//删掉指定的一条记录
func RemoveFromArray(dbname string, collection string, search bson.M, fieldname string, data interface{}) bool {
	return UpdateToDB(dbname, collection, search, bson.M{"$pull": bson.M{fieldname: data}})
}

//删除数组中指定索引的记录
func ArrayRemoveIndex(dbname string, collection string, search bson.M, fieldname string, data interface{}) bool {
	return UpdateToDB(dbname, collection, search, bson.M{"$pull": bson.M{fieldname: data}})
}

//! 查询一条数据
func Find(dbName string, tableName string, find string, find_value interface{}, data interface{}) int {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Find(bson.M{find: find_value}).One(data)
	if err != nil {
		if err == mgo.ErrNotFound {
			gamelog.Warn("Not Find dbName: %s  ntable: %s  find: %s:%v", dbName, tableName, find, find_value)
			return 1
		}

		gamelog.Error3("Find error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v \r\n",
			err.Error(), dbName, tableName, find, find_value)

		return -1
	}

	return 0
}

//! 排序查找
//! order 1 -> 正序  -1 -> 倒序
func Find_Sort(dbName string, tableName string, find string, order int, number int, lst interface{}) int {
	db_session := GetDBSession()
	defer db_session.Close()

	strSort := ""
	if order == 1 {
		strSort = "+" + find
	} else {
		strSort = "-" + find
	}

	collection := db_session.DB(dbName).C(tableName)
	query := collection.Find(nil).Sort(strSort).Limit(number)

	err := query.All(lst)
	if err != nil {
		if err == mgo.ErrNotFound {
			gamelog.Warn("Not Find")
			return 1
		}

		gamelog.Error3("Find_Sort error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s \r\norder: %d\r\nlimit: %d\r\n",
			err.Error(), dbName, tableName, find, order, number)
		return -1
	}

	return 0
}
