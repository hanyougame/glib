package etcdc

import (
	"github.com/jinzhu/copier"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"strings"
)

type Etcd[T any] struct {
	configurator configurator.Configurator[T]
}

// NewEtcd 实例化etcd
func NewEtcd[T any](c Config) *Etcd[T] {
	var cc subscriber.EtcdConf
	_ = copier.Copy(&cc, &c)
	cc.Hosts = strings.Split(c.Host, ",")

	return &Etcd[T]{
		configurator: configurator.MustNewConfigCenter[T](configurator.Config{
			Type: "json",
		}, subscriber.MustNewEtcdSubscriber(cc)),
	}
}

func (ctr *Etcd[T]) GetConfig() (T, error) {
	return ctr.configurator.GetConfig()
}

func (ctr *Etcd[T]) Listener(listener func(ec *Etcd[T])) {
	listener(ctr)
	ctr.configurator.AddListener(func() {
		listener(ctr)
	})
}
