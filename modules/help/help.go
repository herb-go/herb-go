package help

import (
	"fmt"

	"github.com/herb-go/util/cli/app"
)

type Help struct {
	app.BasicModule
}

func (m *Help) Cmd() string {
	return "help"
}

func (m *Help) ID() string {
	return "github.com/herb-go/herb-go/modules/help"
}

func (m *Help) Help(a *app.Application) string {
	help := "Usage %s help [command]\r\n"
	help += "Command list:\r\n"
	for _, v := range *a.Modules {
		help += fmt.Sprintf("  %s - %s\r\n", v.Cmd(), v.Desc(a))
	}
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Help) Desc(a *app.Application) string {
	return "Show module help"
}

func (m *Help) Exec(a *app.Application, args []string) error {
	if len(args) != 1 {
		a.PrintModuleHelp(m)
		return nil
	}
	cmd := args[0]
	module := a.Modules.Get(cmd)
	if module == nil {
		a.Printf("Module \"%s\" not found.\n", cmd)
		a.Println(m)
		return nil
	}
	a.PrintModuleHelp(module)
	return nil
}

var Module = &Help{}

func init() {
	app.Register(Module)
}
