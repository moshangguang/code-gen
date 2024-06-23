package tutorials

import (
	"bytes"
	"code-gen/constants"
	"code-gen/pkg/log"
	"code-gen/pkg/models"
	"code-gen/pkg/models/ddl"
	"code-gen/pkg/types"
	"code-gen/utils/fileUtils"
	"code-gen/utils/strutils"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"xorm.io/core"
)

var mapper = &core.SnakeMapper{}
var codeTemplate = template.New("code_template")

func init() {
	codeTemplate.Funcs(template.FuncMap{"Mapper": mapper.Table2Obj,
		"Type":       typeString,
		"Tag":        tag,
		"UnTitle":    unTitle,
		"gt":         gt,
		"getCol":     getCol,
		"UpperTitle": upTitle,
	})
	var err error
	codeTemplate, err = codeTemplate.Parse(`package {{.Models}}

{{$ilen := len .Imports}}
{{if gt $ilen 0}}
import (
	{{range .Imports}}"{{.}}"{{end}}
)
{{end}}

{{range .Tables}}
type {{Mapper .Name}} struct {
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper $col.Name}}	{{Type $col}} {{Tag $table $col}}
{{end}}
}

func (m {{Mapper .Name}}) TableName() string {
	return "{{ .Name}}"
}


{{end}}`)
	if err != nil {
		panic(err)
	}
}

func XormScene(win fyne.Window) fyne.CanvasObject {
	form := widget.NewForm()
	setting := GetGolangSetting()
	packageEntry := widget.NewEntry()
	packageEntry.Text = setting.PackageName
	packageEntry.OnChanged = func(s string) {
		SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
			setting.PackageName = s
			return setting
		})
	}
	connSlice, _ := mySQLConnectionModel.GetAll()
	names := connSlice.GetNames()
	tableSelect := widget.NewSelectEntry([]string{})
	connEntry := widget.NewSelect(names, func(s string) {

	})
	fileNameEntry := widget.NewEntry()
	fileNameEntry.Text = setting.FileName
	fileNameEntry.OnChanged = func(s string) {
		SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
			setting.FileName = s
			return setting
		})
	}
	dbSelect := widget.NewSelect([]string{}, func(s string) {

	})
	w := ReloadConnWidget{
		connSelect:  connEntry,
		dbSelect:    dbSelect,
		tableSelect: tableSelect,
		win:         win,
	}
	dbSelect.OnChanged = func(s string) {
		SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
			if setting.ConnConfig == nil {
				setting.ConnConfig = make(map[string]types.ConnConfig)
			}
			setting.ConnName = connEntry.Selected
			conf := setting.ConnConfig[connEntry.Selected]
			conf.Database = s
			setting.ConnConfig[connEntry.Selected] = conf
			reloadTableSelect(w, conf, "")
			return setting
		})
	}
	connEntry.OnChanged = func(s string) {
		SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
			setting.ConnName = s
			conf := setting.ConnConfig[s]
			reloadDatabaseSelect(w, conf)
			reloadTableSelect(w, conf, "")
			return setting
		})
	}

	tableSelect.OnChanged = func(s string) {
		SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
			if setting.ConnConfig == nil {
				setting.ConnConfig = make(map[string]types.ConnConfig)
			}
			setting.ConnName = connEntry.Selected
			conf := setting.ConnConfig[connEntry.Selected]
			conf.Database = dbSelect.Selected
			conf.Table = s
			setting.ConnConfig[connEntry.Selected] = conf
			return setting
		})
	}
	tableSelect.OnSubmitted = func(s string) {
		SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
			if setting.ConnConfig == nil {
				setting.ConnConfig = make(map[string]types.ConnConfig)
			}
			conf := setting.ConnConfig[connEntry.Selected]
			conf.Table = s
			setting.ConnConfig[connEntry.Selected] = conf
			reloadTableSelect(w, conf, strings.TrimSpace(s))
			return setting
		})

	}
	if len(names) != 0 {
		connEntry.Selected = names[0]
	}
	if strutils.IsNotEmptyString(setting.ConnName) && arrays.ContainsString(names, setting.ConnName) != -1 {
		connEntry.Selected = setting.ConnName
	}
	if strutils.IsNotEmptyString(connEntry.Selected) {
		if _, ok := connSlice.First(func(connection ddl.MySQLConnection) bool {
			return connection.Name == connEntry.Selected
		}); ok {
			conf := setting.ConnConfig[setting.ConnName]
			reloadDatabaseSelect(w, conf)
			reloadTableSelect(w, conf, "")
		}

	}
	form.Append("包名:", packageEntry)
	form.Append("连接名:", connEntry)
	form.Append("数据库:", dbSelect)
	form.AppendItem(&widget.FormItem{
		Text:     "表名:",
		Widget:   tableSelect,
		HintText: "输入后回车即可模糊搜索",
	})
	form.AppendItem(&widget.FormItem{
		Text:     "文件名:",
		Widget:   fileNameEntry,
		HintText: "默认以表名作为文件名",
	})
	folderLabel := widget.NewLabel("")
	if strutils.IsNotEmptyString(setting.OutputPath) && fileUtils.IsDir(setting.OutputPath) {
		folderLabel.SetText(setting.OutputPath)
	}
	openFolder := widget.NewButtonWithIcon("选择目录", theme.FolderOpenIcon(), func() {
		dlg := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}
			IsDir := fileUtils.IsDir(list.Path())
			if !IsDir {
				return
			}
			folderLabel.SetText(list.Path())
			SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
				setting.OutputPath = list.Path()
				return setting
			})
		}, win)
		gs := GetGolangSetting()
		if strutils.IsNotEmptyString(gs.OutputPath) && fileUtils.IsDir(gs.OutputPath) {
			if forURI, err := storage.ListerForURI(storage.NewFileURI(gs.OutputPath)); err == nil {
				dlg.SetLocation(forURI)
			}
		}
		dlg.Resize(fyne.NewSize(800, 500))
		dlg.SetConfirmText("选择")
		dlg.SetDismissText("关闭")

		dlg.Show()
	})
	clearButton := widget.NewButtonWithIcon("清空", theme.ContentClearIcon(), func() {
		folderLabel.SetText("")
		SaveGolangSettingWithClosure(func(setting types.GolangSetting) types.GolangSetting {
			setting.OutputPath = ""
			return setting
		})
	})
	form.AppendItem(&widget.FormItem{
		Text:     "输出目录:",
		Widget:   container.NewBorder(nil, nil, folderLabel, container.NewHBox(openFolder, clearButton)),
		HintText: "默认输出文件在生成器所在目录",
	})
	form.SubmitText = "生成"

	form.OnSubmit = func() {
		if len(connSlice) == 0 {
			dialog.ShowInformation("错误", "获取不到连接", win)
			return
		}
		conn, ok := connSlice.First(func(connection ddl.MySQLConnection) bool {
			return connection.Name == connEntry.Selected
		})
		if !ok {
			return
		}
		conn.Database = dbSelect.Selected
		engine, err := models.MySQLConnManager.LoadConnection(conn)
		if err != nil {
			dialog.ShowInformation("错误", fmt.Sprintf("获取连接出错,error:%s", err.Error()), win)
			return
		}
		newBytes := bytes.NewBufferString("")

		tableName := strings.TrimSpace(tableSelect.Text)
		table := core.NewTable(tableName, nil)
		tables, err := engine.Dialect().GetTables()
		if err != nil {
			dialog.ShowInformation("错误", "获取所有表出错,err:"+err.Error(), win)
			return
		}
		hasTable := false
		for _, t := range tables {
			if t.Name != tableName {
				continue
			}
			hasTable = true
			break
		}
		if !hasTable {
			dialog.ShowInformation("错误", "表不存在", win)
			return
		}
		_, columns, err := engine.Dialect().GetColumns(tableName)
		queryString, err := engine.SQL(fmt.Sprintf("SHOW COLUMNS FROM `%s`", tableName)).QueryString()
		if err != nil {
			dialog.ShowInformation("错误", "获取表字段出错,err:%s"+err.Error(), win)
			return
		}
		for _, d := range queryString {
			field, ok := d["Field"]
			if !ok {
				continue
			}
			column, ok := columns[field]
			if !ok {
				continue
			}
			table.AddColumn(column)
		}
		gs := GetGolangSetting()
		tables = []*core.Table{table}
		tmpl := &Tmpl{
			Tables:  tables,
			Imports: genGoImports(tables),
			Models:  gs.PackageName,
		}
		err = codeTemplate.Execute(newBytes, tmpl)
		if err != nil {
			dialog.ShowInformation("错误", "生成文件内容出错,err:"+err.Error(), win)
			return
		}
		execPath, err := os.Executable()
		if err != nil {
			dialog.ShowInformation("错误", "获取执行目录出错,err:"+err.Error(), win)
			return
		}
		gs = GetGolangSetting()
		var fileDir string
		if strutils.IsNotEmptyString(gs.OutputPath) && fileUtils.IsDir(gs.OutputPath) {
			fileDir = gs.OutputPath
		}
		if strutils.IsEmptyString(fileDir) {
			fileDir = filepath.Dir(execPath)
		}
		fileName := gs.FileName
		index := strings.Index(fileName, ".")
		if index == 0 {
			fileName = ""
		}
		if index > 0 {
			fileName = fileName[:index]
		}
		if strutils.IsEmptyString(fileName) {
			fileName = tableName
		}
		allBytes, err := io.ReadAll(newBytes)
		if err != nil {
			return
		}
		for i := 0; i < math.MaxInt32; i++ {
			var filePath string
			if i == 0 {
				filePath = filepath.Join(fileDir, strings.Join([]string{fileName, "go"}, "."))
			} else {
				filePath = filepath.Join(fileDir, strings.Join([]string{fmt.Sprintf("%s(%d)", fileName, i), "go"}, "."))
			}
			exist := fileUtils.IsExist(filePath)
			if exist {
				continue
			}
			if err = fileUtils.CreateAndWrite(filePath, allBytes); err != nil {
				dialog.ShowInformation("错误", "创建和写入文件出错,err:"+err.Error(), win)
				return
			}
			cmd := exec.Command("go", "fmt", filePath)
			_ = cmd.Run()
			dialog.ShowInformation("提示", "文件生成成功!", win)
			break
		}
	}
	return form
}

type Tmpl struct {
	Tables  []*core.Table
	Imports map[string]string
	Models  string
}

func getCol(cols map[string]*core.Column, name string) *core.Column {
	return cols[strings.ToLower(name)]
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

func upTitle(src string) string {
	if src == "" {
		return ""
	}

	return strings.ToUpper(src)
}

func GetMySQLConnSlice() ddl.MySQLConnectionSlice {
	slice, err := mySQLConnectionModel.GetAll()
	if err != nil {
		log.Logger.Error("获取所有连接出错", zap.Error(err))
	}
	return slice
}

func GetTables(connection ddl.MySQLConnection, name string) ([]string, error) {
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
		t := d["table_name"]
		if strutils.IsNotEmptyString(name) && !strings.Contains(t, name) {
			continue
		}
		result = append(result, d["table_name"])
	}
	return result, nil
}

type ReloadConnWidget struct {
	connSelect  *widget.Select
	dbSelect    *widget.Select
	tableSelect *widget.SelectEntry
	win         fyne.Window
}

func reloadDatabaseSelect(connWidget ReloadConnWidget, conf types.ConnConfig) {
	connEntry := connWidget.connSelect
	conn, _, _ := mySQLConnectionModel.GetByName(connEntry.Selected)
	databases, err := GetDatabases(conn)
	if err != nil {
		return
	}
	connWidget.dbSelect.Options = databases
	if len(databases) != 0 {
		connWidget.dbSelect.Selected = databases[0]
	}
	if strutils.IsNotEmptyString(conf.Database) && arrays.ContainsString(databases, conf.Database) != -1 {
		connWidget.dbSelect.Selected = conf.Database
	}
	connWidget.dbSelect.Refresh()
}

func reloadTableSelect(connWidget ReloadConnWidget, conf types.ConnConfig, name string) {
	conn, exists, err := mySQLConnectionModel.GetByName(connWidget.connSelect.Selected)
	if err != nil {
		log.Logger.Error("刷新表选项获取连接出错", zap.Error(err))
		return
	}
	if !exists {
		log.Logger.Error("刷新表选项获取不到连接")
		return
	}
	conn.Database = connWidget.dbSelect.Selected
	tables, err := GetTables(conn, name)
	if err != nil {
		dialog.ShowInformation("错误", fmt.Sprintf("获取表出错,error:%s", err.Error()), connWidget.win)
		return
	}
	connWidget.tableSelect.SetOptions(tables)
	if len(tables) != 0 {
		connWidget.tableSelect.SetText(tables[0])
	}
	if strutils.IsNotEmptyString(conf.Table) && arrays.ContainsString(tables, conf.Table) != -1 {
		connWidget.tableSelect.SetText(conf.Table)
	}
}

func GetGolangSetting() types.GolangSetting {
	setting := types.GolangSetting{
		ConnConfig: map[string]types.ConnConfig{},
	}
	config, exists, err := configModel.GetByName(constants.GolangSetting)
	if err != nil {
		return setting
	}
	if !exists {
		return setting
	}
	err = json.Unmarshal([]byte(config.Value), &setting)
	return setting
}

func SaveGolangSettingWithClosure(fn func(types.GolangSetting) types.GolangSetting) {
	setting := GetGolangSetting()
	setting = fn(setting)
	settingBytes, err := json.Marshal(setting)
	if err != nil {
		log.Logger.Error("序列化golang配置出错", zap.Error(err))
		return
	}
	err = configModel.Save(constants.GolangSetting, string(settingBytes))
	if err != nil {
		log.Logger.Error("保存golang配置出错", zap.Error(err))
		return
	}
}

func typeString(col *core.Column) string {
	st := col.SQLType
	name := strings.ToUpper(st.Name)
	switch name {
	case core.TinyInt:
		if strings.HasPrefix(col.Name, "is_") || strings.HasPrefix(col.Name, "has_") {
			return reflect.TypeOf(true).String()
		}
		return reflect.TypeOf(int8(1)).String()
	case core.SmallInt:
		return reflect.TypeOf(int16(1)).String()
	}
	t := core.SQLType2Type(st)
	s := t.String()
	if s == "[]uint8" {
		return "[]byte"
	}
	return s
}
func tag(table *core.Table, col *core.Column) string {
	isNameId := (mapper.Table2Obj(col.Name) == "Id")
	isIdPk := isNameId && typeString(col) == "int64"

	var res []string
	if !col.Nullable {
		if !isIdPk {
			res = append(res, "not null")
		}
	}
	if col.IsPrimaryKey {
		res = append(res, "pk")
	}
	if col.Default != "" {
		res = append(res, "default "+col.Default)
	}
	if col.IsAutoIncrement {
		res = append(res, "autoincr")
	}
	isCreated := col.Name == "created_at"
	if col.SQLType.IsTime() && isCreated {
		res = append(res, "created")
	}

	isUpdated := col.Name == "updated_at"
	if col.SQLType.IsTime() && isUpdated {
		res = append(res, "updated")
	}

	isDeleted := col.Name == "deleted_at"
	if col.SQLType.IsTime() && isDeleted {
		res = append(res, "deleted")
	}

	names := make([]string, 0, len(col.Indexes))
	for name := range col.Indexes {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		index := table.Indexes[name]
		var uistr string
		if index.Type == core.UniqueType {
			uistr = "unique"
		} else if index.Type == core.IndexType {
			uistr = "index"
		}
		if len(index.Cols) > 1 {
			uistr += "(" + index.Name + ")"
		}
		res = append(res, uistr)
	}

	nstr := strings.ToLower(col.SQLType.Name)
	if col.Length != 0 {
		if col.Length2 != 0 {
			nstr += fmt.Sprintf("(%v,%v)", col.Length, col.Length2)
		} else {
			nstr += fmt.Sprintf("(%v)", col.Length)
		}
	} else if len(col.EnumOptions) > 0 { //enum
		nstr += "("
		opts := ""

		enumOptions := make([]string, 0, len(col.EnumOptions))
		for enumOption := range col.EnumOptions {
			enumOptions = append(enumOptions, enumOption)
		}
		sort.Strings(enumOptions)

		for _, v := range enumOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nstr += strings.TrimLeft(opts, ",")
		nstr += ")"
	} else if len(col.SetOptions) > 0 { //enum
		nstr += "("
		opts := ""

		setOptions := make([]string, 0, len(col.SetOptions))
		for setOption := range col.SetOptions {
			setOptions = append(setOptions, setOption)
		}
		sort.Strings(setOptions)

		for _, v := range setOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nstr += strings.TrimLeft(opts, ",")
		nstr += ")"
	}
	res = append(res, nstr)
	if col.Comment != "" {
		res = append(res, fmt.Sprintf("comment('%s')", col.Comment))
	}
	var tags []string
	tags = append(tags, "json:\""+col.Name+"\"")
	if !isCreated && !isUpdated && !isDeleted {
		tags = append(tags, "form:\""+col.Name+"\"")
	}
	if len(res) > 0 {
		tags = append(tags, "xorm:\""+strings.Join(res, " ")+"\"")
	}
	if len(tags) > 0 {
		return "`" + strings.Join(tags, " ") + "`"
	} else {
		return ""
	}
}

func unTitle(src string) string {
	if src == "" {
		return ""
	}

	if len(src) == 1 {
		return strings.ToLower(string(src[0]))
	} else {
		return strings.ToLower(string(src[0])) + src[1:]
	}
}

// gt evaluates the comparison a > b.
func gt(arg1, arg2 interface{}) (bool, error) {
	// > is the inverse of <=.
	lessOrEqual, err := le(arg1, arg2)
	if err != nil {
		return false, err
	}
	return !lessOrEqual, nil
}
func le(arg1, arg2 interface{}) (bool, error) {
	// <= is < or ==.
	lessThan, err := lt(arg1, arg2)
	if lessThan || err != nil {
		return lessThan, err
	}
	return eq(arg1, arg2)
}

func basicKind(v reflect.Value) (kind, error) {
	switch v.Kind() {
	case reflect.Bool:
		return boolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind, nil
	case reflect.Float32, reflect.Float64:
		return floatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return complexKind, nil
	case reflect.String:
		return stringKind, nil
	}
	return invalidKind, errBadComparisonType
}
func eq(arg1 interface{}, arg2 ...interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := basicKind(v1)
	if err != nil {
		return false, err
	}
	if len(arg2) == 0 {
		return false, errNoComparison
	}
	for _, arg := range arg2 {
		v2 := reflect.ValueOf(arg)
		k2, err := basicKind(v2)
		if err != nil {
			return false, err
		}
		if k1 != k2 {
			return false, errBadComparison
		}
		truth := false
		switch k1 {
		case boolKind:
			truth = v1.Bool() == v2.Bool()
		case complexKind:
			truth = v1.Complex() == v2.Complex()
		case floatKind:
			truth = v1.Float() == v2.Float()
		case intKind:
			truth = v1.Int() == v2.Int()
		case stringKind:
			truth = v1.String() == v2.String()
		case uintKind:
			truth = v1.Uint() == v2.Uint()
		default:
			panic("invalid kind")
		}
		if truth {
			return true, nil
		}
	}
	return false, nil
}
func lt(arg1, arg2 interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := basicKind(v1)
	if err != nil {
		return false, err
	}
	v2 := reflect.ValueOf(arg2)
	k2, err := basicKind(v2)
	if err != nil {
		return false, err
	}
	if k1 != k2 {
		return false, errBadComparison
	}
	truth := false
	switch k1 {
	case boolKind, complexKind:
		return false, errBadComparisonType
	case floatKind:
		truth = v1.Float() < v2.Float()
	case intKind:
		truth = v1.Int() < v2.Int()
	case stringKind:
		truth = v1.String() < v2.String()
	case uintKind:
		truth = v1.Uint() < v2.Uint()
	default:
		panic("invalid kind")
	}
	return truth, nil
}

var (
	errBadComparisonType = errors.New("invalid type for comparison")
	errBadComparison     = errors.New("incompatible types for comparison")
	errNoComparison      = errors.New("missing argument for comparison")
)

type kind int

const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	integerKind
	stringKind
	uintKind
)

func genGoImports(tables []*core.Table) map[string]string {
	imports := make(map[string]string)

	for _, table := range tables {
		for _, col := range table.Columns() {
			if typeString(col) == "time.Time" {
				imports["time"] = "time"
			}
		}
	}
	return imports
}
