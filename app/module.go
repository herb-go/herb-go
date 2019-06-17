package app

type Module interface {
	ID() string
	Cmd() string
	Help() string
	Desc() string
	Exec(*Application) error
}

type BasicModule struct {
}

func (m BasicModule) ID() string {
	return ""
}

func (m BasicModule) Cmd() string {
	return ""
}

func (m BasicModule) Help() string {
	return ""
}

func (m BasicModule) Desc() string {
	return ""
}
func (m BasicModule) Exec(*Application) error {
	return nil
}

var Modules = []Module{}

var AddModule = func(m Module) {
	cmd := m.Cmd()
	for k := range Modules {
		if cmd == Modules[k].Cmd() {
			Modules[k] = m
			return
		}
	}
	Modules = append(Modules, m)
}

var GetModule = func(cmd string) Module {
	for k := range Modules {
		if Modules[k].Cmd() == cmd {
			return Modules[k]
		}
	}
	return nil
}
