package mysql

import (
	"code-gen/settings/lombok"
	"code-gen/settings/mysql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Handler struct {
	GetConn     func() (connName string, ok bool)
	GetDatabase func() (dbName string, ok bool)
	GetTable    func() (tableName string, ok bool)
	Callback    []func()
}

func (handler *Handler) CallbackAll() {
	if len(handler.Callback) == 0 {
		return
	}
	for _, fn := range handler.Callback {
		fn()
	}
}
func (handler *Handler) RegisterCallback(fn func()) {
	if handler.Callback == nil {
		handler.Callback = make([]func(), 0)
	}
	handler.Callback = append(handler.Callback, fn)
}
func makeEntityScene(_ *Handler) fyne.CanvasObject {
	lom := lombok.GetLombok()
	getter := &widget.Check{
		Text:    "@Getter",
		Checked: lom.Getter,
		OnChanged: func(b bool) {
			lom.Getter = b
			lom.Save()
		},
	}
	setter := &widget.Check{
		Text:    "@Setter",
		Checked: lom.Setter,
		OnChanged: func(b bool) {
			lom.Setter = b
			lom.Save()
		},
	}
	data := &widget.Check{
		Text:    "@Data",
		Checked: lom.Data,
		OnChanged: func(b bool) {
			lom.Data = b
			lom.Save()
		},
	}

	slf4j := &widget.Check{
		Text:    "@Slf4j",
		Checked: lom.Slf4j,
		OnChanged: func(b bool) {
			lom.Slf4j = b
			lom.Save()
		},
	}

	noArgsConstructor := &widget.Check{
		Text:    "@NoArgsConstructor",
		Checked: lom.NoArgsConstructor,
		OnChanged: func(b bool) {
			lom.NoArgsConstructor = b
			lom.Save()
		},
	}

	allArgsConstructor := &widget.Check{
		Text:    "@AllArgsConstructor",
		Checked: lom.AllArgsConstructor,
		OnChanged: func(b bool) {
			lom.AllArgsConstructor = b
			lom.Save()
		},
	}
	toString := &widget.Check{
		Text:    "@ToString",
		Checked: lom.ToString,
		OnChanged: func(b bool) {
			lom.ToString = b
			lom.Save()
		},
	}
	equalsAndHashCode := &widget.Check{
		Text:    "@EqualsAndHashCode",
		Checked: lom.EqualsAndHashCode,
		OnChanged: func(b bool) {
			lom.EqualsAndHashCode = b
			lom.Save()
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

func makeMyBatisScene(handler *Handler) fyne.CanvasObject {
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
	dictionary := mysql.GetDictionary()
	handler.RegisterCallback(func() {
		connName, ok := handler.GetConn()
		if !ok {
			return
		}
		dbName, ok := handler.GetDatabase()
		if !ok {
			return
		}
		table, ok := handler.GetTable()
		if !ok {
			return
		}
		conn, ok := dictionary.Get(connName)
		if !ok {
			return
		}
		conn = conn.Use(dbName)
		fields := conn.GetFields(table)
		if len(fields) == 0 {
			return
		}

		head := &widget.AccordionItem{
			Title:  fmt.Sprintf("表名:%s", table),
			Detail: newMyBatisItemForm(),
			Open:   true,
		}
		items := make([]*widget.AccordionItem, 0, len(fields)+1)
		items = append(items, head)
		for _, field := range fields {
			items = append(items, &widget.AccordionItem{
				Title:  fmt.Sprintf("字段名:%s", field.ColumnName),
				Detail: newMyBatisItemForm(),
				Open:   true,
			})
		}
		ac.Items = items
		ac.MultiOpen = true
		ac.Refresh()
	})
	return container.NewBorder(nil, nil, nil, right, container.NewVScroll(container.NewVBox(ac)))
}

func newMyBatisItemForm() *widget.Form {
	selectChecks := widget.NewCheckGroup([]string{"=", "!=", ">", "<", ">=", "<=", "in"}, func(strings []string) {

	})
	selectChecks.Horizontal = true
	insertUpdateChecks := widget.NewCheckGroup([]string{
		"插入", "更新", "返回",
	}, func(strings []string) {

	})
	insertUpdateChecks.Horizontal = true
	return widget.NewForm(
		&widget.FormItem{
			Text:   "字段",
			Widget: insertUpdateChecks,
		},
		&widget.FormItem{
			Text:   "条件",
			Widget: selectChecks,
		})
}

const MaxLimit = 100

func JavaScene(win fyne.Window) fyne.CanvasObject {

	dictionary := mysql.GetDictionary()
	mysqlConn := dictionary.All().GetNames()
	connSelect := &widget.Select{
		Options: mysqlConn,
	}
	dbSelect := &widget.Select{}
	tableSelect := &widget.Select{}
	getConn := func() (conn string, ok bool) {
		conn = connSelect.Selected
		ok = len(conn) != 0
		return
	}
	getDatabase := func() (db string, ok bool) {
		db = dbSelect.Selected
		ok = len(db) != 0
		return
	}
	getTable := func() (table string, ok bool) {
		table = tableSelect.Selected
		ok = len(table) != 0
		return
	}
	handler := &Handler{
		GetConn:     getConn,
		GetDatabase: getDatabase,
		GetTable:    getTable,
		Callback:    []func(){},
	}
	tabs := container.NewAppTabs(
		container.NewTabItem("Entity", makeEntityScene(handler)),
		container.NewTabItem("MyBatis", makeMyBatisScene(handler)),
	)
	connSelect.OnChanged = func(s string) {
		conn, ok := dictionary.Get(s)
		if !ok {
			return
		}
		dbSelect.Selected = ""
		databases := conn.PatternDatabases("", MaxLimit)
		dbSelect.Options = databases
		dbSelect.Refresh()
	}

	dbSelect.OnChanged = func(s string) {
		connName, ok := getConn()
		if !ok {
			return
		}
		conn, ok := dictionary.Get(connName)
		if !ok {
			return
		}
		conn = conn.Use(s)
		tables := conn.PatternTables("", MaxLimit)
		tableSelect.Options = tables
		tableSelect.Refresh()
	}
	tableSelect.OnChanged = func(s string) {
		handler.CallbackAll()
	}
	dbSearch := &widget.Entry{
		PlaceHolder: "模糊搜索数据库名",
	}
	tableSearch := &widget.Entry{
		PlaceHolder: "模糊搜索表名",
	}
	form := widget.NewForm(
		&widget.FormItem{
			Text:   "连接",
			Widget: connSelect,
		},
		&widget.FormItem{
			Text: "数据库",
			Widget: container.NewBorder(dbSelect, nil, nil, container.NewGridWithColumns(2, widget.NewButton("清空", func() {
				dbSearch.Text = ""
				dbSearch.Refresh()
			}), widget.NewButton("搜索", func() {
				connName, ok := getConn()
				if !ok {
					return
				}
				conn, ok := dictionary.Get(connName)
				if !ok {
					return
				}
				databases := conn.PatternDatabases(dbSearch.Text, MaxLimit)
				dbSelect.Options = databases
				dbSelect.Refresh()
				tableSelect.Options = []string{}
				tableSelect.Refresh()
			})), dbSearch),
		},

		&widget.FormItem{
			Text: "表",
			Widget: container.NewBorder(tableSelect, nil, nil, container.NewGridWithColumns(2, widget.NewButton("清空", func() {
				tableSearch.Text = ""
				tableSearch.Refresh()
			}), widget.NewButton("搜索", func() {
				connName, ok := getConn()
				if !ok {
					return
				}
				conn, ok := dictionary.Get(connName)
				if !ok {
					return
				}
				dbName, ok := getDatabase()
				if !ok {
					return
				}
				conn = conn.Use(dbName)
				tables := conn.PatternTables(tableSearch.Text, MaxLimit)
				tableSelect.Selected = ""
				tableSelect.Options = tables
				tableSelect.Refresh()
			})), tableSearch),
		},
	)

	return container.NewBorder(form, nil, nil, nil, tabs)
}
