package main

import (
	"context"
	"fmt"
	"ipauthorize/api"
	"os"

	"google.golang.org/grpc"
)

func main() {
	os.Exit(doCall())
}

func doCall() int {
	conn, err := grpc.Dial("127.0.0.1:9079", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Client can't connect. %s", err.Error())
	}
	defer conn.Close()

	client := api.NewIpAuthorizeClient(conn)

	resp, respErr := client.Health(context.Background(), &api.HealthRequest{})
	if respErr != nil {
		fmt.Printf("Health call failed. %s", err.Error())
		return 1
	}

	fmt.Printf("SUccess! '%s'", resp.Now)

	return 0
}
