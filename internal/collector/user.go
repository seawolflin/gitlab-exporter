package collector

import "github.com/prometheus/client_golang/prometheus"

type userCollector struct {
	userStateMetric          *prometheus.Desc
	userLastActivityOnMetric *prometheus.Desc
	userIsAdminMetric        *prometheus.Desc
}

func init() {
	u := newUserCollector()
	prometheus.MustRegister(u)
}

func newUserCollector() *userCollector {
	return &userCollector{
		userStateMetric: prometheus.NewDesc("gitlab_user_state",
			"", []string{"name", "userName", "createAt"}, nil),
		userLastActivityOnMetric: prometheus.NewDesc("gitlab_user_last_activity_on",
			"", []string{"name", "userName", "createAt"}, nil),
		userIsAdminMetric: prometheus.NewDesc("gitlab_user_is_admin",
			"", []string{"name", "userName", "createAt"}, nil),
	}
}

func (collector *userCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.userStateMetric
	ch <- collector.userLastActivityOnMetric
	ch <- collector.userIsAdminMetric
}

func (collector *userCollector) Collect(ch chan<- prometheus.Metric) {

	ch <- prometheus.MustNewConstMetric(collector.userStateMetric, prometheus.CounterValue, 1,
		"Alice", "alice", "2022-01-10")
	ch <- prometheus.MustNewConstMetric(collector.userLastActivityOnMetric, prometheus.CounterValue, 1,
		"Alice", "alice", "2022-01-10")
	ch <- prometheus.MustNewConstMetric(collector.userIsAdminMetric, prometheus.CounterValue, 1,
		"Alice", "alice", "2022-01-10")
}
