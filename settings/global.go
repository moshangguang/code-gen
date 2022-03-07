package settings

import (
	"code-gen/utils/files"
	"code-gen/utils/timestamp"
	"encoding/json"
	"sync"
)

type Global interface {
	GetMySQLConnections() MySQLConnections
	GetLombok() Lombok
	ChangLombokData(open bool)
	ChangeLombokGetter(open bool)
	ChangeLombokSetter(open bool)
	ChangeLombokSlf4j(open bool)
	ChangeNoArgsConstructor(open bool)
	ChangeAllArgsConstructor(open bool)
	ChangeLombokToString(open bool)
	ChangeEqualsAndHashCode(open bool)
}

type global struct {
	*sync.RWMutex `json:"-"`
	MySQLConnect  MySQLConnections
	Java          Java
}

var globalSet = &global{
	RWMutex:      new(sync.RWMutex),
	MySQLConnect: make(MySQLConnections, 0),
	Java: Java{
		Lombok: Lombok{
			Data:               false,
			Getter:             true,
			Setter:             true,
			Slf4j:              true,
			NoArgsConstructor:  true,
			AllArgsConstructor: true,
		},
	},
}

func (global *global) SaveMySQLConnection(conn MySQLConnect) {
	conn.CreateTime = timestamp.Now().TimeStamp()
	save := false
	for i := range global.MySQLConnect {
		connection := global.MySQLConnect[i]
		if connection.Name != conn.Name {
			continue
		}
		conn.CreateTime = connection.CreateTime
		global.MySQLConnect[i] = conn
		save = true
	}
	if !save {
		global.MySQLConnect = append(global.MySQLConnect, conn)
	}
	global.Save()
}
func (global *global) GetMySQLConnections() MySQLConnections {
	global.RLock()
	defer global.RUnlock()
	return global.MySQLConnect
}
func (global *global) GetLombok() Lombok {
	global.RLock()
	defer global.RUnlock()
	return global.Java.Lombok
}
func (global *global) ChangLombokData(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.Data = open
	global.Save()
}

func (global *global) ChangeLombokGetter(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.Getter = open
	global.Save()
}

func (global *global) ChangeLombokSetter(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.Setter = open
	global.Save()
}

func (global *global) ChangeLombokSlf4j(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.Slf4j = open
	global.Save()
}

func (global *global) ChangeNoArgsConstructor(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.NoArgsConstructor = open
	global.Save()
}
func (global *global) ChangeAllArgsConstructor(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.AllArgsConstructor = open
	global.Save()
}
func (global *global) ChangeLombokToString(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.ToString = open
	global.Save()
}

func (global *global) ChangeEqualsAndHashCode(open bool) {
	global.Lock()
	defer global.Unlock()
	global.Java.Lombok.EqualsAndHashCode = open
	global.Save()
}

const SettingFile = "code-gen"

func (global *global) Save() {
	go func() {
		global.Lock()
		defer global.Unlock()
		bytes, err := json.Marshal(global)
		if err != nil {
			return
		}
		files.WriteTempFileContent(SettingFile, bytes)
	}()

}

func GetGlobal() Global {
	return globalSet
}

func init() {
	content := files.GetTempFileContent(SettingFile)
	if len(content) == 0 {
		return
	}
	g := &global{
		RWMutex: new(sync.RWMutex),
	}
	err := json.Unmarshal(content, g)
	if err != nil {
		return
	}
	globalSet = g
}
