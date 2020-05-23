package overseers

import (
	"github.com/herb-go/herb/model/sql/db"
	worker "github.com/herb-go/worker"
	overseer "github.com/herb-go/worker/overseers/dboverseer"
)

//DatabaseWorker empty cache worker.
var DatabaseWorker = db.New()

//DatabaseOverseer cache overseer
var DatabaseOverseer = worker.NewOrverseer("database", &DatabaseWorker)

func init() {
	DatabaseOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().Apply(DatabaseOverseer)
	})
	worker.Appoint(DatabaseOverseer)
}
