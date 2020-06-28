package overseers

import (
	"github.com/herb-go/herb/cache"
	overseer "github.com/herb-go/providers/herb/overseers/cacheproxyoverseer"
	worker "github.com/herb-go/worker"
)

//CacheProxyWorker empty cache worker.
var CacheProxyWorker = cache.NewProxy(nil)

//CacheProxyOverseer cache overseer
var CacheProxyOverseer = worker.NewOrverseer("cacheproxy", &CacheProxyWorker)

func init() {
	CacheProxyOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(CacheProxyOverseer)
	})
	worker.Appoint(CacheProxyOverseer)
}
