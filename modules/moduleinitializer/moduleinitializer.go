package moduleinitializer

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Initializer struct {
	app.BasicModule
	EnvList     string
	Prefix      string
	SlienceMode bool
}

func (m *Initializer) ID() string {
	return "github.com/herb-go/herb-go/modules/moduleinitializer"
}

func (m *Initializer) Cmd() string {
	return "moduleinitializer"
}

func (m *Initializer) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s moduleinitializer [name].
Create new module initializer code.
File below will be created:
	src/modules/[name]/initializer.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Initializer) Desc(a *app.Application) string {
	return "Create new module initializer code"
}
func (m *Initializer) Group(a *app.Application) string {
	return "Application"
}
func (m *Initializer) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().StringVar(&m.EnvList, "envlist", "env1,env2", "Envlist joined with \",\"")
	m.FlagSet().StringVar(&m.Prefix, "prefix", "Modules", "Envs prefix")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *Initializer) Question(a *app.Application) error {
	return nil
}
func (m *Initializer) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		a.PrintModuleHelp(m)
		return nil
	}
	n, err := name.New(true, args...)
	if err != nil {
		return err
	}
	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	result, err := tools.FileExists(mp, n.LowerPath("init.go"))
	if err != nil {
		return err
	}
	if !result {
		return errors.New("module not found")
	}
	err = m.Question(a)
	if err != nil {
		return err
	}
	var envlist = []*name.Name{}
	envs := strings.Split(m.EnvList, ",")
	if len(envs) == 0 {
		return fmt.Errorf("env list \"%s\" is not available", m.EnvList)
	}
	for _, v := range envs {
		if v == "" {
			return fmt.Errorf("env list \"%s\" is not available", m.EnvList)
		}
		e, err := name.New(false, v)
		if err != nil {
			return err
		}
		envlist = append(envlist, e)
	}

	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/moduleinitializer/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n, envlist)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Module initializer  \"%s\" created.\n", n.LowerWithParentDotSeparated)
		return nil
	})
	err = task.ErrosIfAnyFileExists()
	if err != nil {
		return err
	}
	ok, err := task.ConfirmIf(a, !m.SlienceMode)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return task.Exec()

}

func (m *Initializer) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name, envlist []*name.Name) error {
	modelmodule := n.LowerWithParentPath
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath("initializer.go")): "initializer.go.tmpl",
	}
	quotedenvs := make([]string, len(envlist))
	for k := range envlist {
		quotedenvs[k] = "\"" + m.Prefix + "." + envlist[k].Raw + "\""
	}
	envsempty := make([]string, len(envlist))
	for k := range envlist {
		envsempty[k] = envlist[k].Lower + " == \"\""
	}
	data := map[string]interface{}{
		"Name":      n,
		"Envlist":   envlist,
		"Module":    modelmodule,
		"Prefix":    m.Prefix,
		"EnvOr":     strings.Join(envsempty, "||"),
		"EnvParams": strings.Join(quotedenvs, " , "),
	}
	return task.RenderFiles(filesToRender, data)
}

var InitializerModule = &Initializer{}

func init() {
	app.Register(InitializerModule)
}
