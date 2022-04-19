package gitlab

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"github.com/xanzy/go-gitlab"
	"log"
)

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
