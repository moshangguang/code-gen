package settings

import (
	"code-gen/enum/Theme"
	"code-gen/utils/files"
	"encoding/json"
	"fyne.io/fyne/v2"
)

type Global struct {
	MySQLConnect []MySQLConnect
	Java         Java
	Theme        Theme.Theme
}

func (global Global) GetTheme() fyne.Theme {
	return global.Theme.ToTheme()
}

var globalSet = &Global{
	MySQLConnect: make([]MySQLConnect, 0),
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
	Theme: Theme.Dark,
}

func (global *Global) ChangeTheme(theme Theme.Theme) {
	global.Theme = theme
	global.Save()
}

const SettingFile = "code-gen"

func (global Global) Save() {
	bytes, err := json.Marshal(global)
	if err != nil {
		return
	}
	files.WriteTempFileContent(SettingFile, bytes)
}

func GetGlobal() *Global {
	return globalSet
}

func init() {
	content := files.GetTempFileContent(SettingFile)
	if len(content) == 0 {
		return
	}
	g := new(Global)
	err := json.Unmarshal(content, g)
	if err != nil {
		return
	}
	globalSet = g
}
