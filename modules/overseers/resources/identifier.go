package overseers

import (
	"github.com/herb-go/herb/user/httpuser"
	overseer "github.com/herb-go/providers/herb/overseers/identifieroverseer"
	worker "github.com/herb-go/worker"
)

//IdentifierWorker empty identifier worker.
var IdentifierWorker httpuser.Identifier

//IdentifierOverseer identifier overseer
var IdentifierOverseer = worker.NewOrverseer("identifier", &IdentifierWorker)

func init() {
	IdentifierOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(IdentifierOverseer)
	})
	worker.Appoint(IdentifierOverseer)
}
