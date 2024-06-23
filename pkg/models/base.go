package models

import (
	"code-gen/pkg/models/ddl"
	"code-gen/utils/md5"
	"code-gen/utils/strutils"
	"fmt"
	"github.com/go-xorm/xorm"
	"sync"
)

type MySQLConnectionManager struct {
	sync.RWMutex
	engine map[string]*xorm.Engine
}

func (mgr *MySQLConnectionManager) getConnection(dsnMd5 string) (*xorm.Engine, bool) {
	mgr.RLock()
	defer mgr.RUnlock()
	if len(mgr.engine) == 0 {
		return nil, false
	}
	engine, ok := mgr.engine[dsnMd5]
	return engine, ok
}

func (mgr *MySQLConnectionManager) InitConnection(dsn string) (*xorm.Engine, error) {
	dsnMd5 := md5.GetMd5(dsn)
	mgr.Lock()
	defer mgr.Unlock()
	if mgr.engine == nil {
		mgr.engine = make(map[string]*xorm.Engine)
	}
	engine, ok := mgr.engine[dsnMd5]
	if ok {
		return engine, nil
	}
	orm, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = orm.Ping(); err != nil {
		return nil, err
	}
	mgr.engine[dsnMd5] = orm
	return orm, nil
}
func (mgr *MySQLConnectionManager) LoadConnection(connection ddl.MySQLConnection) (*xorm.Engine, error) {
	var dsn string
	if strutils.IsEmptyString(connection.Database) {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/", connection.UserName, connection.Password, connection.Host, connection.Port)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
			connection.UserName,
			connection.Password,
			connection.Host,
			connection.Port,
			connection.Database,
		)
	}
	dsnMd5 := md5.GetMd5(dsn)
	engine, ok := mgr.getConnection(dsnMd5)
	if ok {
		return engine, nil
	}
	return mgr.InitConnection(dsn)
}
func NewMySQLConnectionManager() *MySQLConnectionManager {
	return &MySQLConnectionManager{}
}

var MySQLConnManager = NewMySQLConnectionManager()
