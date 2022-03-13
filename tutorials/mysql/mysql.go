package mysql

import (
	"code-gen/settings/mysql"
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
	dictionary := mysql.GetDictionary()
	connections := dictionary.All()
	connEntry := widget.NewSelectEntry(connections.GetNames())
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
		conn, ok := dictionary.Get(s)
		if !ok {
			return
		}
		hostEntry.Text = conn.Host
		hostEntry.Refresh()
		portEntry.Text = fmt.Sprintf("%d", conn.Port)
		portEntry.Refresh()
		userEntry.Text = conn.User
		userEntry.Refresh()
		passwordEntry.Text = conn.Password
		passwordEntry.Refresh()
	}
	reset := func() {
		connEntry.Text = ""
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
			if len(connEntry.Text) == 0 {
				return
			}
			dictionary.Remove(connEntry.Text)
			reset()
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
			mysqlConn := mysql.Connection{
				Name:     conn,
				Host:     host,
				Port:     port,
				User:     user,
				Password: pwd,
			}
			if err = mysqlConn.MustActive(); err != nil {
				dialog.ShowInformation("错误", err.Error(), win)
				return
			}
			dictionary.Save(mysqlConn)
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

func GetMySQLDatabase(conn mysql.Connection, dbname ...string) (*gorm.DB, error) {
	dbName := "test"
	if len(dbname) != 0 {
		dbName = dbname[0]
	}
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conn.User,
		conn.Password,
		conn.Host,
		conn.Port,
		dbName,
	)
	return gorm.Open("mysql", url)

}
