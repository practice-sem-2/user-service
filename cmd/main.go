package main

import (
	"flag"
	"fmt"
	"github.com/practice-sem-2/user-service/internal/pb"
	"github.com/practice-sem-2/user-service/internal/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	var host string
	var port int

	flag.IntVar(&port, "port", 80, "port on which server will be started")
	flag.StringVar(&host, "host", "0.0.0.0", "host on which server will be started")

	flag.Parse()

	address := fmt.Sprintf("%s:%d", host, port)

	lis, err := net.Listen("tcp", address)

	log.Printf("Start listening on %s", address)

	if err != nil {
		log.Fatalf("can't listen to address: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, server.NewUserServer())

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Grpc serving error: %s", err.Error())
	}
}
