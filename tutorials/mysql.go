package tutorials

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func MySQLScene(_ fyne.Window) fyne.CanvasObject {
	connEntry := widget.NewSelectEntry([]string{"Option A", "Option B", "Option C"})
	connEntry.PlaceHolder = "请输入或选择连接名"
	connEntry.OnChanged = func(s string) {
		fmt.Println(s)
	}
	host := widget.NewEntry()
	port := widget.NewEntry()
	port.SetText("3306")
	port.SetPlaceHolder("3306")
	user := widget.NewEntry()
	user.SetText("root")
	user.SetPlaceHolder("root")
	password := widget.NewEntry()
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
				Widget: host,
			},
			{
				Text:   "端口",
				Widget: port,
			},
			{
				Text:   "用户名",
				Widget: user,
			},
			//{
			//	Text:   "密码",
			//	Widget: password,
			//},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		CancelText: "删除",
		OnSubmit: func() {
		},
		SubmitText: "保存",
	}
	form.Append("密码", password)
	return form
}
