package overseers

import (
	"github.com/herb-go/herb/cache"
	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/worker/overseers/cacheproxyoverseer"
)

//CacheProxyWorker empty cache worker.
var CacheProxyWorker = cache.NewProxy(nil)

//CacheProxyOverseer cache overseer
var CacheProxyOverseer = worker.NewOrverseer("cacheproxy", CacheProxyWorker)

func init() {
	CacheProxyOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().Apply(CacheProxyOverseer)
	})
}
