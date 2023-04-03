package index

import (
	"task/controller/base"
)

func Index(ch *base.ControllerHandle) {
	//_, _ = io.WriteString(w, "index")
	ch.ReturnData()
	return
}
