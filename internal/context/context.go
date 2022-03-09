package context

import (
	"flag"
	"net/url"
	"sync"
)

type context struct {
	gitlabUrl    string
	privateToken string
	url          *url.URL
}

const (
	API_PREFIX = "/api/v4"
)

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
	u := url.URL{}
	url, err := u.Parse(ctx.gitlabUrl)
	if err != nil {
		panic("无效的Gitlab地址")
	}
	ctx.url = url
}

func (ctx *context) Url(path string) {

}
