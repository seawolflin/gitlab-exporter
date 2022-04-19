package gitlab

import (
	"github.com/seawolflin/gitlab-exporter/internal/core/context"
	"github.com/seawolflin/gitlab-exporter/internal/core/initializer"
	"log"
)

func init() {
	initializer.Registry(func() {
		_, err := context.GetInstance().Cron().AddFunc("@every 24h", func() {
			go listUserFromGitlab()
			go listProjectFromGitlab()
			go listCommitFromGitlab()
		})
		if err != nil {
			log.Fatalf("Add Cron err, err: %s", err.Error())
		}
		go listUserFromGitlab()
		go listProjectFromGitlab()
		go listCommitFromGitlab()
	})
}
