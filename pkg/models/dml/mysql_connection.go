package dml

import (
	"code-gen/pkg/models/ddl"
)

type MySQLConnectionModel struct {
}

func (m MySQLConnectionModel) GetAll() (ddl.MySQLConnectionSlice, error) {
	result := make(ddl.MySQLConnectionSlice, 0)
	err := engine.Find(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (m MySQLConnectionModel) GetByName(name string) (info ddl.MySQLConnection, exists bool, err error) {
	exists, err = engine.Where("name = ?", name).Get(&info)
	return
}
func (m MySQLConnectionModel) Insert(conn *ddl.MySQLConnection) error {
	_, err := engine.Insert(conn)
	return err
}
func (m MySQLConnectionModel) DeleteByName(name string) error {
	_, err := engine.Unscoped().Where("name = ?", name).Delete(ddl.MySQLConnection{})
	return err
}
func (m MySQLConnectionModel) Update(conn *ddl.MySQLConnection) error {
	_, err := engine.Table(ddl.MySQLConnection{}).Where("id = ?", conn.Id).Update(conn)
	return err
}
