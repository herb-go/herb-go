package main

import (
	"modules/app"

	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
)

func loadConfigs() {
	//Uncomment next line to print config loading log .
	//config.Debug = true
	config.Lock.RLock()
	app.LoadConfigs()
	config.Lock.RUnlock()
}
func initModules() {
	util.InitModulesOrderByName()
	//Put Your own init code here.
}

//Main app run func.
var run = func() {
	//Replace next line "panic(errFuncWhenRunFuncNotRewrited)" with your own app run function
	panic(errFuncWhenRunFuncNotRewrited)
}

func start() {
	defer util.Recover()
	util.ApplieationLock.Lock()
	defer util.ApplieationLock.Unlock()
	util.UpdatePaths()
	util.MustChRoot()
	loadConfigs()
	initModules()
	app.Development.NotTestingOrPanic()
	util.RegisterDataFolder() //Auto created appdata folder if not exists
	util.MustLoadRegisteredFolders()
	app.Development.InitializeAndPanicIfNeeded()
	run()
}
func main() {
	// Set app root path.
	//Default rootpah is "src/../"
	//You can set os env  "HerbRoot" to overwrite this setting while starting app.
	// util.RootPath = ""
	start()
}
