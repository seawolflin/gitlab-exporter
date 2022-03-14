package gitlab

import (
	"github.com/patrickmn/go-cache"
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

type usersCache struct {
	cache []*gitlab.User
}

func init() {
	initializer.Registry(func() {
		context.GetInstance().OnCacheEvicted("users", listUserFromGitlab)
	})
}

func ListAll() []*gitlab.User {
	c := context.GetInstance().Cache()
	if users, found := c.Get("users"); found {
		log.Println("read from cache.")
		return users.(usersCache).cache
	}

	log.Println("Trigger list users from gitlab.")
	go listUserFromGitlab(c)

	return nil
}

func listUserFromGitlab(c *cache.Cache) {
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

	c.Set("users", usersCache{users}, 12*time.Second)
}
