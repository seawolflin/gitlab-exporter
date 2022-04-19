package models

import (
	"github.com/seawolflin/gitlab-exporter/internal/common/utils"
	"github.com/seawolflin/gitlab-exporter/internal/db"
	"github.com/xanzy/go-gitlab"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	gorm.Model
	GitlabId                       int
	Description                    string
	DefaultBranch                  string
	Public                         bool
	Visibility                     string
	SSHURLToRepo                   string
	HTTPURLToRepo                  string
	WebURL                         string
	ReadmeURL                      string
	TagList                        string
	Topics                         string
	OwnerId                        int
	Name                           string
	NameWithNamespace              string
	Path                           string
	PathWithNamespace              string
	IssuesEnabled                  bool
	OpenIssuesCount                int
	MergeRequestsEnabled           bool
	ApprovalsBeforeMerge           int
	JobsEnabled                    bool
	WikiEnabled                    bool
	SnippetsEnabled                bool
	ResolveOutdatedDiffDiscussions bool
	ContainerExpirationPolicy      ContainerExpirationPolicy `gorm:"embedded;embeddedPrefix:container_"`
	ContainerRegistryEnabled       bool
	ContainerRegistryAccessLevel   string
	GitlabCreatedAt                *time.Time
	LastActivityAt                 *time.Time
	CreatorID                      int
	Namespace                      ProjectNamespace `gorm:"embedded;embeddedPrefix:namespace_"`
	ImportStatus                   string
	ImportError                    string
	//Permissions                               *Permissions
	//MarkedForDeletionAt                       *ISOTime
	EmptyRepo                                 bool
	Archived                                  bool
	AvatarURL                                 string
	LicenseURL                                string
	License                                   ProjectLicense `gorm:"embedded;embeddedPrefix:project_license_"`
	SharedRunnersEnabled                      bool
	ForksCount                                int
	StarCount                                 int
	RunnersToken                              string
	PublicBuilds                              bool
	AllowMergeOnSkippedPipeline               bool
	OnlyAllowMergeIfPipelineSucceeds          bool
	OnlyAllowMergeIfAllDiscussionsAreResolved bool
	RemoveSourceBranchAfterMerge              bool
	PrintingMergeRequestLinkEnabled           bool
	LFSEnabled                                bool
	RepositoryStorage                         string
	RequestAccessEnabled                      bool
	MergeMethod                               string
	ForkedFromProject                         ForkParent `gorm:"embedded;embeddedPrefix:fork_parent_"`
	Mirror                                    bool
	MirrorUserID                              int
	MirrorTriggerBuilds                       bool
	OnlyMirrorProtectedBranches               bool
	MirrorOverwritesDivergedBranches          bool
	PackagesEnabled                           bool
	ServiceDeskEnabled                        bool
	ServiceDeskAddress                        string
	IssuesAccessLevel                         string
	RepositoryAccessLevel                     string
	MergeRequestsAccessLevel                  string
	ForkingAccessLevel                        string
	WikiAccessLevel                           string
	BuildsAccessLevel                         string
	SnippetsAccessLevel                       string
	PagesAccessLevel                          string
	OperationsAccessLevel                     string
	AnalyticsAccessLevel                      string
	AutocloseReferencedIssues                 bool
	SuggestionCommitMessage                   string
	AutoCancelPendingPipelines                string
	CIForwardDeploymentEnabled                bool
	SquashOption                              string
	//SharedWithGroups                          []struct {
	//	GroupID          int    `json:"group_id"`
	//	GroupName        string `json:"group_name"`
	//	GroupAccessLevel int    `json:"group_access_level"`
	//} `json:"shared_with_groups"`
	Statistics        ProjectStatistics `gorm:"embedded;embeddedPrefix:statistics_"`
	Links             Links             `gorm:"embedded;embeddedPrefix:links_"`
	CIConfigPath      string
	CIDefaultGitDepth int
	//CustomAttributes                         []*CustomAttribute `gorm:"embedded;embeddedPrefix:custom_attributes_"`
	ComplianceFrameworks                     string
	BuildCoverageRegex                       string
	BuildTimeout                             int
	IssuesTemplate                           string
	MergeRequestsTemplate                    string
	KeepLatestArtifact                       bool
	MergePipelinesEnabled                    bool
	MergeTrainsEnabled                       bool
	RestrictUserDefinedVariables             bool
	MergeCommitTemplate                      string
	SquashCommitTemplate                     string
	AutoDevopsDeployStrategy                 string
	AutoDevopsEnabled                        bool
	BuildGitStrategy                         string
	EmailsDisabled                           bool
	ExternalAuthorizationClassificationLabel string
	RequirementsAccessLevel                  int
	SecurityAndComplianceAccessLevel         int
}

type ContainerExpirationPolicy struct {
	Cadence         string
	KeepN           int
	OlderThan       string
	NameRegexDelete string
	NameRegexKeep   string
	Enabled         bool
	NextRunAt       *time.Time
}

type ProjectNamespace struct {
	ID        int
	Name      string
	Path      string
	Kind      string
	FullPath  string
	AvatarURL string
	WebURL    string
}

type ProjectLicense struct {
	Key       string
	Name      string
	Nickname  string
	HTMLURL   string
	SourceURL string
}

type ForkParent struct {
	HTTPURLToRepo     string
	ID                int
	Name              string
	NameWithNamespace string
	Path              string
	PathWithNamespace string
	WebURL            string
}

type StorageStatistics struct {
	StorageSize      int64
	RepositorySize   int64
	LfsObjectsSize   int64
	JobArtifactsSize int64
}

type ProjectStatistics struct {
	StorageStatistics
	CommitCount int
}

type Links struct {
	Self          string
	Issues        string
	MergeRequests string
	RepoBranches  string
	Labels        string
	Events        string
	Members       string
}

type CustomAttribute struct {
	Key   string
	Value string
}

func (p Project) AddOrUpdate(project *gitlab.Project) {
	db.DB.Where("gitlab_id = ?", project.ID).First(&p)

	utils.CopyStruct(project, &p)

	p.GitlabId = project.ID
	p.Visibility = string(project.Visibility)
	p.TagList = utils.ListToString(project.TagList)
	p.Topics = utils.ListToString(project.Topics)

	if project.ContainerExpirationPolicy != nil {
		utils.CopyStruct(project.ContainerExpirationPolicy, &p.ContainerExpirationPolicy)
	}
	p.ContainerRegistryAccessLevel = string(project.ContainerRegistryAccessLevel)
	if project.Namespace != nil {
		utils.CopyStruct(project.Namespace, &p.Namespace)
	}
	if project.License != nil {
		utils.CopyStruct(project.License, &p.License)
	}
	p.MergeMethod = string(project.MergeMethod)
	if project.ForkedFromProject != nil {
		utils.CopyStruct(project.ForkedFromProject, &p.ForkedFromProject)
	}
	p.IssuesAccessLevel = string(project.IssuesAccessLevel)
	p.RepositoryAccessLevel = string(project.RepositoryAccessLevel)
	p.MergeRequestsAccessLevel = string(project.MergeRequestsAccessLevel)
	p.ForkingAccessLevel = string(project.ForkingAccessLevel)
	p.WikiAccessLevel = string(project.WikiAccessLevel)
	p.BuildsAccessLevel = string(project.BuildsAccessLevel)
	p.SnippetsAccessLevel = string(project.SnippetsAccessLevel)
	p.PagesAccessLevel = string(project.PagesAccessLevel)
	p.OperationsAccessLevel = string(project.OperationsAccessLevel)
	p.AnalyticsAccessLevel = string(project.AnalyticsAccessLevel)
	p.SquashOption = string(project.SquashOption)
	if project.Statistics != nil {
		utils.CopyStruct(project.Statistics, &p.Statistics)
	}
	if project.Links != nil {
		utils.CopyStruct(project.Links, &p.Links)
	}

	p.ComplianceFrameworks = utils.ListToString(project.ComplianceFrameworks)

	db.DB.Save(&p)
}

func (p Project) QueryAll() []*Project {
	var projects []*Project
	db.DB.Find(&projects)

	return projects
}

func (p Project) QueryAllProjectId() []int {
	var results []int
	db.DB.Model(&Project{}).Select("gitlab_id").Find(&results)

	return results
}
