package tutorials

import (
	"code-gen/tutorials/mysql"
	"fyne.io/fyne/v2"
)

type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var Tutorials = map[string]Tutorial{
	"MySQL": {
		Title: "MySQL",
		View:  mysql.MySQLScene,
	},
	"MySQLJava": {
		Title: "Java",
		View:  mysql.JavaScene,
	},
}

var TutorialIndex = map[string][]string{
	"": {
		"MySQL",
	},
	"MySQL": {
		"MySQLJava",
	},
}
