package tutorials

import (
	"code-gen/pkg/models"
	"code-gen/pkg/models/ddl"
	"code-gen/pkg/models/dml"
	"code-gen/utils/strutils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cast"
	"github.com/wxnacy/wgo/arrays"
	"strings"
)

var (
	mySQLConnectionModel = new(dml.MySQLConnectionModel)
	configModel          = new(dml.ConfigModel)
)

func MySQLEditScene(win fyne.Window) fyne.CanvasObject {
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
	saveBtn := widget.NewButtonWithIcon("保存", theme.ConfirmIcon(), func() {

	})
	deleteBtn := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {

	})
	w := MySQLConnWidget{
		connEntry:     connEntry,
		hostEntry:     hostEntry,
		portEntry:     portEntry,
		userEntry:     userEntry,
		passwordEntry: passwordEntry,
		saveBtn:       saveBtn,
		resetBtn:      deleteBtn,
	}
	connEntry.OnChanged = func(s string) {
		reloadConnSelect(w, s)
	}
	reloadConnSelect(w, "")
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
	deleteBtn.OnTapped = func() {
		dialog.ShowCustomConfirm("删除账号", "确认", "取消", widget.NewLabel("确认删除连接?"), func(b bool) {
			if !b {
				return
			}
			err := mySQLConnectionModel.DeleteByName(connEntry.Selected)
			if err != nil {
				dialog.ShowInformation("错误", "删除连接出错,err:"+err.Error(), win)
				return
			}
			reloadConnSelect(w, "")
		}, win)

	}
	saveBtn.OnTapped = func() {
		conn, _, _ := mySQLConnectionModel.GetByName(connEntry.Selected)
		conn.Host = strings.TrimSpace(hostEntry.Text)
		conn.Port = cast.ToInt(strings.TrimSpace(portEntry.Text))
		conn.UserName = strings.TrimSpace(userEntry.Text)
		conn.Password = strings.TrimSpace(passwordEntry.Text)
		_, err := models.MySQLConnManager.LoadConnection(conn)
		if err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("连接输出库出错,error:%s", err.Error()), win)
			return
		}

		if err = mySQLConnectionModel.Update(&conn); err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("更新连接出错,error:%s", err.Error()), win)
			return
		}
		dialog.ShowInformation("成功", "保存成功", win)
	}
	saveBtn.Importance = widget.HighImportance
	form.Append("", container.NewBorder(nil, nil, nil, container.NewHBox(deleteBtn, saveBtn)))
	return form
}

type MySQLConnWidget struct {
	connEntry     *widget.Select
	hostEntry     *widget.Entry
	portEntry     *widget.Entry
	userEntry     *widget.Entry
	passwordEntry *widget.Entry
	saveBtn       *widget.Button
	resetBtn      *widget.Button
}

func reloadConnSelect(connWidget MySQLConnWidget, name string) {

	connSlice := GetMySQLConnSlice()

	if strutils.IsEmptyString(name) && len(connSlice) != 0 {
		name = connSlice[0].Name
	}

	names := connSlice.GetNames()
	connEntry := connWidget.connEntry
	hostEntry := connWidget.hostEntry
	portEntry := connWidget.portEntry
	userEntry := connWidget.userEntry
	saveBtn := connWidget.saveBtn
	resetBtn := connWidget.resetBtn
	passwordEntry := connWidget.passwordEntry
	connEntry.Options = names

	if len(connSlice) == 0 {
		connEntry.Selected = "暂无连接"
		saveBtn.Disable()
		resetBtn.Disable()
	} else {
		saveBtn.Enable()
		resetBtn.Enable()
	}
	if len(names) != 0 && strutils.IsNotEmptyString(name) && arrays.ContainsString(names, name) != -1 {
		connEntry.Selected = name
	}
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
	saveBtn.Refresh()
	resetBtn.Refresh()

}
func MySQLAddScene(win fyne.Window) fyne.CanvasObject {
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
		_, exists, err := mySQLConnectionModel.GetByName(name)
		if err != nil {
			dialog.ShowInformation("错误", "判断连接名是否存在出错,err:%s"+err.Error(), win)
			return
		}
		if exists {
			dialog.ShowInformation("错误", "连接名已存在", win)
			return
		}
		conn := ddl.MySQLConnection{
			Name:     name,
			Host:     host,
			Port:     port,
			UserName: user,
			Password: pwd,
		}
		err = mySQLConnectionModel.Insert(&conn)
		if err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("插入链接出错,error:%s", err.Error()), win)
			return
		}
		dialog.ShowInformation("成功", "保存成功", win)

	})
	saveBtn.Importance = widget.HighImportance
	form.Append("", container.NewBorder(nil, nil, nil, container.NewHBox(resetBtn, saveBtn)))
	return form
}
