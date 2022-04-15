package db

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"github.com/seawolflin/gitlab-exporter/internal/models"
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

	migrate()
}

func migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		panic("数据库迁移失败! err: " + err.Error())
	}

	//pkg, err := importer.Default().Import("github.com/seawolflin/gitlab-exporter/internal/models")
	//if err != nil {
	//	panic("数据库迁移失败")
	//}
	//
	//for _, declName := range pkg.Scope().Names() {
	//	fmt.Println(declName)
	//}
}
