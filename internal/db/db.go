package db

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	initializer.Registry(func() {
		Open()
	})
}

func Open() {
	var err error
	DB, err = gorm.Open(sqlite.Open("gitlab.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("连接数据库失败！ err: " + err.Error())
	}
}

func Migrate(database ...interface{}) {
	err := DB.AutoMigrate(database...)
	if err != nil {
		panic("数据库迁移失败! err: " + err.Error())
	}
}
