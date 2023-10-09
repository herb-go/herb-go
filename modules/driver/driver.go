package driver

import (
	"fmt"
	"strings"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
)

type Driver struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Driver) ID() string {
	return "github.com/herb-go/herb-go/modules/driver"
}

func (m *Driver) Cmd() string {
	return "driver"
}

func (m *Driver) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s driver.
Add driver file.
	File below will be created:
	src/modules/driver/[name].go

Supported drivers:
	%s
`
	return fmt.Sprintf(help, a.Config.Cmd, strings.Join(GetRegisteredDrivers(), " , "))
}

func (m *Driver) Desc(a *app.Application) string {
	return "Add driver file."
}
func (m *Driver) Group(a *app.Application) string {
	return "Application"
}
func (m *Driver) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Driver) Question(a *app.Application) error {
	return nil
}
func (m *Driver) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}

	if len(args) != 1 {
		a.PrintModuleHelp(m)
		return nil
	}

	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	if f := DriverInitFuncs[args[0]]; f != nil {
		return f(a, a.Cwd, mp, m.SlienceMode)
	}
	a.PrintModuleHelp(m)
	return nil

}

var DriverModule = &Driver{}

func init() {
	app.Register(DriverModule)
}
