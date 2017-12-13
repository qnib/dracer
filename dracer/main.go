// +build go1.7

package dracer

import (
	"fmt"
	"os"

	"github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"time"
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
	do DracerOptions
	dockerSocket 	string
	debug			bool
}

func NewDracer(opts ...DracerOption) DockerTracer {
	options := defaultDracerOptions
	for _, o := range opts {
		o(&options)
	}
	return DockerTracer{
		do: options,
	}
}

func (dt *DockerTracer) Run() {
	fmt.Println("huhu")
	// Create our HTTP collector.
	collector, err := zipkin.NewHTTPCollector(dt.do.ZipkinEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, dt.do.Debug, hostPort, "docker")

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

	// Create Root Span for duration of the interaction with svc1
	span := opentracing.StartSpan("Run")

	t1 := span.Tracer()
	s1 := t1.StartSpan("sub1")
	time.Sleep(1)
	s1.Finish()
	// Call the Concat Method
	span.LogEvent("Call Concat")

	// Call the Sum Method
	span.LogEvent("Call Sum")

	// Finish our CLI span
	span.Finish()

	// Close collector to ensure spans are sent before exiting.
	collector.Close()
}
