package overseers

import (
	"github.com/herb-go/datasource/sql/db"
	overseer "github.com/herb-go/datasource-drivers/overseers/dboverseer"
	worker "github.com/herb-go/worker"
)

//DatabaseWorker empty database worker.
var DatabaseWorker = db.New()

//DatabaseOverseer database overseer
var DatabaseOverseer = worker.NewOrverseer("database", &DatabaseWorker)

func init() {
	DatabaseOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(DatabaseOverseer)
	})
	worker.Appoint(DatabaseOverseer)
}
