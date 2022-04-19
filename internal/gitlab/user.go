package gitlab

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"github.com/xanzy/go-gitlab"
	"log"
)

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
