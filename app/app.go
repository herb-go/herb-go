package app

type ApplicationConfig struct {
	Name          string
	Cmd           string
	Version       string
	IntroTemplate string
}
type Application struct {
	Config *ApplicationConfig
	Args   []string
	Envs   []string
}

func (a *Application) Run() {

}

var Config = &ApplicationConfig{
	Name:          "Herb-go cli tool",
	Cmd:           "herb-go",
	Version:       "0.1",
	IntroTemplate: "",
}

func Run(config *ApplicationConfig, args []string, envs []string) {
	var app = &Application{
		Config: config,
		Args:   args,
		Envs:   envs,
	}
	app.Run()
}
