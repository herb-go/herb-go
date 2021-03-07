package main

import (
	"os"
	"path/filepath"

	"github.com/herb-go/herb-go/config"
	_ "github.com/herb-go/herb-go/modules"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if filepath.Base(path) == "src" {
		result, err := tools.FileExists(filepath.Join(path, "main.go"))
		if err != nil {
			panic(err)
		}
		if result {
			err = os.Chdir(filepath.Dir(path))
			if err != nil {
				panic(err)
			}
		}
	}
	App := app.NewApplication(config.Config)
	App.Args = os.Args
	App.Env = app.OsEnv
	App.Modules = app.RegisteredModules
	App.Run()
}
