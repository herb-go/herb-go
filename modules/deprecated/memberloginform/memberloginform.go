package memberloginform

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/herb-go/util/cli/name"

	"github.com/herb-go/herb-go/modules/project"
	"github.com/herb-go/util/cli/app"
	"github.com/herb-go/util/cli/app/tools"
)

type LoginForm struct {
	app.BasicModule
	FormID          string
	Caseinsensitive bool
	AccountKeyword  string
	SlienceMode     bool
}

func (m *LoginForm) ID() string {
	return "github.com/herb-go/herb-go/modules/memberloginform"
}

func (m *LoginForm) Cmd() string {
	return "memberloginform"
}

func (m *LoginForm) Help(a *app.Application) string {
	m.Init(a, &[]string{})
	help := `Usage %s memberloginform [name].
Create new member login form and action.
File below will be created:
	src/modules/[name]/[id]/forms/login.form
	src/modules/[name]/[id]/actions/login.action
`
	return fmt.Sprintf(help, a.Config.Cmd)
}

func (m *LoginForm) Desc(a *app.Application) string {
	return "Create new member login form and action"
}
func (m *LoginForm) Group(a *app.Application) string {
	return "Auth"
}
func (m *LoginForm) Init(a *app.Application, args *[]string) error {
	if m.FlagSet().Parsed() {
		return nil
	}
	m.FlagSet().BoolVar(&m.SlienceMode, "s", false, "Slience mode")
	m.FlagSet().StringVar(&m.FormID, "id", "login", "Form id")
	m.FlagSet().BoolVar(&m.Caseinsensitive, "caseinsensitive", false, "Username case insensitive mode")
	m.FlagSet().StringVar(&m.AccountKeyword, "accountkeyword", "", "Member account keyword.UID will be used for login if no account keyword given.")
	err := m.FlagSet().Parse(*args)
	if err != nil {
		return err
	}
	*args = m.FlagSet().Args()
	return nil
}
func (m *LoginForm) Question(a *app.Application) error {
	return nil
}
func (m *LoginForm) Exec(a *app.Application, args []string) error {
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
	formid, err := name.New(false, m.FormID)
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
		return errors.New("member modules not found")
	}
	app, err := tools.FindLib(a.Getenv("GOPATH"), "github.com/herb-go/herb-go")
	if err != nil {
		return err
	}

	task := tools.NewTask(filepath.Join(app, "/modules/memberloginform/resources"), a.Cwd)

	err = m.Render(a, a.Cwd, mp, task, n, formid)
	if err != nil {
		return err
	}
	task.AddJob(func() error {
		a.Printf("Member login form  \"%s\" created.\n", n.LowerWithParentDotSeparated)
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

func (m *LoginForm) Render(a *app.Application, appPath string, mp string, task *tools.Task, n *name.Name, fid *name.Name) error {
	filesToRender := map[string]string{
		filepath.Join(mp, n.LowerPath(fid.Lower, "forms", fid.Lower+"form.go")):     "form.go.tmpl",
		filepath.Join(mp, n.LowerPath(fid.Lower, "actions", fid.Lower+"action.go")): "action.go.tmpl",
	}
	data := map[string]interface{}{
		"Name":            n,
		"FormID":          fid,
		"Caseinsensitive": m.Caseinsensitive,
		"AccountKeyword":  m.AccountKeyword,
	}
	return task.RenderFiles(filesToRender, data)
}

var LoginFormModule = &LoginForm{}

func init() {
	app.Register(LoginFormModule)
}
