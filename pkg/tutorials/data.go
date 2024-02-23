package tutorials

import (
	"code-gen/pkg/tutorials/mysql"
	"fyne.io/fyne/v2"
)

type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var Tutorials = map[string]Tutorial{
	"MySQL": {
		Title: "MySQL",
		View:  mysql.MySQLConnections,
	},
	"MySQLAddConn": {
		Title: "新增连接",
		View:  mysql.MySQLScene,
	},
	"MySQLConnManager": {
		Title: "管理连接",
		View:  mysql.MySQLConnections,
	},
	"Golang": {
		Title: "Golang",
		View:  mysql.MySQLConnections,
	},
	//"MySQLJava": {
	//	Title: "Java",
	//	View:  mysql.JavaScene,
	//},
}

var TutorialIndex = map[string][]string{
	"": {
		"MySQL",
		"Golang",
	},
	"MySQL": {
		"MySQLAddConn",
		"MySQLConnManager",
	},
}
