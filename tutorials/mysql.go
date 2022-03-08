package tutorials

import (
	"code-gen/settings"
	"code-gen/utils/timestamp"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"strings"
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
	connEntry.OnChanged = func(s string) {
		connection, ok := global.GetMySQLConnection(s)
		if !ok {
			return
		}
		hostEntry.Text = connection.Host
		hostEntry.Refresh()
		portEntry.Text = fmt.Sprintf("%d", connection.Port)
		portEntry.Refresh()
		userEntry.Text = connection.User
		userEntry.Refresh()
		passwordEntry.Text = connection.Password
		passwordEntry.Refresh()
	}
	reset := func() {
		connEntry.Text = ""
		connEntry.SetOptions(global.GetMySQLConnections().GetNames())
		connEntry.Refresh()
		hostEntry.Text = ""
		hostEntry.Refresh()
		portEntry.Text = "3306"
		portEntry.Refresh()
		userEntry.Text = "root"
		userEntry.Refresh()
		passwordEntry.Text = ""
		passwordEntry.Refresh()
	}
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
			ok := global.RemoveMySQLConnection(connEntry.Text)
			reset()
			if ok {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "MySQL",
					Content: "删除连接成功...",
				})
			}

		},
		CancelText: "删除",
		OnSubmit: func() {
			user := strings.TrimSpace(userEntry.Text)
			pwd := strings.TrimSpace(passwordEntry.Text)
			host := strings.TrimSpace(hostEntry.Text)
			conn := strings.TrimSpace(connEntry.Text)
			if conn == "" {
				dialog.ShowInformation("错误", "请输入连接名", win)
				return
			}

			portText := strings.TrimSpace(portEntry.Text)
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
			global.SaveMySQLConnection(settings.MySQLConnect{
				Name:       conn,
				Host:       host,
				Port:       port,
				User:       user,
				Password:   pwd,
				CreateTime: timestamp.Now().TimeStamp(),
			})
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "MySQL",
				Content: "MySQL数据库连接成功...",
			})
		},
		SubmitText: "保存",
	}
	form.Append("密码", passwordEntry)
	form.Append("", container.NewBorder(nil, nil, nil, widget.NewButton("重置", func() {
		reset()
	})),
	)
	return form
}
