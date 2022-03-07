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
type MySQLConnections []MySQLConnect

func (connections MySQLConnections) GetNames() []string {
	result := make([]string, 0, len(connections))
	for _, connection := range connections {
		result = append(result, connection.Name)
	}
	return result

}
