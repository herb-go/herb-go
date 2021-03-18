package overseers

import (
	"github.com/herb-go/datamodules/ncache"
	overseer "github.com/herb-go/herb-drivers/overseers/ncacheoverseer"
	worker "github.com/herb-go/worker"
)

//NCacheWorker empty cache worker.
var NCacheWorker = ncache.New()

//NCacheOverseer ncache overseer
var NCacheOverseer = worker.NewOrverseer("cache", &NCacheWorker)

func init() {
	NCacheOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(NCacheOverseer)
	})
	worker.Appoint(NCacheOverseer)
}
