package overseers

import (
	"github.com/herb-go/herb/cache"
	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/worker/overseers/cacheoverseer"
)

//CacheableWorker empty cache worker.
var CacheableWorker cache.Cacheable

//CacheableOverseer cache overseer
var CacheableOverseer = worker.NewOrverseer("cacheable", CacheableWorker)

func init() {
	CacheableOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().Apply(CacheableOverseer)
	})
}
