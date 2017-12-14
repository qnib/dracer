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

func (ts *TraceSupervisor) Run(CntAction string) {
	log.Printf("[II] Start listener for: '%s' [%s]", ts.CntName, ts.CntID[:12])
	span := opentracing.StartSpan(ts.CntName)
	if CntAction == "create" {
		span.LogEvent("create")
	} else {
		span.LogEvent("discovered")
	}
	for {
		select {
		case msg := <-ts.Com:
			span.LogEvent(msg.Action)
			switch msg.Action {
			case "destroy":
				span.Finish()
			}
		}
	}
}