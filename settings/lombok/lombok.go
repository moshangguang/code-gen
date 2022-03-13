package lombok

import "code-gen/utils/files"

type Lombok struct {
	Getter             bool
	Setter             bool
	Slf4j              bool
	NoArgsConstructor  bool
	AllArgsConstructor bool
	Data               bool
	ToString           bool
	EqualsAndHashCode  bool
}

var lombok = &Lombok{
	Getter:             true,
	Setter:             true,
	Slf4j:              false,
	NoArgsConstructor:  false,
	AllArgsConstructor: false,
	Data:               false,
	ToString:           false,
	EqualsAndHashCode:  false,
}

func (lombok *Lombok) Save() {
	_ = files.Marshal(SettingFile, *lombok)
}

func GetLombok() *Lombok {
	return lombok
}

const SettingFile = "lombok_code_gen.ini"

func init() {
	_, _ = files.Unmarshal(SettingFile, lombok)
}
