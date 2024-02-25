package tutorials

import (
	"code-gen/constants"
	"code-gen/help"
	"code-gen/pkg/db"
	"code-gen/pkg/models"
	"code-gen/pkg/models/ddl"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/wxnacy/wgo/arrays"
)

func GolangScene(win fyne.Window) fyne.CanvasObject {
	form := widget.NewForm()
	setting := ddl.GolangSetting{}
	db.GetDatabase().Unmarshal(constants.GolangSetting, &setting)
	packageEntry := widget.NewEntry()
	packageEntry.Text = setting.PackageName
	packageEntry.OnChanged = func(s string) {
		gs := ddl.GolangSetting{}
		db.GetDatabase().Unmarshal(constants.GolangSetting, &gs)
		gs.PackageName = s
		db.GetDatabase().Save(constants.GolangSetting, gs)
	}
	connSlice := GetMySQLConnSlice()
	names := connSlice.GetNames()
	tableSelect := widget.NewSelectEntry([]string{})
	connEntry := widget.NewSelect(names, func(s string) {

	})
	dbSelect := widget.NewSelect([]string{}, func(s string) {
		gs := ddl.GolangSetting{}
		db.GetDatabase().Unmarshal(constants.GolangSetting, &gs)
		gs.DatabaseName = s
		db.GetDatabase().Save(constants.GolangSetting, gs)
		if conn, ok := GetMySQLConnSlice().First(func(connection ddl.MySQLConnection) bool {
			return connection.Name == connEntry.Selected
		}); ok {
			conn.Database = s
			tables, err := GetTables(conn)
			if err != nil {
				//todo:log
			}
			if len(tables) != 0 {
				tableSelect.SetOptions(tables)
				tableSelect.SetText(tables[0])
			}
		}
	})
	connEntry.OnChanged = func(s string) {
		gs := ddl.GolangSetting{}
		db.GetDatabase().Unmarshal(constants.GolangSetting, &gs)
		gs.ConnName = s
		db.GetDatabase().Save(constants.GolangSetting, gs)
	}

	tableSelect.OnSubmitted = func(s string) {

	}
	if len(names) != 0 {
		connEntry.Selected = names[0]
	}
	if help.IsNotEmptyString(setting.ConnName) && arrays.ContainsString(names, setting.ConnName) != -1 {
		connEntry.Selected = setting.ConnName
	}
	if help.IsNotEmptyString(connEntry.Selected) {
		if connection, ok := connSlice.First(func(connection ddl.MySQLConnection) bool {
			return connection.Name == connEntry.Selected
		}); ok {
			databases, err := GetDatabases(connection)
			if err != nil {
				//todo:log
			} else if len(databases) != 0 {
				dbSelect.Options = databases
				dbSelect.Selected = databases[0]
				if help.IsNotEmptyString(setting.DatabaseName) && arrays.ContainsString(databases, setting.DatabaseName) != -1 {
					dbSelect.Selected = setting.DatabaseName
				}
			}
			connection.Database = dbSelect.Selected
			tables, err := GetTables(connection)
			if err != nil {
				//todo:log
			} else if len(tables) != 0 {
				tableSelect.SetOptions(tables)
				tableSelect.SetText(tables[0])
			}
		}

	}
	form.Append("包名:", packageEntry)
	form.Append("连接名:", connEntry)

	form.Append("数据库:", dbSelect)
	form.Append("表名:", tableSelect)
	form.SubmitText = "生成"
	form.OnSubmit = func() {

	}
	return form
}
func GetDatabases(connection ddl.MySQLConnection) ([]string, error) {
	engine, err := models.MySQLConnManager.LoadConnection(connection)
	if err != nil {
		return make([]string, 0), err
	}
	queryString, err := engine.SQL("show databases").QueryString()
	if err != nil {
		return make([]string, 0), err
	}
	result := make([]string, 0, len(queryString))
	for _, s := range queryString {
		result = append(result, s["Database"])
	}
	return result, nil
}
func saveGolangSetting(func(setting ddl.GolangSetting) ddl.GolangSetting) {

}

func GetMySQLConnSlice() ddl.MySQLConnectionSlice {
	connSlice := make(ddl.MySQLConnectionSlice, 0)
	db.GetDatabase().Unmarshal(constants.MySQLConnection, &connSlice)
	return connSlice
}

func GetTables(connection ddl.MySQLConnection) ([]string, error) {
	engine, err := models.MySQLConnManager.LoadConnection(connection)
	if err != nil {
		return make([]string, 0), err
	}
	queryString, err := engine.SQL("SELECT table_name table_name FROM information_schema.tables WHERE table_schema = ?", connection.Database).QueryString()
	if err != nil {
		return make([]string, 0), err
	}
	result := make([]string, 0)
	for _, d := range queryString {
		result = append(result, d["table_name"])
	}
	return result, nil
}
