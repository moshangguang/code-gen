package tutorials

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeEntityScene() fyne.CanvasObject {
	lombokCheck := container.NewGridWithColumns(
		3,
		widget.NewCheck("@Getter", func(b bool) {

		}),
		widget.NewCheck("@Setter", func(b bool) {

		}),
		widget.NewCheck("@Data", func(b bool) {

		}),
		widget.NewCheck("@Slf4j", func(b bool) {

		}),
		widget.NewCheck("@NoArgsConstructor", func(b bool) {

		}),
		widget.NewCheck("@AllArgsConstructor", func(b bool) {

		}),
	)
	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "lombok",
				Widget: lombokCheck,
			},
		},
		OnSubmit: func() {

		},
		SubmitText: "生成",
	}
	return form
}

func makeMyBatisScene() fyne.CanvasObject {
	sqlSelect := widget.NewSelect([]string{
		"insert incr",
		"insert one",
		"insert many",
		"insert duplicate",
		"replace one",
		"replace many",
		"delete",
		"delete by id",
		"delete by ids",
		"update",
		"update by id",
		"update by ids",
		"select",
	}, func(s string) {

	})
	sqlSelect.SetSelectedIndex(0)
	genButton := &widget.Button{
		Text:       "生成",
		Icon:       theme.ConfirmIcon(),
		Importance: widget.HighImportance,
		OnTapped: func() {

		},
	}
	rightBox := container.NewVBox(
		sqlSelect,
		&widget.Entry{
			PlaceHolder: "method",
		},
		&widget.Entry{
			PlaceHolder: "primary key",
		},
		&widget.Button{
			Text: "取消勾选",
			Icon: theme.CancelIcon(),
			OnTapped: func() {

			},
		},
		genButton,
	)
	right := container.NewVScroll(rightBox)
	selectChecks := widget.NewCheckGroup([]string{"=", "!=", ">", "<", ">=", "<="}, func(strings []string) {

	})
	selectChecks.Horizontal = true
	insertUpdateChecks := widget.NewCheckGroup([]string{
		"插入", "更新",
	}, func(strings []string) {

	})
	insertUpdateChecks.Horizontal = true
	form := widget.NewForm(
		&widget.FormItem{
			Text:   "插入/更新项",
			Widget: insertUpdateChecks,
		},
		&widget.FormItem{
			Text:   "查询项",
			Widget: selectChecks,
		})
	ac := widget.NewAccordion(
		widget.NewAccordionItem("B", widget.NewLabel("Two")),
		&widget.AccordionItem{
			Title:  "C",
			Detail: form,
			Open:   true,
		},
	)
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	ac.Append(widget.NewAccordionItem("E", &widget.Check{Text: "E"}))
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	ac.Append(widget.NewAccordionItem("E", &widget.Check{Text: "E"}))
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	ac.Append(widget.NewAccordionItem("E", &widget.Check{Text: "E"}))
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	ac.Append(widget.NewAccordionItem("E", &widget.Check{Text: "E"}))
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	ac.Append(widget.NewAccordionItem("E", &widget.Check{Text: "E"}))
	ac.MultiOpen = true
	return container.NewBorder(nil, nil, nil, right, container.NewVScroll(container.NewVBox(ac)))

}
func MySQLJavaScene(_ fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Entity", makeEntityScene()),
		container.NewTabItem("MyBatis", makeMyBatisScene()),
	)
	form := widget.NewForm(
		&widget.FormItem{
			Text: "连接",
			Widget: widget.NewSelect([]string{}, func(s string) {

			}),
		},
		&widget.FormItem{
			Text:   "数据库",
			Widget: widget.NewSelectEntry([]string{}),
		},
		&widget.FormItem{
			Text:   "表",
			Widget: widget.NewSelectEntry([]string{}),
		},
	)
	//locations := makeTabLocationSelect(tabs.SetTabLocation)
	//border := container.NewBorder(locations, nil, nil, nil, tabs)
	return container.NewBorder(form, nil, nil, nil, tabs)
}
