package models

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"github.com/seawolflin/gitlab-exporter/internal/db"
)

func init() {
	// 必须在db初始化之后，再执行models的初始化
	initializer.RegistryByOrder(func() {
		Migrate()
	}, initializer.DEFUALT_ORDER+1)
}

func Migrate() {
	db.Migrate(&User{}, &Project{}, &Commit{})
}
