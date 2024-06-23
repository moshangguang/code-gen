package types

type GolangSetting struct {
	PackageName string                `json:"package_name"`
	FileName    string                `json:"file_name"`
	ConnName    string                `json:"conn_name"`
	OutputPath  string                `json:"output_path"`
	ConnConfig  map[string]ConnConfig `json:"conn_config"`
}
type ConnConfig struct {
	Database string `json:"database"`
	Table    string `json:"table"`
}
