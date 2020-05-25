package overseers

import (
	"github.com/herb-go/herb/middleware/router"
	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/worker/overseers/routeroverseer"
)

//RouterWorker empty router worker.
var RouterWorker router.Router

//RouterOverseer router overseer
var RouterOverseer = worker.NewOrverseer("router", &RouterWorker)

func init() {
	RouterOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().Apply(RouterOverseer)
	})
	worker.Appoint(RouterOverseer)
}
