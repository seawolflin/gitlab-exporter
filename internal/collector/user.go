package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)
import gitlabUser "github.com/seawolflin/gitlab-exporter/internal/gitlab"

type userCollector struct {
	userStateMetric *prometheus.Desc
}

func init() {
	u := newUserCollector()
	prometheus.MustRegister(u)
}

func newUserCollector() *userCollector {
	return &userCollector{
		userStateMetric: prometheus.NewDesc("gitlab_user_state",
			"", []string{"name", "userName", "createAt", "isAdmin", "lastActivityOn"}, nil),
	}
}

func (collector *userCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.userStateMetric
}

func (collector *userCollector) Collect(ch chan<- prometheus.Metric) {
	users := gitlabUser.ListAll()

	for _, user := range users {
		ch <- prometheus.MustNewConstMetric(collector.userStateMetric, prometheus.GaugeValue,
			convertStateToValue(user.State),
			user.Name, user.Username, user.CreatedAt.String(), strconv.FormatBool(user.IsAdmin),
			func() string {
				if user.LastActivityOn != nil {
					return user.LastActivityOn.String()
				}
				return ""
			}())
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
