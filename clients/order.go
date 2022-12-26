package clients

import (
	"context"
	"io"
	"log"

	uuid "github.com/satori/go.uuid"
	pb "github.com/shinemost/grpc-up-client/pbs"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func AddOrder(conn *grpc.ClientConn, ctx context.Context) {
	c := pb.NewOrderManagementClient(conn)
	id := uuid.NewV4()
	description := "xxxxxxx"
	price := float32(699.00)
	destination := "hefei"

	r, err := c.AddOrder(ctx, &pb.Order{
		Id:          id.String(),
		Items:       []string{"jjcai5", "xfchen12", "qinkun12"},
		Description: description,
		Price:       price,
		Destination: destination,
	})
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order ID: %s added successfully", r.GetValue())

	searchStream, _ := c.SearchOrders(ctx, &wrapperspb.StringValue{Value: "xfchen12"})
	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			break
		}
		log.Print("search Result:", searchOrder)
	}
}
