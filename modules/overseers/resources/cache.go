package overseers

import (
	"github.com/herb-go/herb/cache"
	overseer "github.com/herb-go/herb-drivers/overseers/cacheoverseer"
	worker "github.com/herb-go/worker"
)

//CacheWorker empty cache worker.
var CacheWorker = cache.New()

//CacheOverseer cache overseer
var CacheOverseer = worker.NewOrverseer("cache", &CacheWorker)

func init() {
	CacheOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(CacheOverseer)
	})
	worker.Appoint(CacheOverseer)
}
