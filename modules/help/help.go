package help

import "github.com/herb-go/herb-go/app"

type Help struct {
	app.BasicModule
}

func init() {
	app.AddModule(&Help{})
}
