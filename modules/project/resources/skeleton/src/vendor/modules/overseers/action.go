package overseers

import (
	"github.com/herb-go/herb/middleware/action"
	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/providers/herb/overseers/actionoverseer"
)

//ActionWorker empty cache worker.
var ActionWorker = action.New(nil)

//ActionOverseer cache overseer
var ActionOverseer = worker.NewOrverseer("action", &ActionWorker)

func init() {
	ActionOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().Apply(ActionOverseer)
	})
	worker.Appoint(ActionOverseer)
}
