package help

import "github.com/herb-go/herb-go/app"

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
