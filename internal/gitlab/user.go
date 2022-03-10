package gitlab

import (
	"log"
	"time"

	"github.com/seawolflin/gitlab-exporter/internal/context"
	"github.com/xanzy/go-gitlab"
)

type usersCache struct {
	cache []*gitlab.User
}

func ListAll() []*gitlab.User {
	c := context.GetInstance().Cache()
	if users, found := c.Get("users"); found {
		log.Println("read from cache.")
		return users.(usersCache).cache
	}

	userService := context.GetInstance().GitlabClient().Users

	var users []*gitlab.User
	for {
		us, _, err := userService.ListUsers(&gitlab.ListUsersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 1,
			},
		})
		if err != nil {
			log.Fatalf("获取用户列表失败, %v", err)
		}
		if len(users) == 0 {
			break
		}
		users = append(users, us...)
	}

	c.Set("users", usersCache{users}, 12*time.Hour)

	return users
}
