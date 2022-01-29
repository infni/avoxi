package api

type IPAuthrorizeConfig struct {
	// list version and exit
	DisplayVesion bool

	// the gRPC port to listen on
	GrpcPort string

	// the http port the service will listen on
	HttpPort string
}
