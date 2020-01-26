package loggers

import (
	"modules/app"
	"os"
	"os/signal"

	"github.com/herb-go/logger"

	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "100loggers"

//MyLogger my logger
// var MyLogger = logger.PrintLogger.SubLogger().SetID("mylogger")

var reopenSignals []os.Signal

func reopenLoggers() {
	logger.ResetBuiltinLoggers()
	// MyLogger.Reopen()

}
func listenForReopenLoggers() {
	var signalChan = make(chan os.Signal)
	signal.Notify(signalChan, reopenSignals...)
	go func() {
		<-util.QuitChan()
		signal.Stop(signalChan)
	}()
	go func() {
		for {
			select {
			case <-signalChan:
				go reopenLoggers()
			case <-util.QuitChan():
				return
			}
		}
	}()
}

func init() {
	util.RegisterModule(ModuleName, func() {
		util.Must(app.Loggers.ApplyToBuiltinLoggers())
		util.ErrorLogger = logger.Error
		// util.Must(app.Loggers.ApplyTo(MyLogger, "mylogger"))
		if app.Development.Debug {
			logger.EnableDevelopmengLoggers()
		}
		listenForReopenLoggers()
	})
}
