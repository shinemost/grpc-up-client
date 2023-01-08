package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/shinemost/grpc-up-client/tracer"
	"go.opencensus.io/plugin/ocgrpc"

	"github.com/shinemost/grpc-up-client/clients"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	//address = "grpc-server-service.local:50051"
	address  = "localhost:50051"
	hostname = "localhost"
	caFile   = "certs/ca.crt"
	crtFile  = "certs/clinet.pem"
	keyFile  = "certs/clinet.key"
)

func main() {
	// initTracing()
	// var wg sync.WaitGroup
	// wg.Add(1)
	// req := prometheus.NewRegistry()
	// grpcMetrics := grpc_prometheus.NewClientMetrics()
	// req.MustRegister(grpcMetrics)
	//基于公钥证书创建TLS证书
	// creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	// if err != nil {
	// 	log.Fatalf("failde to load credentials:%v", err)
	// }

	// Register stats and trace exporters to export
	// the collected data.
	// view.RegisterExporter(&exporter.PrintExporter{})

	// Register the view to collect gRPC client stats.
	// if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
	// 	log.Fatal(err)
	// }

	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load client key pair : %s", err)
	}

	//auth := models.BasicAuth{
	//	Username: "admin",
	//	Password: "admin",
	//}
	// auth := oauth.NewOauthAccess(models.FetchToken())

	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("cloud not find ca cert: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to add ca cerWts")
	}

	jaegertracer, closer, err := tracer.NewTracer("grpc-client")
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	conn, err := grpc.Dial(
		address,
		// grpc.WithPerRPCCredentials(auth),
		// grpc.WithTransportCredentials(creds),
		// grpc.WithUnaryInterceptor(interceptor.OrderUnaryClientInterceptor),
		// grpc.WithStreamInterceptor(interceptor.ClientStreamInterceptor),
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				ServerName:   hostname,
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
			}),
		),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		// grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithUnaryInterceptor(grpcopentracing.UnaryClientInterceptor(grpcopentracing.WithTracer(jaegertracer))),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// httpServer := &http.Server{
	// 	Handler: promhttp.HandlerFor(req, promhttp.HandlerOpts{}),
	// 	Addr:    "localhost:9094",
	// }

	// go func() {
	// 	defer wg.Done()
	// 	if err := httpServer.ListenAndServe(); err != nil {
	// 		log.Fatal("Unalbe to start a http server")
	// 	}
	// }()

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
	for i := 0; i < 10; i++ {
		clients.P(conn, ctx)
		time.Sleep(3 * time.Second)
	}

	//clients.UpdateOrders(conn, ctx, cancel)
	//clients.SearchOrders(conn, ctx)
	//clients.ProcessOrders(conn, ctx)
	// wg.Wait()
}

// func main() {
// 	clients.StartClient()
// }
