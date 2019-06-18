package project

import "github.com/herb-go/herb-go/app"

type Project struct {
	app.BasicModule
}

func (m *Project) ID() string {
	return "github.com/herb-go/herb-go/modules/project"
}

func (m *Project) Cmd() string {
	return "new"
}

func (m *Project) Help() string {
	return ""
}

func (m *Project) Desc() string {
	return ""
}
func (m *Project) Exec(*app.Application) error {
	return nil
}

var Module = &Project{}

func init() {
	app.Register(Module)
}
