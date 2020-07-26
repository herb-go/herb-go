package overseers

import (
	overseer "github.com/herb-go/herb-drivers/overseers/cacheproxyoverseer"
	"github.com/herb-go/herb/cache"
	worker "github.com/herb-go/worker"
)

//CacheProxyWorker empty cacheproxy worker.
var CacheProxyWorker = cache.NewProxy(nil)

//CacheProxyOverseer cacheproxy overseer
var CacheProxyOverseer = worker.NewOrverseer("cacheproxy", &CacheProxyWorker)

func init() {
	CacheProxyOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(CacheProxyOverseer)
	})
	worker.Appoint(CacheProxyOverseer)
}
