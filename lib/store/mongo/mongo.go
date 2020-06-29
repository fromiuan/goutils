package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo"
	_ "github.com/go-sql-driver/mysql"
)

//Container for the database
var MConn *mgo.Database

// SetupMongo
func SetupMongo(mongoURL, mongoDBName string) {
	MConn = getMgoConnection(mongoURL, mongoDBName)
	if MConn == nil {
		fmt.Errorf("can not initialized mongo connection:%s and %s\n", mongoURL, mongoDBName)
		os.Exit(1)
	}
}

func getMgoConnection(url, dbname string) *mgo.Database {
	retry := 3
	if retry >= 0 {
		session, err := mgo.DialWithTimeout(url, time.Second*3) //连接数据库
		if err != nil {
			log.Printf("sleep 3 sec; mgo dial err:%v\n", err)
			time.Sleep(3 * time.Second)
			retry--
		} else {
			return session.DB(dbname)
		}
	}
	return nil
}
