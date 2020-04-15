package overseers

import (
	"github.com/herb-go/herb/middleware/action"
	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/worker/overseers/actionoverseer"
)

//ActionWorker empty cache worker.
//You should set your empter worker instead of worker.DummyTeam
var ActionWorker = action.New(nil)

//ActionOverseer cache overseer
var ActionOverseer = worker.NewOrverseer("action", ActionWorker)

func init() {
	ActionOverseer.WithInitFunc(func() error {
		return overseer.New().Apply(ActionOverseer)
	})
}
