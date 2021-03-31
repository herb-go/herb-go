package overseers

import (
	overseer "github.com/herb-go/datamodule-drivers/overseers/cachestoragebuilderoverseer"
	"github.com/herb-go/datamodules/herbcache"
	worker "github.com/herb-go/worker"
)

//HerbCacheStorageBuilderWorker empty cache worker.
var HerbCacheStorageBuilderWorker func(*herbcache.Storage, func(v interface{}) error) error

//HerbCacheStorageBuilderOverseer ncache overseer
var HerbCacheStorageBuilderOverseer = worker.NewOrverseer("cache", &HerbCacheStorageBuilderWorker)

func init() {
	HerbCacheStorageBuilderOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(HerbCacheStorageBuilderOverseer)
	})
	worker.Appoint(HerbCacheStorageBuilderOverseer)
}
