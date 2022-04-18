package user

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"github.com/xanzy/go-gitlab"
	"log"
)

func init() {
	initializer.Registry(func() {
		_, err := context.GetInstance().Cron().AddFunc("@every 24h", func() {
			listUserFromGitlab()
		})
		if err != nil {
			log.Fatalf("Add User Cron err, err: %s", err.Error())
		}
		go listUserFromGitlab()
	})
}

func ListAll() []*models.User {
	users := models.User{}.QueryAll()

	if len(users) > 0 {
		log.Println("read from local db.")
		return users
	}

	return nil
}

func listUserFromGitlab() {
	log.Println("listUserFromGitlab")

	userService := context.GetInstance().GitlabClient().Users

	var users []*gitlab.User
	page := 1
	for {
		us, _, err := userService.ListUsers(&gitlab.ListUsersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			log.Fatalf("获取用户列表失败, %v", err)
		}
		page++
		if len(us) == 0 {
			break
		}
		users = append(users, us...)
	}

	for _, user := range users {
		models.User{}.AddOrUpdate(user)
	}
}
