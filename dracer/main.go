// +build go1.7

package dracer

import (
	"context"
	"fmt"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/docker/docker/api/types/events"
	"strings"
)

const (
	// Endpoint to send Zipkin spans to.
	ZipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"

	// Docker socket
	DockerSocket = "unix:///var/run/docker.sock"

	// Host + port of our service.
	hostPort = "0.0.0.0:0"

	// Debug mode.
	Debug = false
	// same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
	sameSpan = true
	// make Tracer generate 128 bit traceID's for root spans.
	traceID128Bit = true
)

type DockerTracer struct {
	do 			DracerOptions
	engCli 		*client.Client
	msgs 		<-chan events.Message
	errs 		<-chan error
	sMap 		map[string]TraceSupervisor

}

func NewDracer(opts ...DracerOption) DockerTracer {
	options := defaultDracerOptions
	for _, o := range opts {
		o(&options)
	}
	return DockerTracer{
		do: options,
		sMap: make(map[string]TraceSupervisor),
	}
}

func (dt *DockerTracer) Connect() {
	var err error
	dt.engCli, err = client.NewClient(dt.do.DockerSocket, "v1.29", nil, nil)
	if err != nil {
		fmt.Printf("Could not connect docker/docker/client to '%s': %v", dt.do.DockerSocket, err)
		return
	}
	dt.msgs, dt.errs = dt.engCli.Events(context.Background(), types.EventsOptions{})

}

func (dt *DockerTracer) StartSupervisor(CntID, CntName, CntAction string) {
	ts := TraceSupervisor{
		CntID: CntID,
		CntName: CntName,
		Com: make(chan events.Message),
	}
	dt.sMap[CntID] = ts
	go ts.Run(CntAction)
}

func (dt *DockerTracer) Run() {
	dt.Connect()
	collector, _ := zipkin.NewHTTPCollector(dt.do.ZipkinEndpoint)
	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, dt.do.Debug, hostPort, "docker")

	var err error
	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}
	// Explicitly set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(tracer)


	for {
		select {
		case dMsg := <-dt.msgs:
			switch dMsg.Type {
			case "container":
				if strings.HasPrefix(dMsg.Action, "exec_") || strings.HasPrefix(dMsg.Action, "health_status") {
					continue
				}
				fmt.Printf("%s.%s: %s\n", dMsg.Type, dMsg.Action, dMsg.Actor.ID)
				if _, ok := dt.sMap[dMsg.Actor.ID]; !ok {
					dt.StartSupervisor(dMsg.Actor.ID, dMsg.Actor.Attributes["name"], dMsg.Action)
				}
				dt.sMap[dMsg.Actor.ID].Com <- dMsg
			case "service":
				//dt.RecordServiceEvent(dMsg)
			}
		case dErr := <-dt.errs:
			if dErr != nil {
				continue
			}
		}
	}
}
