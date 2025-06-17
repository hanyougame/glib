package httpc

import (
	"context"
	"github.com/go-resty/resty/v2"
	"sync"
)

// 定义http 引擎
var engine *resty.Client
var once sync.Once

func Do(ctx context.Context, fs ...func(cli *resty.Client)) *resty.Request {
	once.Do(func() {
		engine = MustClient()
	})
	for _, f := range fs {
		f(engine)
	}
	return engine.R().SetContext(ctx)
}

// MustClient new http client
func MustClient() *resty.Client {
	return resty.New()
}
