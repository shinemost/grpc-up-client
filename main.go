package main

import (
	"context"
	"github.com/shinemost/grpc-up-client/interceptor"
	"log"
	"time"

	"github.com/shinemost/grpc-up-client/clients"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.OrderUnaryClientInterceptor),
		grpc.WithStreamInterceptor(interceptor.ClientStreamInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	clientDeadLine := time.Now().Add(time.Second * 2)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadLine)
	defer cancel()

	clients.AddOrder(conn, ctx)
	//clients.UpdateOrders(conn, ctx)
	//clients.SearchOrders(conn, ctx)
	//clients.ProcessOrders(conn, ctx)

}
