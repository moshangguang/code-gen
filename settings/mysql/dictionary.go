package mysql

import "code-gen/utils/files"

type dictionary struct {
	Connections Connections
}

func (dictionary *dictionary) Get(name string) (Connection, bool) {
	return dictionary.Connections.GetByName(name)
}

func (dictionary *dictionary) Save(connection Connection) {
	add := true
	for i := range dictionary.Connections {
		if dictionary.Connections[i].Name == connection.Name {
			add = false
			dictionary.Connections[i] = connection
		}
	}
	if add {
		dictionary.Connections = append(dictionary.Connections, connection)
	}
	dictionary.GlobalSave()
}
func (dictionary *dictionary) Remove(name string) {
	dictionary.Connections = dictionary.Connections.Filter(func(connection Connection) bool {
		return connection.Name != name
	})
	dictionary.GlobalSave()
}
func (dictionary *dictionary) GlobalSave() {
	_ = files.Marshal(SettingFile, *dictionary)
}

func (dictionary *dictionary) All() Connections {
	return dictionary.Connections
}

var dict = &dictionary{
	Connections: []Connection{},
}

func GetDictionary() Dictionary {
	return dict
}

const SettingFile = "mysql_code_gen.ini"

func init() {
	_, _ = files.Unmarshal(SettingFile, &dict)
}
