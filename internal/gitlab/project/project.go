package project

import (
	"github.com/patrickmn/go-cache"
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"github.com/xanzy/go-gitlab"
	"log"
	"time"
)

type projectsCache struct {
	cache []*gitlab.Project
}

const CacheKey = "projects"

func init() {
	initializer.Registry(func() {
		context.GetInstance().OnCacheEvicted(CacheKey, listProjectFromGitlab)
	})
}

func ListAll() []*gitlab.Project {
	c := context.GetInstance().Cache()
	if users, found := c.Get(CacheKey); found {
		log.Println("read from cache.")
		return users.(projectsCache).cache
	}

	log.Println("Trigger list users from gitlab.")
	go listProjectFromGitlab(c)

	return nil
}

func listProjectFromGitlab(c *cache.Cache) {
	log.Println("listProjectFromGitlab")

	projectService := context.GetInstance().GitlabClient().Projects

	var projects []*gitlab.Project
	page := 1
	statistics := true
	for {
		us, _, err := projectService.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
			Statistics: &statistics,
		})
		if err != nil {
			log.Fatalf("获取项目列表失败, %v", err)
		}
		page++
		if len(us) == 0 {
			break
		}
		projects = append(projects, us...)
	}

	c.Set(CacheKey, projectsCache{projects}, 30*time.Minute)
}
