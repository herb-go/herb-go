package overseers

import (
	"github.com/herb-go/herb/cache"
	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/worker/overseers/cacheoverseer"
)

//CacheWorker empty cache worker.
var CacheWorker = cache.New()

//CacheOverseer cache overseer
var CacheOverseer = worker.NewOrverseer("cache", CacheWorker)

func init() {
	CacheOverseer.WithInitFunc(func() error {
		return overseer.New().Apply(CacheOverseer)
	})
}
