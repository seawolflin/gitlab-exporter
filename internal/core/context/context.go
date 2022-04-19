package context

import (
	"flag"
	"github.com/robfig/cron/v3"
	"github.com/xanzy/go-gitlab"
	"net/url"
	"sync"
)

type context struct {
	gitlabUrl    string
	privateToken string
	gitlabClient *gitlab.Client
	c            *cron.Cron
}

func (ctx *context) Cron() *cron.Cron {
	return ctx.c
}

func (ctx *context) GitlabClient() *gitlab.Client {
	return ctx.gitlabClient
}

var ins *context
var once = sync.Once{}

func GetInstance() *context {
	once.Do(func() {
		ins = &context{}
	})
	return ins
}

func init() {
	flag.StringVar(&GetInstance().gitlabUrl, "url", "", "Gitlab Url")
	flag.StringVar(&GetInstance().privateToken, "token", "", "Gitlab Private Token")
}

func (ctx *context) Parse() {
	flag.Parse()
	ctx.check()

	client, err := gitlab.NewClient(ins.privateToken, gitlab.WithBaseURL(ins.gitlabUrl))
	if err != nil {
		panic(err.Error())
	}
	ctx.gitlabClient = client

	ctx.c = cron.New()
	ctx.c.Start()
}

func (ctx *context) check() {
	if len(ctx.gitlabUrl) <= 0 {
		panic("url不能为空")
	}
	_, err := url.Parse(ctx.gitlabUrl)
	if err != nil {
		panic("无效的Gitlab url")
	}
	if len(ctx.privateToken) <= 0 {
		panic("token不能为空")
	}
}
