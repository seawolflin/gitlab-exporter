package project

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"github.com/xanzy/go-gitlab"
	"log"
)

func init() {
	initializer.Registry(func() {
		_, err := context.GetInstance().Cron().AddFunc("@every 24m", func() {
			listProjectFromGitlab()
		})
		if err != nil {
			log.Fatalf("Add User Cron err, err: %s", err.Error())
		}

		go listProjectFromGitlab()
	})
}

func ListAll() []*models.Project {
	projects := models.Project{}.QueryAll()

	if len(projects) > 0 {
		log.Println("read from local db.")
		return projects
	}

	return nil
}

func listProjectFromGitlab() {
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

		for _, project := range projects {
			models.Project{}.AddOrUpdate(project)
		}
	}
}
