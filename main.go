package main

import (
	"os"

	"github.com/herb-go/herb-go/config"
	_ "github.com/herb-go/herb-go/modules"
	"github.com/herb-go/util/cli/app"
)

func main() {
	App := app.NewApplication(config.Config)
	App.Args = os.Args
	App.Env = app.OsEnv
	App.Modules = app.RegisteredModules
	App.Run()
}
