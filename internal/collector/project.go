package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/seawolflin/gitlab-exporter/internal/common/utils"
	"github.com/seawolflin/gitlab-exporter/internal/models"
)

type projectCollector struct {
	projectCommitCountMetric      *prometheus.Desc
	projectForksCountMetric       *prometheus.Desc
	projectStarCountMetric        *prometheus.Desc
	projectLastActivityAtMetric   *prometheus.Desc
	projectLfsEnableMetric        *prometheus.Desc
	projectStorageSizeMetric      *prometheus.Desc
	projectRepositorySizeMetric   *prometheus.Desc
	projectLfsObjectsSizeMetric   *prometheus.Desc
	projectJobArtifactsSizeMetric *prometheus.Desc
}

func init() {
	u := newProjectCollector()
	prometheus.MustRegister(u)
}

func newProjectCollector() *projectCollector {
	return &projectCollector{
		projectCommitCountMetric: prometheus.NewDesc("gitlab_project_commit_count",
			"commit数量", []string{"repo", "name"}, nil),
		projectForksCountMetric: prometheus.NewDesc("gitlab_project_forks_count",
			"forks数量", []string{"repo", "name"}, nil),
		projectStarCountMetric: prometheus.NewDesc("gitlab_project_star_count",
			"star数量", []string{"repo", "name"}, nil),
		projectLastActivityAtMetric: prometheus.NewDesc("gitlab_project_last_activity_at",
			"最近活跃时间", []string{"repo", "name"}, nil),
		projectLfsEnableMetric: prometheus.NewDesc("gitlab_project_lfs_enabled",
			"是否开启lfs", []string{"repo", "name"}, nil),
		projectStorageSizeMetric: prometheus.NewDesc("gitlab_project_storage_size",
			"存储大小", []string{"repo", "name"}, nil),
		projectRepositorySizeMetric: prometheus.NewDesc("gitlab_project_repository_size",
			"仓库大小", []string{"repo", "name"}, nil),
		projectJobArtifactsSizeMetric: prometheus.NewDesc("gitlab_project_job_artifacts_size",
			"Job产物大小", []string{"repo", "name"}, nil),
	}
}

func (collector *projectCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.projectCommitCountMetric
	ch <- collector.projectForksCountMetric
	ch <- collector.projectStarCountMetric
	ch <- collector.projectLastActivityAtMetric
	ch <- collector.projectLfsEnableMetric
	ch <- collector.projectStorageSizeMetric
	ch <- collector.projectRepositorySizeMetric
	ch <- collector.projectJobArtifactsSizeMetric
}

func (collector *projectCollector) Collect(ch chan<- prometheus.Metric) {
	projects := models.Project{}.QueryAll()

	for _, project := range projects {
		ch <- prometheus.MustNewConstMetric(collector.projectCommitCountMetric, prometheus.GaugeValue,
			float64(project.Statistics.CommitCount), project.PathWithNamespace, project.Name)
		ch <- prometheus.MustNewConstMetric(collector.projectForksCountMetric, prometheus.GaugeValue,
			float64(project.ForksCount), project.PathWithNamespace, project.Name)
		ch <- prometheus.MustNewConstMetric(collector.projectStarCountMetric, prometheus.GaugeValue,
			float64(project.StarCount), project.PathWithNamespace, project.Name)
		ch <- prometheus.MustNewConstMetric(collector.projectLastActivityAtMetric, prometheus.GaugeValue,
			float64(project.LastActivityAt.Unix()), project.PathWithNamespace, project.Name)
		ch <- prometheus.MustNewConstMetric(collector.projectLfsEnableMetric, prometheus.GaugeValue,
			utils.ConvertBoolToValue(project.LFSEnabled), project.PathWithNamespace, project.Name)
		ch <- prometheus.MustNewConstMetric(collector.projectStorageSizeMetric, prometheus.GaugeValue,
			float64(project.Statistics.StorageSize), project.PathWithNamespace, project.Name)
		ch <- prometheus.MustNewConstMetric(collector.projectRepositorySizeMetric, prometheus.GaugeValue,
			float64(project.Statistics.RepositorySize), project.PathWithNamespace, project.Name)
		ch <- prometheus.MustNewConstMetric(collector.projectJobArtifactsSizeMetric, prometheus.GaugeValue,
			float64(project.Statistics.JobArtifactsSize), project.PathWithNamespace, project.Name)
	}
}
