package overseers

import (
	"net/http"

	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/worker/overseers/middlewareoverseer"
)

//MiddlewareWorker empty middleware worker.
var MiddlewareWorker func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

//MiddlewareOverseer middleware overseer
var MiddlewareOverseer = worker.NewOrverseer("router", &MiddlewareWorker)

func init() {
	MiddlewareOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().Apply(MiddlewareOverseer)
	})
	worker.Appoint(MiddlewareOverseer)
}
