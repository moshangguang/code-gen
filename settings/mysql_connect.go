package settings

import "code-gen/utils/timestamp"

type MySQLConnect struct {
	Name       string
	Host       string
	Port       int
	User       string
	Password   string
	CreateTime timestamp.Timestamp
}
