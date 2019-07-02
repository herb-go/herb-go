package main

import (
	"os"

	"github.com/herb-go/herb-go/app"
	_ "github.com/herb-go/herb-go/modules"
)

func main() {
	App := app.NewApplication(app.Config)
	App.Args = os.Args
	App.Env = app.OsEnv
	App.Modules = app.RegisteredModules
	App.Run()
}
