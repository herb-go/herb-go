package config

import "github.com/herb-go/util/cli/app"

var Config = app.NewApplicationConfig()

func InitConcif() {
	Config.Name = "Herb-go cli tool"
	Config.Cmd = "herb-go"
	Config.Version = "0.1"
	Config.IntroTemplate = "{{.Config.Name}} Version {{.Config.Version}}\nCli tool to create herb-go app.\nType \"{{.Config.Cmd}} {{.HelpModuleCmd}}\" to get help."
}
