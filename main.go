package main

import (
	"os"
	"log"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"
	"github.com/qnib/dracer/dracer"

)

var (
	dockerSocketFlag = cli.StringFlag{
		Name:  "docker-socket",
		Value: dracer.DockerSocket,
		Usage: "Docker host to connect to.",
		EnvVar: "DOXY_DOCKER_SOCKET",
	}
	zipkinEndpointFlag = cli.StringFlag{
		Name:  "zipkin-endpoint",
		Value: dracer.ZipkinHTTPEndpoint,
		Usage: "Zipkin endpoint to connect to.",
		EnvVar: "DRACER_ZIPKIN_ENDPOINT",
	}
	debugFlag = cli.BoolFlag{
		Name: "debug",
		Usage: "Print proxy requests",
		EnvVar: "DOXY_DEBUG",
	}
)

func EvalOptions(cfg *config.Config) (po []dracer.DracerOption) {
	zipkinEndpoint, _ := cfg.String("zipkin-endpoint")
	po = append(po, dracer.WithZipkinEndpoint(zipkinEndpoint))
	dockerSock, _ := cfg.String("docker-socket")
	po = append(po, dracer.WithDockerSocket(dockerSock))
	debug, _ := cfg.Bool("debug")
	po = append(po, dracer.WithDebugValue(debug))
	return
}

func RunApp(ctx *cli.Context) {
	log.Printf("[II] Start Version: %s", ctx.App.Version)
	cfg := config.NewConfig([]config.Provider{config.NewCLI(ctx, true)})
	po := EvalOptions(cfg)
	p := dracer.NewDracer(po...)
	p.Run()
}

func main() {
	app := cli.NewApp()
	app.Name = "Tracer for docker events comming from the docker engine."
	app.Usage = "dracer [options]"
	app.Version = "0.0.0"
	app.Flags = []cli.Flag{
		debugFlag,
		dockerSocketFlag,
		zipkinEndpointFlag,
	}
	app.Action = RunApp
	app.Run(os.Args)
}