package ddl

import "time"

type Config struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
}

func (Config) TableName() string {
	return "config"
}
