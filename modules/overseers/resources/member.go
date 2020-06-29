package overseers

import (
	"github.com/herb-go/member"
	overseer "github.com/herb-go/providers/herb/overseers/memberoverseer"
	worker "github.com/herb-go/worker"
)

//MemberWorker empty member worker.
var MemberWorker = member.New()

//MemberOverseer member overseer
var MemberOverseer = worker.NewOrverseer("member", &MemberWorker)

func init() {
	MemberOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(MemberOverseer)
	})
	worker.Appoint(MemberOverseer)
}
