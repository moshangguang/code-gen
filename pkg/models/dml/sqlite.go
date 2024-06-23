package dml

import (
	"code-gen/pkg/models/ddl"
	"code-gen/utils/fileUtils"
	"code-gen/utils/runtime"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wxnacy/wgo/arrays"
	"path/filepath"
)

var engine *xorm.Engine

func Init() {
	var err error
	dbPath := filepath.Join(fileUtils.GetRootDir(), "code-gen", "codegen.db")
	engine, err = xorm.NewEngine("sqlite3", dbPath)
	runtime.PanicError(err)
	err = engine.Ping()
	runtime.PanicError(err)
	tableNames, err := GetAllTableNames()
	runtime.PanicError(err)
	if arrays.ContainsString(tableNames, ddl.MySQLConnection{}.TableName()) == -1 {
		InitMySQLConnection()
	}
	if arrays.ContainsString(tableNames, ddl.Config{}.TableName()) == -1 {
		InitConfig()
	}
}
func InitMySQLConnection() {
	exec, err := engine.Exec(`CREATE TABLE "mysql_connection" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT NOT NULL,
  "host" TEXT NOT NULL,
  "port" TEXT NOT NULL,
  "user_name" TEXT NOT NULL,
  "password" TEXT NOT NULL,
  "created_at" DATE NOT NULL,
  "updated_at" DATE NOT NULL,
  "deleted_at" DATE
)
`)
	runtime.PanicError(err)
	_, err = exec.RowsAffected()
	runtime.PanicError(err)
	exec, err = engine.Exec(`CREATE UNIQUE INDEX "UEQ_name"
ON "mysql_connection" (
  "name"
)`)
	runtime.PanicError(err)
	_, err = exec.RowsAffected()
	runtime.PanicError(err)

}

func InitConfig() {
	exec, err := engine.Exec(`CREATE TABLE "config" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" TEXT NOT NULL,
  "value" TEXT NOT NULL,
  "created_at" DATE NOT NULL,
  "updated_at" DATE NOT NULL,
  "deleted_at" DATE
)
`)
	runtime.PanicError(err)
	_, err = exec.RowsAffected()
	runtime.PanicError(err)
	exec, err = engine.Exec(`CREATE UNIQUE INDEX "UEQ_config_name"
ON "config" (
  "name"
)`)
	runtime.PanicError(err)
	_, err = exec.RowsAffected()
	runtime.PanicError(err)
}

func GetAllTableNames() ([]string, error) {
	result := make([]string, 0)
	tables, err := engine.Dialect().GetTables()
	if err != nil {
		return nil, err
	}
	for _, t := range tables {
		result = append(result, t.Name)
	}
	return result, nil

}
