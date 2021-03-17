package overseers

import (
	"github.com/herb-go/deprecated/cache"
	overseer "github.com/herb-go/deprecated/herb-drivers/overseers/cacheproxyoverseer"
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
