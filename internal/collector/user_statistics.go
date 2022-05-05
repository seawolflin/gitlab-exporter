package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/seawolflin/gitlab-exporter/internal/common/utils"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"strconv"
)

type userStatisticsCollector struct {
	userDailyCommitCntMetric        *prometheus.Desc
	userDailyAdditionLinesMetric    *prometheus.Desc
	userDailyDeletionLinesMetric    *prometheus.Desc
	userDailyTotalModifyLinesMetric *prometheus.Desc
}

func init() {
	u := newUserStatisticsCollector()
	prometheus.MustRegister(u)
}

func newUserStatisticsCollector() *userStatisticsCollector {
	return &userStatisticsCollector{
		userDailyCommitCntMetric: prometheus.NewDesc("gitlab_user_daily_commit_cnt",
			"用户每日提交数", []string{"id", "name", "userName", "email", "date"}, nil),
		userDailyAdditionLinesMetric: prometheus.NewDesc("gitlab_user_daily_addition_lines",
			"用户每日增加行数", []string{"id", "name", "userName", "email", "date"}, nil),
		userDailyDeletionLinesMetric: prometheus.NewDesc("gitlab_user_daily_deletion_lines",
			"用户每日删除行数", []string{"id", "name", "userName", "email", "date"}, nil),
		userDailyTotalModifyLinesMetric: prometheus.NewDesc("gitlab_user_daily_total_modify_lines",
			"用户每日总共修改行数", []string{"id", "name", "userName", "email", "date"}, nil),
	}
}

func (collector *userStatisticsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.userDailyCommitCntMetric
	ch <- collector.userDailyAdditionLinesMetric
	ch <- collector.userDailyDeletionLinesMetric
	ch <- collector.userDailyTotalModifyLinesMetric
}

func (collector *userStatisticsCollector) Collect(ch chan<- prometheus.Metric) {
	yesterdayStart := utils.TodayStart().AddDate(0, 0, -1)
	stats := models.Commit{}.QueryStats(&yesterdayStart)

	for email, commitStats := range stats {
		user := models.User{}.QueryByEmail(email)
		ch <- prometheus.MustNewConstMetric(collector.userDailyCommitCntMetric, prometheus.GaugeValue,
			float64(commitStats.CommitCnt),
			strconv.Itoa(user.GitlabId), commitStats.AuthorName, user.Username, email, commitStats.CommitDate)
		ch <- prometheus.MustNewConstMetric(collector.userDailyAdditionLinesMetric, prometheus.GaugeValue,
			float64(commitStats.Additions),
			strconv.Itoa(user.GitlabId), commitStats.AuthorName, user.Username, email, commitStats.CommitDate)
		ch <- prometheus.MustNewConstMetric(collector.userDailyDeletionLinesMetric, prometheus.GaugeValue,
			float64(commitStats.Deletions),
			strconv.Itoa(user.GitlabId), commitStats.AuthorName, user.Username, email, commitStats.CommitDate)
		ch <- prometheus.MustNewConstMetric(collector.userDailyTotalModifyLinesMetric, prometheus.GaugeValue,
			float64(commitStats.Total),
			strconv.Itoa(user.GitlabId), commitStats.AuthorName, user.Username, email, commitStats.CommitDate)
	}
}
