package mysql_writer

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type MysqlWriterData struct {
	Database string
	Sql      string
}

type MysqlWriter struct {
	databaseConn  map[string]*sql.DB
	databaseCache map[string]chan string
	WorkerNumber  int
	Username      string
	Password      string
	Address       string
	Port          int
}

func (mw *MysqlWriter) CreateMysqlDB(database string) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mw.Username, mw.Password, mw.Address, &mw.Port, database)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (mw *MysqlWriter) SyncWrite(data *MysqlWriterData) error {
	database := data.Database
	db, exist := mw.databaseConn[database]
	if !exist {
		var err error
		db, err = mw.CreateMysqlDB(database)
		if err != nil {
			return err
		}
		mw.databaseConn[database] = db
	}
	result, err := db.Exec(data.Sql)
	if err != nil {
		return err
	}
	affect, _ := result.RowsAffected()
	logrus.Debugf("exec %s have success, effect row %d", data.Sql, affect)
	return nil
}

func (mw *MysqlWriter) NewAsyncWriteWorker(database string, workerNumber, sqlCacheSize int) error {
	cache, exist := mw.databaseCache[database]
	if !exist {
		cache = make(chan string, sqlCacheSize)
		mw.databaseCache[database] = cache
		db, err := mw.CreateMysqlDB(database)
		if err != nil {
			return err
		}
		mw.databaseConn[database] = db
	}
	return nil
}

func (mw *MysqlWriter) asyncWriter(workerId int, database string) {
	db := mw.databaseConn[database]
	cache := mw.databaseCache[database]
	for sqlString := range cache {
		result, err := db.Exec(sqlString)
		if err != nil {
			logrus.Errorf("exec sql %s have an err: %v", err)
			continue
		}
		affect, err := result.RowsAffected()
		if err != nil {
			logrus.Errorf("get rows affected have an err: %v", err)
			continue
		}
		logrus.Infof("exec success, workerId: %d, database: %v, affect rows: %d", workerId, database, affect)
	}
}

func (mw *MysqlWriter) AsyncWrite(data *MysqlWriterData) {
	database := data.Database
	cache, exist := mw.databaseCache[database]
	if !exist {
		errMsg := fmt.Sprintf("You must call NewAsyncWriteWorker befor call AsyncWrite")
		panic(errMsg)
	}
	cache <- data.Sql
}
