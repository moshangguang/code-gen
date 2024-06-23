package files

import (
	"code-gen/utils/runtime"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	runtime.PanicError(err)
	return ioutil.ReadAll(f)
}

func OpenFile(filePth string) (file *os.File, err error) {
	file, err = os.Open(filePth)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(filePth)
	}
	return
}
func WriteTempFileContent(fileName string, content []byte) {
	filePath := fmt.Sprintf("%s%s%s", os.TempDir(), string(os.PathSeparator), fileName)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	runtime.PanicError(err)
	defer file.Close()
	n, _ := file.Seek(0, os.SEEK_END)
	_, err = file.WriteAt(content, n)
	runtime.PanicError(err)
}
func GetTempFileContent(fileName string) []byte {
	filePath := fmt.Sprintf("%s%s%s", os.TempDir(), string(os.PathSeparator), fileName)
	_, err := os.Stat(filePath)
	var file *os.File
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(filePath)
	} else {
		file, err = os.Open(filePath)
	}
	runtime.PanicError(err)
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	runtime.PanicError(err)
	return body
}
func Unmarshal(tmpFile string, obj interface{}) (bool, error) {
	content := GetTempFileContent(tmpFile)
	if len(content) == 0 {
		return false, nil
	}
	err := json.Unmarshal(content, obj)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func Marshal(tmpFile string, obj interface{}) error {
	if bytes, err := json.Marshal(obj); err != nil {
		return err
	} else {
		WriteTempFileContent(tmpFile, bytes)
	}
	return nil
}
