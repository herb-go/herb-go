package overseers

import (
	"github.com/herb-go/protecter/authenticator"
	overseer "github.com/herb-go/providers/herb/overseers/authenticatorfactoryoverseer"
	worker "github.com/herb-go/worker"
)

//AuthenticatorFactoryWorker authenticator factory worker.
var AuthenticatorFactoryWorker authenticator.AuthenticatorFactory

//AuthenticatorFactoryOverseer authenticator factory overseer
var AuthenticatorFactoryOverseer = worker.NewOrverseer("authenticatorfactory", &AuthenticatorFactoryWorker)

func init() {
	AuthenticatorFactoryOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(AuthenticatorFactoryOverseer)
	})
	worker.Appoint(AuthenticatorFactoryOverseer)
}
