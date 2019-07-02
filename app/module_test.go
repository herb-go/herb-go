package app

import "testing"

type testModule struct {
	BasicModule
	id     string
	cmd    string
	help   string
	desc   string
	output string
}

func (m *testModule) ID() string {
	return m.id
}
func (m *testModule) Cmd() string {
	return m.cmd
}
func (m *testModule) Help() string {
	return m.help
}
func (m *testModule) Desc() string {
	return m.desc
}
func (m *testModule) Exec(a *Application, args []string) error {
	if m.output != "" {
		a.Println(m.output)
	}
	return nil
}

func newTestModule(cmd string) *testModule {
	return &testModule{
		cmd: cmd,
	}
}
func TestBasicModule(t *testing.T) {
	type testBasicModule struct {
		BasicModule
	}
	tm := &testBasicModule{}
	if tm.Cmd() != "" {
		t.Fatal(tm)
	}
	if tm.Desc() != "" {
		t.Fatal(tm)
	}
	if tm.ID() != "" {
		t.Fatal(tm)
	}
	if tm.Help() != "" {
		t.Fatal(tm)
	}
	if err := tm.Exec(NewApplication(nil), nil); err != nil {
		t.Fatal(tm)
	}
}
func TestModule(t *testing.T) {
	modules := NewModules()
	tmtest := newTestModule("test")
	tmtest2 := newTestModule("test2")
	modules.Add(tmtest)
	modules.Add(tmtest2)
	modules.Add(tmtest2)
	m := modules.Get("test")
	if m != tmtest {
		t.Fatal(m)
	}
	m = modules.Get("test2")
	if m != tmtest2 {
		t.Fatal(m)
	}
	m = modules.Get("notexist")
	if m != nil {
		t.Fatal(m)
	}
}

func TestRegisterModule(t *testing.T) {
	defer func() {
		RegisteredModules = NewModules()
	}()
	tmtest := newTestModule("test")
	tmtest2 := newTestModule("test2")
	Register(tmtest)
	Register(tmtest2)
	m := RegisteredModules.Get("test")
	if m != tmtest {
		t.Fatal(m)
	}
	m = RegisteredModules.Get("test2")
	if m != tmtest2 {
		t.Fatal(m)
	}
	m = RegisteredModules.Get("notexist")
	if m != nil {
		t.Fatal(m)
	}
}
