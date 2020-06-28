package overseers

import (
	"github.com/herb-go/herb/middleware/router"
	overseer "github.com/herb-go/providers/herb/overseers/routeroverseer"
	worker "github.com/herb-go/worker"
)

//RouterFactoryWorker empty router factory worker.
var RouterFactoryWorker *router.Factory

//RouterFactoryOverseer router overseer
var RouterFactoryOverseer = worker.NewOrverseer("router", &RouterFactoryWorker)

func init() {
	RouterFactoryOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(RouterFactoryOverseer)
	})
	worker.Appoint(RouterFactoryOverseer)
}
