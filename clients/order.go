package clients

import (
	"context"
	"io"
	"log"

	pb "github.com/shinemost/grpc-up-client/pbs"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func AddOrder(conn *grpc.ClientConn, ctx context.Context) {
	c := pb.NewOrderManagementClient(conn)
	id := "103"
	description := "11"
	price := float32(199.00)
	destination := "hefeixinzhan"
	items := []string{"jjcai", "xfchen1", "qinkun"}

	r, err := c.AddOrder(ctx, &pb.Order{
		Id:          id,
		Items:       items,
		Description: description,
		Price:       price,
		Destination: destination,
	})
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order ID: %s added successfully", r.GetValue())

}

func SearchOrders(conn *grpc.ClientConn, ctx context.Context) {
	c := pb.NewOrderManagementClient(conn)
	searchStream, _ := c.SearchOrders(ctx, &wrapperspb.StringValue{Value: "iPad Pro"})
	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			break
		}
		log.Print("search Result:", searchOrder)
	}
}

func UpdateOrders(conn *grpc.ClientConn, ctx context.Context) {
	c := pb.NewOrderManagementClient(conn)
	updateOrderStream, err := c.UpdateOrders(ctx)
	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _,%v", c, err)
	}

	updOrder1 := pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Google Pixel Book"}, Destination: "Mountain View, CA", Price: 1100.00}
	updOrder2 := pb.Order{Id: "103", Items: []string{"Apple Watch S4", "Mac Book Pro", "iPad Pro"}, Destination: "San Jose, CA", Price: 2800.00}
	updOrder3 := pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub", "iPad Mini"}, Destination: "Mountain View, CA", Price: 2200.00}

	if err := updateOrderStream.Send(&updOrder1); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateOrderStream, updOrder1, err)
	}

	if err := updateOrderStream.Send(&updOrder2); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateOrderStream, updOrder2, err)
	}

	if err := updateOrderStream.Send(&updOrder3); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateOrderStream, updOrder3, err)
	}

	updateRes, err := updateOrderStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateOrderStream, err, nil)
	}
	log.Printf("Update Orders Res : %s", updateRes)
}
