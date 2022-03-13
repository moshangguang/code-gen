package mysql

import (
	"code-gen/utils/strutils"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

type Dictionary interface {
	All() Connections
	Get(name string) (Connection, bool)
	Save(connection Connection)
	Remove(name string)
	GlobalSave()
}

type Field struct {
	gorm.Model
	ColumnName string
	DataType   string
}
type Connection struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (connection Connection) Use(database string) Connection {
	connection.Database = database
	return connection
}
func (connection Connection) MustGetDB() (*gorm.DB, error) {
	dbName := strings.TrimSpace(connection.Database)
	if len(dbName) == 0 {
		dbName = "test"
	}
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		connection.User,
		connection.Password,
		connection.Host,
		connection.Port,
		dbName,
	)
	return gorm.Open("mysql", url)
}
func (connection Connection) GetFields(table string) []Field {
	result := make([]Field, 0)
	db, ok := connection.GetDB()
	if !ok {
		return result
	}

	defer db.Close()
	db.Raw("select column_name,data_type from information_schema.columns where table_name=? and table_schema=?",
		table,
		connection.Database,
	).Scan(&result)
	for _, field := range result {
		if len(field.ColumnName) == 0 || len(field.DataType) == 0 {
			return make([]Field, 0)
		}
	}
	return result
}
func (connection Connection) GetDB() (*gorm.DB, bool) {
	db, err := connection.MustGetDB()
	if err == nil {
		return db, true
	}
	return nil, false
}

func (connection Connection) MustActive() error {
	db, err := connection.MustGetDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (connection Connection) SimpleRaw(sql string) []string {
	result := make([]string, 0)
	db, ok := connection.GetDB()
	if !ok {
		return result
	}
	defer db.Close()
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return result
	}
	for rows.Next() {
		var item string
		if err = rows.Scan(&item); err != nil {
			return result
		} else {
			result = append(result, item)
		}
	}
	return result
}

func (connection Connection) SimplePatternRaw(sql, pattern string) []string {
	data := connection.SimpleRaw(sql)
	if len(pattern) == 0 {
		return data
	}
	result := make([]string, 0, len(data)/2)
	for _, d := range data {
		if strings.Contains(d, pattern) {
			result = append(result, d)
		}
	}
	return result
}
func (connection Connection) PatternTables(pattern string, limit int) []string {
	data := connection.SimpleRaw("show tables")
	result := strutils.PrefixPattern(data, pattern)
	if len(result) <= limit {
		return result
	}
	return result[:limit]
}
func (connection Connection) PatternDatabases(pattern string, limit int) []string {
	data := connection.SimpleRaw("show databases")
	result := strutils.PrefixPattern(data, pattern)
	if len(result) <= limit {
		return result
	}
	return result[:limit]
}

type Connections []Connection
type ConnectionFilter func(Connection) bool

func (connections Connections) GetNames() []string {
	result := make([]string, 0, len(connections))
	for _, connection := range connections {
		result = append(result, connection.Name)
	}
	return result

}
func (connections Connections) First(filter ConnectionFilter) (Connection, bool) {
	for _, connection := range connections {
		if filter(connection) {
			return connection, true
		}
	}
	return Connection{}, false
}
func (connections Connections) Filter(filter ConnectionFilter) Connections {
	result := make(Connections, 0, len(connections))
	for _, connection := range connections {
		if filter(connection) {
			result = append(result, connection)
		}
	}
	return result
}
func (connections Connections) ContainsName(name string) bool {
	_, ok := connections.First(func(connection Connection) bool {
		return connection.Name == name
	})
	return ok
}
func (connections Connections) GetByName(name string) (Connection, bool) {
	return connections.First(func(connection Connection) bool {
		return connection.Name == name
	})
}
func (connections Connections) RemoveByName(name string) Connections {
	return connections.Filter(func(connection Connection) bool {
		return connection.Name != name
	})
}
