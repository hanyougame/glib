package redisx

import (
	"context"
	"fmt"
	"github.com/hanyougame/glib/stores/redisx/config"
	"testing"
)

func TestRedis(t *testing.T) {
	Must(config.Config{
		Addrs:    []string{"k8s-k1-k1redisr-f4e8a8a11b-8ba79625458908d4.elb.ap-east-1.amazonaws.com:6379"},
		Debug:    true,
		Trace:    true,
		Password: "ObGWx6761jcUV5um5wEeFYvv",
	})
	r, err := Engine.Get(context.Background(), "test").Result()
	fmt.Println(err)
	fmt.Println(r)
}
