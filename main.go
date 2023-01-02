package main

import (
	"context"
	"log"

	"github.com/shinemost/grpc-up-client/clients"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = "certs/localhost.crt"
)

func main() {
	//基于公钥证书创建TLS证书
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("failde to load credentials:%v", err)
	}

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(creds),
		// grpc.WithUnaryInterceptor(interceptor.OrderUnaryClientInterceptor),
		// grpc.WithStreamInterceptor(interceptor.ClientStreamInterceptor),
	)
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
	//cancel()

	//RPC客户端的多路复用，多个客户端共用一个连接
	clients.P(conn, ctx)
	//clients.UpdateOrders(conn, ctx, cancel)
	//clients.SearchOrders(conn, ctx)
	//clients.ProcessOrders(conn, ctx)

}

// func main() {
// 	clients.StartClient()
// }
