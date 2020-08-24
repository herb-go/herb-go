package overseers

import (
	"github.com/herb-go/usersystem"
	overseer "github.com/herb-go/usersystem-drivers/overseers/usersystemdirectivefactoryoverseer"
	worker "github.com/herb-go/worker"
)

//UserSystemDirectiveFactoryWorker usersystem directive factory worker.
var UserSystemDirectiveFactoryWorker func(loader func(v interface{}) error) (usersystem.Directive, error)

//UserSystemDirectiveFactoryOverseer usersystem directive factory overseer
var UserSystemDirectiveFactoryOverseer = worker.NewOrverseer("memberdirectivefactory", &UserSystemDirectiveFactoryWorker)

func init() {
	UserSystemDirectiveFactoryOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(UserSystemDirectiveFactoryOverseer)
	})
	worker.Appoint(UserSystemDirectiveFactoryOverseer)
}
