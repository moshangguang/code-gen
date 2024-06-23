package ddl

import "time"

type MySQLConnection struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Database  string    `json:"database" xorm:"-"`
	UserName  string    `json:"user_name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
}

func (MySQLConnection) TableName() string {
	return "mysql_connection"
}

type MySQLConnectionSlice []MySQLConnection

func (slice MySQLConnectionSlice) First(filter MySQLConnectionFilter) (MySQLConnection, bool) {
	for _, item := range slice {
		if filter(item) {
			return item, true
		}
	}
	return MySQLConnection{}, false
}

func (slice MySQLConnectionSlice) Index(filter MySQLConnectionFilter) (index int) {
	for i, item := range slice {
		if filter(item) {
			return i
		}
	}
	return -1
}

func (slice MySQLConnectionSlice) Filter(filter MySQLConnectionFilter) MySQLConnectionSlice {
	result := make(MySQLConnectionSlice, 0, len(slice)/2)
	for _, item := range slice {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result
}
func (slice MySQLConnectionSlice) GetNames() []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		result = append(result, item.Name)
	}
	return result
}

type MySQLConnectionFilter func(MySQLConnection) bool
