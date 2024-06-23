package tutorials

import (
	"fyne.io/fyne/v2"
)

type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var Tutorials = map[string]Tutorial{
	"MySQL": {
		Title: "MySQL",
		View:  MySQLEditScene,
	},
	"MySQLAddConn": {
		Title: "新增连接",
		View:  MySQLAddScene,
	},
	"MySQLConnManager": {
		Title: "管理连接",
		View:  MySQLEditScene,
	},
	"Xorm": {
		Title: "Xorm",
		View:  XormScene,
	},
}

var TutorialIndex = map[string][]string{
	"": {
		"MySQL",
		"Xorm",
	},
	"MySQL": {
		"MySQLAddConn",
		"MySQLConnManager",
	},
}
