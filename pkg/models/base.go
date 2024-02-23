package models

type MySQLConnection struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database,omitempty"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
