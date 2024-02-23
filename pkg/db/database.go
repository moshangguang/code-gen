package db

import (
	"code-gen/utils/desUtils"
	"code-gen/utils/fileUtils"
	"code-gen/utils/logger"
	"code-gen/utils/md5"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"sync"
)

var database = &Database{
	RWMutex: new(sync.RWMutex),
	data:    map[string]Value{},
}

func init() {
}
func GetDatabase() *Database {
	return database
}

type Database struct {
	*sync.RWMutex
	data map[string]Value
}
type Value struct {
	dataBytes []byte
	exists    bool
}

func (table *Database) GetFilePath(key string) string {
	return filepath.Join(fileUtils.GetConfigDir(), md5.GetMd5(key))
}

func (table *Database) Remove(key string) {
	table.Lock()
	defer table.Unlock()
	table.data[key] = Value{
		dataBytes: nil,
		exists:    false,
	}
	filePath := table.GetFilePath(key)
	fileUtils.Remove(filePath)
}
func (table *Database) SaveAll() {
	table.Lock()
	defer table.Unlock()
	for key := range table.data {
		table.saveByKey(key, table.data[key].dataBytes)
	}
}
func (table *Database) saveByKey(key string, dataBytes []byte) {
	filePath := table.GetFilePath(key)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Logger.Error("保存文件时打开文件失败",
			zap.String("key", key),
			zap.Error(errors.WithStack(err)))
		return

	}
	defer fileUtils.Close(file)

	fileBytes, err := desUtils.DefaultECBEncrypt(dataBytes)
	if err != nil {
		logger.Logger.Fatal("保存文件时生成保护数据失败",
			zap.String("key", key),
			zap.Error(errors.WithMessage(err, "保存文件时写入文件失败")))
		return
	}
	if _, err = file.Write(fileBytes); err != nil {
		logger.Logger.Fatal("保存文件时写入文件失败",
			zap.String("key", key),
			zap.Error(errors.WithMessage(err, "保存文件时写入文件失败")))
	}
}

func (table *Database) Save(key string, v interface{}) {
	if v == nil {
		return
	}
	dataBytes, err := json.Marshal(v)
	if err != nil {
		logger.Logger.Fatal("数据序列化文件失败",
			zap.String("key", key),
			zap.Error(errors.WithMessage(err, "数据序列化文件失败")))
		return
	}

	table.Lock()
	defer table.Unlock()
	table.saveByKey(key, dataBytes)

	table.data[key] = Value{
		dataBytes: dataBytes,
		exists:    true,
	}
}
func (table *Database) getOnly(key string) ([]byte, bool) {
	val, ok := table.data[key]
	if ok && val.exists {
		result := make([]byte, len(val.dataBytes))
		copy(result, val.dataBytes)
		return result, true
	}
	return nil, false
}

func (table *Database) RGet(key string) ([]byte, bool) {
	table.RLock()
	defer table.RUnlock()
	return table.getOnly(key)
}

func (table *Database) Get(key string) ([]byte, bool) {
	dataBytes, ok := table.RGet(key)
	if ok {
		return dataBytes, true
	}
	table.Lock()
	defer table.Unlock()
	dataBytes, ok = table.getOnly(key)
	if ok {
		return dataBytes, true
	}
	filePath := table.GetFilePath(key)

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		table.data[key] = Value{
			dataBytes: nil,
			exists:    false,
		}
		return nil, false
	}
	dataBytes, err = os.ReadFile(filePath)
	if err != nil {
		logger.Logger.Error("读取配置文件失败", zap.String("key", key), zap.Error(errors.WithStack(err)))
		table.data[key] = Value{
			dataBytes: nil,
			exists:    false,
		}
		return nil, false
	}
	dataBytes, err = desUtils.DefaultECBDecrypt(dataBytes)
	if err != nil {
		table.data[key] = Value{
			dataBytes: nil,
			exists:    false,
		}
		return nil, false
	}
	table.data[key] = Value{
		dataBytes: dataBytes,
		exists:    true,
	}
	return table.getOnly(key)
}

func (table *Database) GetOrCreateString(key string, defaultString string) (str string) {
	dataBytes, ok := table.Get(key)
	if ok {
		if err := json.Unmarshal(dataBytes, &str); err == nil {
			return
		}
	}
	table.Save(key, defaultString)
	return defaultString
}

func (table *Database) GetBool(key string, defaultValue ...bool) bool {
	d := false
	if len(defaultValue) != 0 {
		d = defaultValue[0]
	}

	dataBytes, ok := table.Get(key)
	if !ok {
		return d
	}
	val := false
	err := json.Unmarshal(dataBytes, &val)
	if err != nil {
		logger.Logger.Error("获取bool键值对失败", zap.String("key", key), zap.Error(errors.WithStack(err)))
		return d
	}
	return val
}

func (table *Database) GetStringArrays(key string) []string {
	result := make([]string, 0)
	dataBytes, ok := table.Get(key)
	if !ok {
		return result
	}
	if err := json.Unmarshal(dataBytes, &result); err != nil {
		logger.Logger.Error("获取string arrays键值对失败", zap.String("key", key), zap.Error(errors.WithStack(err)))
	}
	return result
}

func (table *Database) GetString(key string, defaultValue ...string) string {
	d := ""
	if len(defaultValue) != 0 {
		d = defaultValue[0]
	}
	dataBytes, ok := table.Get(key)
	if !ok {
		return d
	}
	val := ""
	err := json.Unmarshal(dataBytes, &val)
	if err != nil {
		logger.Logger.Error("获取string键值对失败", zap.String("key", key), zap.Error(errors.WithStack(err)))
		return d
	}
	return val
}

func (table *Database) Unmarshal(key string, val interface{}) {
	dataBytes, ok := table.Get(key)
	if !ok {
		return
	}
	if err := json.Unmarshal(dataBytes, val); err != nil {
		logger.Logger.Error("反序列化配置失败", zap.String("key", key), zap.Error(errors.WithStack(err)))
	}
}

func (table *Database) UnmarshalWithOk(key string, val interface{}) bool {
	dataBytes, ok := table.Get(key)
	if !ok {
		return false
	}
	if err := json.Unmarshal(dataBytes, val); err != nil {
		logger.Logger.Error("反序列化配置失败", zap.String("key", key), zap.Error(errors.WithStack(err)))
		return false
	}
	return true
}
