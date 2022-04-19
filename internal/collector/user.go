package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/seawolflin/gitlab-exporter/internal/common/utils"
	"github.com/seawolflin/gitlab-exporter/internal/models"
	"strconv"
)

type userCollector struct {
	userStateMetric          *prometheus.Desc
	userCreateAtMetric       *prometheus.Desc
	userIsAdminMetric        *prometheus.Desc
	userLastActivityOnMetric *prometheus.Desc
}

func init() {
	u := newUserCollector()
	prometheus.MustRegister(u)
}

func newUserCollector() *userCollector {
	return &userCollector{
		userStateMetric: prometheus.NewDesc("gitlab_user_state",
			"用户状态，1 - active, 2 - blocked, -1 - 未知", []string{"id", "name", "userName", "email"}, nil),
		userCreateAtMetric: prometheus.NewDesc("gitlab_user_create_at",
			"用户创建时间，单位：s", []string{"id", "name", "userName", "email"}, nil),
		userIsAdminMetric: prometheus.NewDesc("gitlab_user_is_admin",
			"用户是否是管理员", []string{"id", "name", "userName", "email"}, nil),
		userLastActivityOnMetric: prometheus.NewDesc("gitlab_user_last_activity_on",
			"用户活跃时间，单位：s", []string{"id", "name", "userName", "email"}, nil),
	}
}

func (collector *userCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.userStateMetric
	ch <- collector.userCreateAtMetric
	ch <- collector.userIsAdminMetric
	ch <- collector.userLastActivityOnMetric
}

func (collector *userCollector) Collect(ch chan<- prometheus.Metric) {
	users := models.User{}.QueryAll()

	for _, user := range users {
		ch <- prometheus.MustNewConstMetric(collector.userStateMetric, prometheus.GaugeValue,
			convertStateToValue(user.State),
			strconv.Itoa(user.GitlabId), user.Name, user.Username, user.Email)
		ch <- prometheus.MustNewConstMetric(collector.userCreateAtMetric, prometheus.GaugeValue,
			float64(user.GitlabCreatedAt.Unix()),
			strconv.Itoa(user.GitlabId), user.Name, user.Username, user.Email)
		ch <- prometheus.MustNewConstMetric(collector.userIsAdminMetric, prometheus.GaugeValue,
			utils.ConvertBoolToValue(user.IsAdmin),
			strconv.Itoa(user.GitlabId), user.Name, user.Username, user.Email)

		if user.LastActivityOn != nil {
			lastActivityOn := *user.LastActivityOn
			ch <- prometheus.MustNewConstMetric(collector.userLastActivityOnMetric, prometheus.GaugeValue,
				float64(lastActivityOn.Unix()),
				strconv.Itoa(user.GitlabId), user.Name, user.Username, user.Email)
		}
	}
}

func convertStateToValue(state string) float64 {
	switch state {
	case "active":
		return 1
	case "blocked":
		return 2
	default:
		return -1
	}
}
