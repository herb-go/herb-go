package overseers

import (
	"modules/app"

	"github.com/herb-go/util"
	"github.com/herb-go/worker"
)

//MustInitOverseers init overseers.
//Panic if any error raised
func MustInitOverseers() {
	util.Must(app.PresetWorkers.Apply())
	util.Must(app.Workers.Apply())
	util.Must(worker.InitOverseers())
	util.Must(worker.TrainWorkers())
}
