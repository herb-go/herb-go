package overseers

import (
	"github.com/herb-go/member"
	overseer "github.com/herb-go/providers/herb/overseers/memberoverseer"
	worker "github.com/herb-go/worker"
)

//MemberWorker empty cache worker.
var MemberWorker = member.New()

//MemberOverseer cache overseer
var MemberOverseer = worker.NewOrverseer("member", &MemberWorker)

func init() {
	MemberOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().Apply(MemberOverseer)
	})
	worker.Appoint(MemberOverseer)
}
