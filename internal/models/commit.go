package models

import (
	"github.com/seawolflin/gitlab-exporter/internal/common/utils"
	"github.com/seawolflin/gitlab-exporter/internal/db"
	"github.com/xanzy/go-gitlab"
	"gorm.io/gorm"
	"time"
)

type Commit struct {
	gorm.Model
	Date            *time.Time `gorm:"index"`
	GitlabId        string     `gorm:"uniqueIndex"`
	ShortID         string     `gorm:"index"`
	Title           string
	AuthorName      string
	AuthorEmail     string
	AuthoredDate    *time.Time
	CommitterName   string
	CommitterEmail  string
	CommittedDate   *time.Time
	GitlabCreatedAt *time.Time
	Message         string
	Stats           CommitStats `gorm:"embedded;embeddedPrefix:stats_"`
	Status          string
	ProjectID       int
	WebURL          string
}

type CommitStats struct {
	Additions int
	Deletions int
	Total     int
}

func (c Commit) AddOrUpdate(commit *gitlab.Commit, date *time.Time) {
	db.DB.Where("gitlab_id = ?", commit.ID).First(&c)

	utils.CopyStruct(commit, &c)

	c.GitlabId = commit.ID
	c.GitlabCreatedAt = commit.CreatedAt
	if commit.Status != nil {
		c.Status = string(*commit.Status)
	}
	c.Date = date
	utils.CopyStruct(commit.Stats, &c.Stats)

	db.DB.Save(&c)
}

type UserCommitStats struct {
	CommitStats
	CommitCnt  int
	CommitDate string
}

func (c Commit) QueryStats(date *time.Time) map[string]UserCommitStats {
	var results []struct {
		AuthorEmail string
		UserCommitStats
	}

	db.DB.
		Model(&Commit{}).
		Select("author_email, sum(stats_additions) as additions,"+
			" sum(stats_deletions) as deletions, sum(stats_total) as total, count(1) as commit_cnt").
		Group("author_email").Where("date = ?", date).
		Find(&results)

	var stats = make(map[string]UserCommitStats)
	for _, result := range results {
		stats[result.AuthorEmail] = UserCommitStats{
			CommitCnt:  result.CommitCnt,
			CommitDate: date.Format("2006-01-02"),
			CommitStats: CommitStats{
				Additions: result.Additions,
				Deletions: result.Deletions,
				Total:     result.Total,
			},
		}
	}

	return stats
}
