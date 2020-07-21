package protecter

import (
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/httpinfo"
	"github.com/herb-go/herb-go/modules/overseers"
	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Protecter struct {
	app.BasicModule
	SlienceMode bool
}

func (m *Protecter) ID() string {
	return "github.com/herb-go/herb-go/modules/protecter"
}

func (m *Protecter) Cmd() string {
	return "protecter"
}

func (m *Protecter) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s protecter <name>.
Create protecter module and config files.
Default name is "protecter".
File below will be created:
	config/<name>.toml
	system/configskeleton/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/init.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Protecter) Desc(a *app.Application) string {
	return "Create protecter module and config files."
}
func (m *Protecter) Group(a *app.Application) string {
	return "Auth"
}
func (m *Protecter) Init(a *app.Application, args *[]string) error {
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
func (m *Protecter) Question(a *app.Application) error {
	return nil
}
func (m *Protecter) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
	if err != nil {
		return err
	}
	var n *name.Name

	if len(args) == 0 {
		fmt.Println("No protecter module name given.\"protecter\" is used")
		n, err = name.New(true, "protecter")
	} else {
		n, err = name.New(true, args...)
	}
	if err != nil {
		return err
	}

	mp, err := project.GetModuleFolder(a.Cwd)
	if err != nil {
		return err
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}
	err = overseers.OverseerModule.Exec(a, []string{"-s", "authenticator"})
	if err != nil {
		return err
	}
	result, err := tools.FileExists(mp, "httpinfo", "init.go")
	if err != nil {
		return err
	}
	if !result {
		err = httpinfo.HTTPInfoModule.Exec(a, []string{"-s"})
		if err != nil {
			return err
		}

	}
	task := tools.NewTask(filepath.Join(app, "/modules/protecter/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Protecter  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *Protecter) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name) error {
	err := tools.CopyIfNotExist(filepath.Join(task.SrcFolder, "hired.authenticator.go.tmpl"), mp, "hired", "authenticator.go")
	if err != nil {
		return err
	}
	filesToRender := map[string]string{
		filepath.Join("config", n.LowerWithParentDotSeparated+".toml"):                   "protecter.toml.tmpl",
		filepath.Join("system", "configskeleton", n.LowerWithParentDotSeparated+".toml"): "protecter.toml.tmpl",
		filepath.Join(mp, n.LowerPath(n.Lower+".go")):                                    "protecter.go.tmpl",
		filepath.Join(mp, "app", n.LowerWithParentDotSeparated+".go"):                    "app.protecter.go.tmpl",
	}
	return task.RenderFiles(filesToRender, n)
}

var ProtecterModule = &Protecter{}

func init() {
	app.Register(ProtecterModule)
}
