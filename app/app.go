package app

import (
	"fmt"
	"io"
	"os"
)

type ApplicationConfig struct {
	Name          string
	Cmd           string
	Version       string
	IntroTemplate string
}
type Application struct {
	Config  *ApplicationConfig
	Args    []string
	Envs    []string
	Cwd     string
	Modules *Modules
	Stdout  io.Writer
	Stdin   io.Reader
}

func (a *Application) Print(args ...interface{}) (n int, err error) {
	return fmt.Fprint(a.Stdout, args...)
}

func (a *Application) Println(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(a.Stdout, args...)
}

func (a *Application) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(a.Stdout, format, args...)
}
func (a *Application) Run() {

}

func NewApplication(config *ApplicationConfig) *Application {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &Application{
		Config:  config,
		Cwd:     cwd,
		Modules: NewModules(),
		Stdout:  os.Stdout,
		Stdin:   os.Stdin,
	}
}

var Config = &ApplicationConfig{
	Name:          "Herb-go cli tool",
	Cmd:           "herb-go",
	Version:       "0.1",
	IntroTemplate: "",
}
