package mysql

import (
	_ "github.com/go-sql-driver/mysql"

	"sync"
	"time"

	"github.com/gocraft/dbr"
	"github.com/pkg/errors"
)

var (
	defaultMaxOpenConns = 50     // 最大连接数
	defaultMaxIdleConns = 10     // 最大空闲连接
	defaultInstanceName = "base" // 默认的实例名称

	defaultSession *dbr.Session
	pool           = make(map[string]*dbr.Session)
	l              sync.RWMutex
)

type Config struct {
	Instance   string
	DriverName string
	DataSource string
}

func Init(conf []*Config) error {
	l.Lock()
	defer l.Unlock()

	var err error
	for _, v := range conf {
		if _, ok := pool[v.Instance]; !ok {
			pool[v.Instance], err = NewSession(v.DriverName, v.DataSource)
			if err != nil {
				return err
			}
		}
	}

	if session, ok := pool[defaultInstanceName]; ok {
		defaultSession = session
	}
	return nil
}

func NewSession(driverName string, dataSource string) (*dbr.Session, error) {
	conn, err := dbr.Open(driverName, dataSource, NewDBLog(driverName))
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(defaultMaxOpenConns)
	conn.SetMaxIdleConns(defaultMaxIdleConns)
	conn.SetConnMaxLifetime(time.Second * 60)
	return conn.NewSession(nil), nil
}

// 获取默认session
func DefaultSession() *dbr.Session {
	return defaultSession
}

// 设置数据库最大连接数
func SetMaxOpenConns(instanceName string, maxOpenConns int) error {
	session, err := GetInstanceSession(instanceName)

	if err != nil {
		return err
	}
	session.SetMaxOpenConns(maxOpenConns)
	return nil
}

// 设置数据库最大空闲连接数
func SetMaxIdleConns(instanceName string, maxIdleConns int) error {
	session, err := GetInstanceSession(instanceName)
	if err != nil {
		return err
	}
	session.SetMaxIdleConns(maxIdleConns)
	return nil
}

// 获取指定实例名称的session
func GetInstanceSession(name string) (*dbr.Session, error) {
	l.RLock()
	defer l.RUnlock()
	if session, ok := pool[name]; ok {
		return session, nil
	}
	return nil, errors.New("unknown DataBase alias name :" + name)
}
