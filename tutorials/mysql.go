package tutorials

import (
	"code-gen/settings"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

func MySQLScene(win fyne.Window) fyne.CanvasObject {
	global := settings.GetGlobal()
	connects := global.GetMySQLConnections()
	connEntry := widget.NewSelectEntry(connects.GetNames())
	connEntry.PlaceHolder = "请输入或选择连接名"
	hostEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	portEntry.SetText("3306")
	portEntry.SetPlaceHolder("3306")
	userEntry := widget.NewEntry()
	userEntry.SetText("root")
	userEntry.SetPlaceHolder("root")
	passwordEntry := widget.NewEntry()
	largeText := widget.NewMultiLineEntry()
	largeText.Hidden = true
	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "连接名",
				Widget: connEntry,
			},
			{
				Text:   "主机",
				Widget: hostEntry,
			},
			{
				Text:   "端口",
				Widget: portEntry,
			},
			{
				Text:   "用户名",
				Widget: userEntry,
			},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		CancelText: "删除",
		OnSubmit: func() {
			user := userEntry.Text
			pwd := passwordEntry.Text
			host := hostEntry.Text
			conn := connEntry.Text
			if conn == "" {
				dialog.ShowInformation("错误", "请输入连接名", win)
				return
			}

			portText := portEntry.Text
			port, err := strconv.Atoi(portText)
			if err != nil {
				dialog.ShowInformation("错误", fmt.Sprintf("端口出错,error:%s", err.Error()), win)
				return
			}
			url := fmt.Sprintf("%s:%s@tcp(%s:%d)/test?charset=utf8", user, pwd, host, port)
			db, err := gorm.Open("mysql", url)
			if err != nil {
				dialog.ShowInformation("错误", fmt.Sprintf("连接MySQL出错,error:%s", err.Error()), win)
				return
			}
			defer db.Close()
		},
		SubmitText: "保存",
	}
	form.Append("密码", passwordEntry)
	return form
}
