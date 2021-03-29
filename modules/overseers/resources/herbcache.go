package overseers

import (
	"github.com/herb-go/datamodules/herbcache"
	overseer "github.com/herb-go/herb-drivers/overseers/herbcacheoverseer"
	worker "github.com/herb-go/worker"
)

//NCacheWorker empty cache worker.
var HerbCacheWorker = herbcache.New()

//NCacheOverseer ncache overseer
var HerbCacheOverseer = worker.NewOrverseer("cache", &HerbCacheWorker)

func init() {
	HerbCacheOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(HerbCacheOverseer)
	})
	worker.Appoint(HerbCacheOverseer)
}
