package persist

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/herb-go/herb-go/modules/project"

	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type Version struct {
	app.BasicModule
	SlienceMode bool
	SemverMode  bool
}

type renderData struct {
	Year  string
	Month string
	Day   string
}

func (m *Version) ID() string {
	return "github.com/herb-go/herb-go/modules/version"
}

func (m *Version) Cmd() string {
	return "version"
}

func (m *Version) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s version.
Create version module.
Default name is "persist".
File below will be created:
	src/modules/version/version.go
	src/modules/version.go
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *Version) Desc(a *app.Application) string {
	return "Create version module"
}
func (m *Version) Group(a *app.Application) string {
	return "Misc"
}
func (m *Version) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().BoolVar(&m.SemverMode, "semver", false, "Semantic version mode")

	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}

func (m *Version) Exec(a *app.Application, args []string) error {
	err := m.Init(a, &args)
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
	task := tools.NewTask(filepath.Join(app, "/modules/version/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Version module created.\n")
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

func (m *Version) Render(a *app.Application, appPath string, mp string, task *tools.Task) error {
	var tmpl = ""
	if m.SemverMode {
		tmpl = "semver.go.tmpl"
	} else {
		tmpl = "version.go.tmpl"
	}
	filesToRender := map[string]string{
		filepath.Join(mp, "version.go"):            "versioninit.go.tmpl",
		filepath.Join(mp, "version", "version.go"): tmpl,
	}
	now := time.Now()
	data := renderData{
		Year:  now.Format("2006"),
		Month: now.Format("1"),
		Day:   now.Format("2"),
	}
	return task.RenderFiles(filesToRender, data)
}

var VersionModule = &Version{}

func init() {
	app.Register(VersionModule)
}
