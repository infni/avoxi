package main

import (
	"context"
	"flag"
	"fmt"
	golog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ipauthorize/api"
	"ipauthorize/internal/pkg/countycodes"
	"ipauthorize/internal/pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var version = "defined at build time"

const (
	defaultGrpcPort = "9079"
	defaulthttpPort = "9080"
)

func main() {

	// get configuration
	cfg := &api.IPAuthrorizeConfig{}
	flag.StringVar(&cfg.GrpcPort, "grpc-port", defaultGrpcPort, "gRPC server endpoint (default: "+defaultGrpcPort+")")
	flag.StringVar(&cfg.HttpPort, "http-port", defaulthttpPort, "http port to bind to (defaults to "+defaulthttpPort+")")
	flag.BoolVar(&cfg.DisplayVesion, "version", false, "display the version and exit")

	flag.Parse()

	if cfg.DisplayVesion {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if execute(cfg) {
		fmt.Print("Exited\n")
		os.Exit(0)
	}

	// the error has previously been written out.  Just exit.
	os.Exit(1)
}

func execute(cfg *api.IPAuthrorizeConfig) bool {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logFilename := fmt.Sprintf("IpAuthorize_%s.log", time.Now().Format("20060102150405"))

	logFilePointer, logError := os.Create(logFilename)
	if logError != nil {
		fmt.Printf("Failed to open log file '%v'. Error: %v\n", logFilename, logError)
		return false
	}

	logger := log.NewFileLogger(golog.New(logFilePointer, "", 0))
	defer logFilePointer.Close()

	addlInfo := log.AddlInfo{
		"cfg":     cfg,
		"version": version,
	}

	// register service
	g := grpc.NewServer()

	comparer := countycodes.NewContryCodeComparer("GeoLite2-Country.mmdb")
	api.RegisterIpAuthorizeServer(g, api.NewIPAuthorizev1(comparer))

	// run the gRPC server in the background
	if listen, err := net.Listen("tcp", ":"+cfg.GrpcPort); err == nil {
		go func() {
			fmt.Println("starting gRPC server...")
			if err := g.Serve(listen); err != nil {
				fmt.Printf("failed to serve: %s", err)
			}
		}()
	} else {
		logger.LogCritical(fmt.Sprintf(" : gRPC Listener aborted : %v", err.Error()), addlInfo)
		return false
	}

	// start the http server
	server := runtime.NewServeMux()
	if err := api.RegisterIpAuthorizeHandlerFromEndpoint(ctx, server, "localhost:"+cfg.GrpcPort, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		fmt.Printf("Failed to register IpAuhtorize as gRPC enpoint for MUX : %v\n", err)
		return false
	}

	fmt.Println("starting HTTP server...")

	httpServer := &http.Server{Addr: ":" + cfg.HttpPort, Handler: server}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogCritical(fmt.Sprintf(" : HTTP Proxy Server aborted : %v", err.Error()), addlInfo)
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(
		sigs,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	s := <-sigs
	addlInfo["sigterm"] = s.String()
	logger.Log("Shutting down from terminal signal.", addlInfo)

	go func() {
		localAddlInfo := addlInfo
		s := <-sigs
		localAddlInfo["sigterm"] = s.String()
		logger.Log("Additional signals recieved.", localAddlInfo)
	}()

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	fmt.Println("Stopping HTTP server gracefully")
	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		logger.LogCritical(fmt.Sprintf("Failed to shutdown http server : %v", err.Error()), addlInfo)
	}

	fmt.Println("Stopping gRPC server gracefully")
	g.GracefulStop()

	cancel()

	return true
}
