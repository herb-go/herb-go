package main

import (
	"os"

	"github.com/herb-go/herb-go/app"
	_ "github.com/herb-go/herb-go/modules"
)

func main() {
	app.Run(app.Config, os.Args, os.Environ())
}
