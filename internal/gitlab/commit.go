package gitlab

import (
	"github.com/seawolflin/gitlab-exporter/internal/common/utils"
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"github.com/xanzy/go-gitlab"
	"log"
)

func listCommitFromGitlab() {
	log.Println("listCommitFromGitlab")

	commitService := context.GetInstance().GitlabClient().Commits

	projectIds := models.Project{}.QueryAllProjectId()

	yesterdayStart := utils.TodayStart().AddDate(0, 0, -1)
	yesterdayEnd := utils.TodayEnd().AddDate(0, 0, -1)
	withStats := true

	for _, projectId := range projectIds {
		page := 1
		for {
			commits, _, err := commitService.ListCommits(projectId, &gitlab.ListCommitsOptions{
				ListOptions: gitlab.ListOptions{
					Page:    page,
					PerPage: 100,
				},
				Since:     &yesterdayStart,
				Until:     &yesterdayEnd,
				WithStats: &withStats,
			})
			if err != nil {
				log.Printf("获取项目%d的提交列表失败, %v", projectId, err)
				break
			}
			page++
			if len(commits) == 0 {
				break
			}

			for _, commit := range commits {
				models.Commit{}.AddOrUpdate(commit, &yesterdayStart)
			}
		}
	}
}
