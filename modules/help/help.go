package help

import "github.com/herb-go/util/cli/app"

type Help struct {
	app.BasicModule
}

func (m *Help) Cmd() string {
	return "help"
}

var Module = &Help{}

func init() {
	app.Register(Module)
}
