package tutorials

import (
	"code-gen/constants"
	"code-gen/pkg/db"
	"code-gen/pkg/models"
	"code-gen/pkg/models/ddl"
	"code-gen/settings/mysql"
	"code-gen/utils/strutils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"strings"
	"time"
)

func MySQLEditScene(win fyne.Window) fyne.CanvasObject {
	data := GetMySQLConnSlice()
	if len(data) == 0 {
		return logoTextScreen("暂无连接")
	}
	connEntry := widget.NewSelect([]string{}, func(s string) {

	})
	hostEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	portEntry.Validator = validation.NewRegexp(`^[0-9]*$`, "请输入正确端口号")
	portEntry.SetText("3306")
	//portEntry.SetPlaceHolder("3306")
	userEntry := widget.NewEntry()
	userEntry.SetText("root")
	passwordEntry := widget.NewPasswordEntry()
	w := MySQLConnWidget{
		connEntry:     connEntry,
		hostEntry:     hostEntry,
		portEntry:     portEntry,
		userEntry:     userEntry,
		passwordEntry: passwordEntry,
	}
	connEntry.OnChanged = func(s string) {
		reloadConnSelect(w, nil, s)
	}
	reloadConnSelect(w, data, "")
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
	resetBtn := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {
		data = GetMySQLConnSlice()
		data = data.Filter(func(connection ddl.MySQLConnection) bool {
			return connection.Name != connEntry.Selected
		})
		db.GetDatabase().Save(constants.MySQLConnection, &data)
		reloadConnSelect(w, nil, "")
	})
	saveBtn := widget.NewButtonWithIcon("保存", theme.ConfirmIcon(), func() {
		data = GetMySQLConnSlice()
		index := data.Index(func(connection ddl.MySQLConnection) bool {
			return connection.Name == connEntry.Selected
		})
		if index == -1 {
			return
		}
		conn := data[index]
		conn.Host = strings.TrimSpace(hostEntry.Text)
		conn.Port = cast.ToInt(strings.TrimSpace(portEntry.Text))
		conn.UserName = strings.TrimSpace(userEntry.Text)
		conn.Password = strings.TrimSpace(passwordEntry.Text)
		_, err := models.MySQLConnManager.LoadConnection(conn)
		if err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("连接输出库出错,error:%s", err.Error()), win)
			return
		}
		data[index] = conn
		db.GetDatabase().Save(constants.MySQLConnection, &data)
		dialog.ShowInformation("成功", "保存成功", win)
	})
	saveBtn.Importance = widget.HighImportance
	form.Append("", container.NewBorder(nil, nil, nil, container.NewHBox(resetBtn, saveBtn)))
	return form
}

type MySQLConnWidget struct {
	connEntry     *widget.Select
	hostEntry     *widget.Entry
	portEntry     *widget.Entry
	userEntry     *widget.Entry
	passwordEntry *widget.Entry
}

func reloadConnSelect(connWidget MySQLConnWidget, connSlice ddl.MySQLConnectionSlice, name string) {
	if connSlice == nil {
		connSlice = GetMySQLConnSlice()
	}
	if strutils.IsEmptyString(name) && len(connSlice) != 0 {
		name = connSlice[0].Name
	}

	names := connSlice.GetNames()
	connEntry := connWidget.connEntry
	hostEntry := connWidget.hostEntry
	portEntry := connWidget.portEntry
	userEntry := connWidget.userEntry
	passwordEntry := connWidget.passwordEntry
	connEntry.Options = names
	connEntry.Selected = name
	hostEntry.Text = ""
	portEntry.Text = ""
	userEntry.Text = ""
	passwordEntry.Text = ""
	conn, ok := connSlice.First(func(connection ddl.MySQLConnection) bool {
		return connection.Name == name
	})
	if ok {
		hostEntry.Text = conn.Host
		portEntry.Text = cast.ToString(conn.Port)
		userEntry.Text = conn.UserName
		passwordEntry.Text = conn.Password
	}
	connEntry.Refresh()
	hostEntry.Refresh()
	portEntry.Refresh()
	userEntry.Refresh()
	passwordEntry.Refresh()

}
func MySQLAddScene(win fyne.Window) fyne.CanvasObject {
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
	passwordEntry := widget.NewPasswordEntry()
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
		data := GetMySQLConnSlice()
		_, ok := data.First(func(connection ddl.MySQLConnection) bool {
			return connection.Name == name
		})
		if ok {
			dialog.ShowInformation("错误", "连接名已存在", win)
			return
		}

		conn := ddl.MySQLConnection{
			Id:       time.Now().UnixNano(),
			Name:     name,
			Host:     host,
			Port:     port,
			UserName: user,
			Password: pwd,
		}
		_, err = models.MySQLConnManager.LoadConnection(conn)
		if err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("连接输出库出错,error:%s", err.Error()), win)
			return
		}
		data = append(data, ddl.MySQLConnection{
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
