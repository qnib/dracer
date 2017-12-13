package dracer

import (
	"log"

	"github.com/docker/docker/api/types/events"
	"github.com/opentracing/opentracing-go"
)

type TraceSupervisor struct {
	CntID 	string 			 	// ContainerID
	CntName string			 	// sanatized name of container
	Com 	chan events.Message // Channel to communicate with goroutine
}

func (ts *TraceSupervisor) Run() {
	log.Printf("[II] Start listener for: '%s' [%s]", ts.CntName, ts.CntID[:12])
	span := opentracing.StartSpan(ts.CntName)
	span.LogEvent("create")
	defer span.Finish()
	for {
		select {
		case msg := <-ts.Com:
			span.LogEvent(msg.Action)
			switch msg.Action {
			case "destroy":
				return
			}
		}
	}
}