package overseers

import (
	overseer "github.com/herb-go/herb-drivers/overseers/identifieroverseer"
	"github.com/herb-go/herb/user/httpuser"
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
