package help

import "github.com/herb-go/herb-go/app"

type Help struct {
	app.BasicModule
}

var Module = &Help{}

func init() {
	app.Register(Module)
}
