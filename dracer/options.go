package dracer


type DracerOption func(*DracerOptions)

type DracerOptions struct {
	DockerSocket 	string
	ZipkinEndpoint  string
	Debug 			bool
}

var defaultDracerOptions = DracerOptions {
	DockerSocket: DockerSocket,
	ZipkinEndpoint: ZipkinHTTPEndpoint,
	Debug: false,
}

func WithDockerSocket(s string) DracerOption {
	return func(o *DracerOptions) {
		o.DockerSocket = s
	}
}

func WithZipkinEndpoint(s string) DracerOption {
	return func(o *DracerOptions) {
		o.ZipkinEndpoint = s
	}
}

func WithDebugValue(d bool) DracerOption {
	return func(o *DracerOptions) {
		o.Debug = d
	}
}