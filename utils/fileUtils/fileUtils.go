package fileUtils

import (
	"code-gen/utils/exceptUtils"
	"code-gen/utils/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"
)

func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return true
}
func IsNotExist(filePath string) bool {
	return !IsExist(filePath)
}

func IsDir(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func CreateAndWrite(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer Close(file)
	_, err = file.Write(data)
	return err
}

func Close(file *os.File) {
	if file == nil {
		return
	}
	err := file.Close()
	if err != nil {
		logger.Logger.Error("关闭文件失败",
			zap.String("fileName", file.Name()),
			zap.Error(errors.WithMessage(err, "关闭文件失败")))
	}
}
func GetRootDir() string {
	current, err := user.Current()
	if err != nil {
		panic(err)
	}
	return current.HomeDir
}
func GetConfigDir() string {
	return filepath.Join(GetRootDir(), "code-gen")
}
func Remove(filePath string) {
	if IsNotExist(filePath) {
		return
	}
	exceptUtils.CatchErrorWithMessage("删除文件出错", func() error {
		return os.Remove(filePath)
	})
}
func RemoveAll(filePath string) {
	if IsNotExist(filePath) {
		return
	}
	exceptUtils.CatchErrorWithMessage("删除所有文件出错", func() error {
		return os.RemoveAll(filePath)
	})
}
func GetTempDir() string {
	return os.TempDir()
}
