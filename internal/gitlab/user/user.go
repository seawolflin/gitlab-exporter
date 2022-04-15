package user

import (
	"github.com/patrickmn/go-cache"
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"github.com/xanzy/go-gitlab"
	"log"
	"time"
)

type usersCache struct {
	cache []*gitlab.User
}

const CacheKey = "users"

func init() {
	initializer.Registry(func() {
		context.GetInstance().OnCacheEvicted(CacheKey, listUserFromGitlab)

		_, err := context.GetInstance().Cron().AddFunc("@every 24h", func() {
			listUserFromGitlab(context.GetInstance().Cache())
		})
		if err != nil {
			log.Fatalf("Add User Cron err, err: %s", err.Error())
		}
	})
}

func ListAll() []*gitlab.User {
	c := context.GetInstance().Cache()
	if users, found := c.Get(CacheKey); found {
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

	for _, user := range users {
		u := models.User{}
		u.AddOrUpdate(user)
	}

	c.Set(CacheKey, usersCache{users}, 24*time.Hour)
}
