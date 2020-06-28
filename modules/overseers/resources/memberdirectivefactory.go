package overseers

import (
	"github.com/herb-go/member"
	overseer "github.com/herb-go/providers/herb/overseers/memberdirectivefactoryoverseer"
	worker "github.com/herb-go/worker"
)

//MemberDirectiveFactoryWorker member directive factory worker.
var MemberDirectiveFactoryWorker func(loader func(v interface{}) error) (member.Directive, error)

//MemberDirectiveFactoryOverseer member directive factory overseer
var MemberDirectiveFactoryOverseer = worker.NewOrverseer("memberdirectivefactory", &MemberDirectiveFactoryWorker)

func init() {
	MemberOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(MemberDirectiveFactoryOverseer)
	})
	worker.Appoint(MemberDirectiveFactoryOverseer)
}
