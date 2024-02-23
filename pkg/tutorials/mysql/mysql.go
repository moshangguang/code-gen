package mysql

import (
	"code-gen/constants"
	"code-gen/pkg/db"
	"code-gen/pkg/models"
	"code-gen/settings/mysql"
	"code-gen/utils/strutils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"strings"
	"time"
)

func MySQLConnections(win fyne.Window) fyne.CanvasObject {
	return logoTextScreen("暂无连接")
}
func MySQLScene(win fyne.Window) fyne.CanvasObject {
	dictionary := mysql.GetDictionary()
	connEntry := widget.NewEntry()
	connEntry.PlaceHolder = "请输入连接名"
	hostEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	portEntry.Validator = validation.NewRegexp(`^[0-9]*$`, "请输入正确端口号")
	portEntry.SetText("3306")
	//portEntry.SetPlaceHolder("3306")
	userEntry := widget.NewEntry()
	userEntry.SetText("root")
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
				Text:   "连接名:",
				Widget: connEntry,
			},
			{
				Text:   "主机:",
				Widget: hostEntry,
			},
			{
				Text:   "端口:",
				Widget: portEntry,
			},
			{
				Text:   "用户名:",
				Widget: userEntry,
			},
		},
	}
	form.Append("密码:", passwordEntry)
	resetBtn := widget.NewButtonWithIcon("清空", theme.ContentClearIcon(), func() {
		reset()
	})
	saveBtn := widget.NewButtonWithIcon("保存", theme.ConfirmIcon(), func() {
		user := strings.TrimSpace(userEntry.Text)
		pwd := strings.TrimSpace(passwordEntry.Text)
		host := strings.TrimSpace(hostEntry.Text)
		name := strings.TrimSpace(connEntry.Text)
		if strutils.IsEmptyString(name) {
			dialog.ShowInformation("错误", "请输入连接名", win)
			return
		}

		port, err := cast.ToIntE(strings.TrimSpace(portEntry.Text))
		if err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("端口出错,error:%s", err.Error()), win)
			return
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, pwd, host, port)
		orm, err := xorm.NewEngine("mysql", dsn)
		err = orm.Ping()
		if err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("连接输出库出错,error:%s", err.Error()), win)
			return
		}
		data := make([]models.MySQLConnection, 0)
		db.GetDatabase().Unmarshal(constants.MySQLConnection, &data)
		data = append(data, models.MySQLConnection{
			Id:       time.Now().UnixNano(),
			Name:     name,
			Host:     host,
			Port:     port,
			UserName: user,
			Password: pwd,
		})
		db.GetDatabase().Save(constants.MySQLConnection, data)
		dialog.ShowInformation("成功", "保存成功", win)

	})
	saveBtn.Importance = widget.HighImportance
	form.Append("", container.NewBorder(nil, nil, nil, container.NewHBox(resetBtn, saveBtn)))
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
