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
type MySQLConnectFilter func(MySQLConnect) bool

func (connections MySQLConnections) GetNames() []string {
	result := make([]string, 0, len(connections))
	for _, connection := range connections {
		result = append(result, connection.Name)
	}
	return result

}
func (connections MySQLConnections) Filter(filter MySQLConnectFilter) MySQLConnections {
	result := make(MySQLConnections, 0, len(connections))
	for _, connection := range connections {
		if filter(connection) {
			result = append(result, connection)
		}
	}
	return result
}
func (connections MySQLConnections) ContainsName(name string) bool {
	for _, connection := range connections {
		if connection.Name == name {
			return true
		}
	}
	return false
}

func (connections MySQLConnections) RemoveByName(name string) MySQLConnections {
	return connections.Filter(func(connect MySQLConnect) bool {
		return connect.Name != name
	})
}
