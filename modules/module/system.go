package module

import (
	"fmt"

	"github.com/herb-go/util/cli/app"
)

type Systems struct {
	*Module
}

func (m *Systems) ID() string {
	return "github.com/herb-go/herb-go/modules/systems"
}

func (m *Systems) Cmd() string {
	return "systems"
}

func (m *Systems) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s systems [name].
Create module with default location "systems".
File below will be created:
	src/modules/systems/<name>/init.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Systems) Desc(a *app.Application) string {
	return "Create new systems module code"
}

var SystemModule = &Systems{
	Module: &Module{
		DefaultLocation: "systems",
	},
}

func init() {
	app.Register(SystemModule)
}
