package overseers

import (
	overseer "github.com/herb-go/datamodule-drivers/overseers/herbcacheoverseer"
	"github.com/herb-go/datamodules/herbcache"
	worker "github.com/herb-go/worker"
)

//HerbCacheWorker empty cache worker.
var HerbCacheWorker = herbcache.New()

//HerbCacheOverseer herbcache overseer
var HerbCacheOverseer = worker.NewOrverseer("cache", &HerbCacheWorker)

func init() {
	HerbCacheOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(HerbCacheOverseer)
	})
	worker.Appoint(HerbCacheOverseer)
}
