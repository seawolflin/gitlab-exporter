package models

import (
	"github.com/seawolflin/gitlab-exporter/internal/common/utils"
	"github.com/seawolflin/gitlab-exporter/internal/db"
	"github.com/xanzy/go-gitlab"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	GitlabId                       int `gorm:"uniqueIndex"`
	Username                       string
	Email                          string
	Name                           string
	State                          string
	WebURL                         string
	GitlabCreatedAt                *time.Time
	Bio                            string
	Location                       string
	PublicEmail                    string
	Skype                          string
	Linkedin                       string
	Twitter                        string
	WebsiteURL                     string
	Organization                   string
	JobTitle                       string
	ExternUID                      string
	Provider                       string
	ThemeID                        int
	LastActivityOn                 *time.Time
	ColorSchemeID                  int
	IsAdmin                        bool
	AvatarURL                      string
	CanCreateGroup                 bool
	CanCreateProject               bool
	ProjectsLimit                  int
	CurrentSignInAt                *time.Time
	CurrentSignInIP                string
	LastSignInAt                   *time.Time
	LastSignInIP                   string
	ConfirmedAt                    *time.Time
	TwoFactorEnabled               bool
	Note                           string
	External                       bool
	PrivateProfile                 bool
	SharedRunnersMinutesLimit      int
	ExtraSharedRunnersMinutesLimit int
	UsingLicenseSeat               bool
}

func (u User) AddOrUpdate(user *gitlab.User) {
	db.DB.Where("gitlab_id = ?", user.ID).First(&u)

	utils.CopyStruct(user, &u)

	u.GitlabId = user.ID
	u.GitlabCreatedAt = user.CreatedAt
	u.LastActivityOn = (*time.Time)(user.LastActivityOn)

	db.DB.Save(&u)
}

func (u User) Query(gitlabId int) *User {
	db.DB.Where("gitlab_id = ?", gitlabId).First(&u)

	return &u
}

func (u User) QueryAll() []*User {
	var users []*User
	db.DB.Find(&users)

	return users
}
