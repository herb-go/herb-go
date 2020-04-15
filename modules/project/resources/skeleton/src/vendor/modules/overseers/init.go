package overseers

import (
	"github.com/herb-go/util"
	"github.com/herb-go/worker"
)

//MustInitOverseers init overseers.
//Panic if any error raised
func MustInitOverseers() {
	util.Must(worker.InitOverseers())
}
