package app

import (
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
