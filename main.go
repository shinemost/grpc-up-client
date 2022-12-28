package main

import (
	"context"
	"github.com/shinemost/grpc-up-client/clients"
	"github.com/shinemost/grpc-up-client/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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
	//clientDeadLine := time.Now().Add(time.Second * 2)
	//过期时间与截止时间不同，timeout每个请求单独设置，可能客户端发起一个请求会调用多个服务，那么过期时间会叠加，而截止时间deadline则是绝对时间
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), clientDeadLine)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Cancelling the RPC
	cancel()
	clients.AddOrder(conn, ctx)
	//clients.UpdateOrders(conn, ctx)
	//clients.SearchOrders(conn, ctx)
	//clients.ProcessOrders(conn, ctx)

}
