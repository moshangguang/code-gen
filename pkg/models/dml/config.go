package dml

import (
	"code-gen/pkg/models/ddl"
)

type ConfigModel struct {
}

func (m ConfigModel) GetByName(name string) (config ddl.Config, exists bool, err error) {
	exists, err = engine.Where("name = ?", name).Get(&config)
	return
}

func (m ConfigModel) Save(name, value string) error {
	config, exists, err := m.GetByName(name)
	if err != nil {
		return err
	}
	if exists {
		config.Value = value
		if _, err = engine.Where("id = ?", config.Id).Update(config); err != nil {
			return err
		}
	} else {
		if _, err = engine.Insert(&ddl.Config{
			Name:  name,
			Value: value,
		}); err != nil {
			return err
		}
	}
	return nil
}
