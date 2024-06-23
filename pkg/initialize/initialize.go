package initialize

import (
	"code-gen/pkg/models/dml"
	"code-gen/utils/fileUtils"
	"fmt"
	"os"
)

func Init() {
	_, err := InitHome()
	if err != nil {
		panic(err)
	}
	dml.Init()
}
func InitHome() (string, error) {
	dir := fileUtils.GetConfigDir()
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return "", err
			}
		}
		stat, err = os.Stat(dir)
	}
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("home is not dir")
	}
	return dir, nil
}
