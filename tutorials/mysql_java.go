package tutorials

import (
	"code-gen/settings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeEntityScene() fyne.CanvasObject {
	global := settings.GetGlobal()
	lombok := global.GetLombok()
	getter := &widget.Check{
		Text:    "@Getter",
		Checked: lombok.Getter,
		OnChanged: func(b bool) {
			global.ChangeLombokGetter(b)
		},
	}
	setter := &widget.Check{
		Text:    "@Setter",
		Checked: lombok.Setter,
		OnChanged: func(b bool) {
			global.ChangeLombokSetter(b)
		},
	}
	data := &widget.Check{
		Text:    "@Data",
		Checked: lombok.Data,
		OnChanged: func(b bool) {
			global.ChangLombokData(b)
		},
	}

	slf4j := &widget.Check{
		Text:    "@Slf4j",
		Checked: lombok.Slf4j,
		OnChanged: func(b bool) {
			global.ChangeLombokSlf4j(b)
		},
	}

	noArgsConstructor := &widget.Check{
		Text:    "@NoArgsConstructor",
		Checked: lombok.NoArgsConstructor,
		OnChanged: func(b bool) {
			global.ChangeNoArgsConstructor(b)
		},
	}

	allArgsConstructor := &widget.Check{
		Text:    "@AllArgsConstructor",
		Checked: lombok.AllArgsConstructor,
		OnChanged: func(b bool) {
			global.ChangeAllArgsConstructor(b)
		},
	}
	toString := &widget.Check{
		Text:    "@ToString",
		Checked: lombok.ToString,
		OnChanged: func(b bool) {
			global.ChangeLombokToString(b)
		},
	}
	equalsAndHashCode := &widget.Check{
		Text:    "@EqualsAndHashCode",
		Checked: lombok.EqualsAndHashCode,
		OnChanged: func(b bool) {
			global.ChangeEqualsAndHashCode(b)
		},
	}
	lombokCheck := container.NewGridWithColumns(
		3,
		getter,
		setter,
		data,
		slf4j,
		noArgsConstructor,
		allArgsConstructor,
		toString,
		equalsAndHashCode,
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
		"insert dup",
		"replace one",
		"replace many",
		"delete",
		"update",
		"select count",
		"select one",
		"select many",
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
			PlaceHolder: "entity",
		},
		genButton,
	)
	right := container.NewVScroll(rightBox)
	selectChecks := widget.NewCheckGroup([]string{"=", "!=", ">", "<", ">=", "<=", "in"}, func(strings []string) {

	})
	selectChecks.Horizontal = true
	insertUpdateChecks := widget.NewCheckGroup([]string{
		"插入", "更新", "返回",
	}, func(strings []string) {

	})
	insertUpdateChecks.Horizontal = true
	form := widget.NewForm(
		&widget.FormItem{
			Text:   "字段",
			Widget: insertUpdateChecks,
		},
		&widget.FormItem{
			Text:   "条件",
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
