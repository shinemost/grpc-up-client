package clients

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/shinemost/grpc-up-client/pbs"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func AddOrder(conn *grpc.ClientConn, ctx context.Context) {
	c := pb.NewOrderManagementClient(conn)
	id := "100"
	description := "测试错误的ID造成的异常"
	price := float32(199.00)
	destination := "hefeixinzhan"
	items := []string{"jjcai", "xfchen1", "qinkun"}

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"kn", "vn",
	)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)
	ctxA := metadata.AppendToOutgoingContext(mdCtx, "k1", "v1", "k1", "v2", "k2", "v3")

	var header, trailer metadata.MD

	//睡眠5秒，过deadline
	//time.Sleep(time.Second * 5)
	r, err := c.AddOrder(ctxA, &pb.Order{
		Id:          id,
		Items:       items,
		Description: description,
		Price:       price,
		Destination: destination,
	}, grpc.Header(&header), grpc.Trailer(&trailer), grpc.UseCompressor(gzip.Name))

	// Reading the headers
	if t, ok := header["timestamp"]; ok {
		log.Printf("timestamp from header:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Fatal("timestamp expected but doesn't exist in header")
	}
	if l, ok := header["location"]; ok {
		log.Printf("location from header:\n")
		for i, e := range l {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Fatal("location expected but doesn't exist in header")
	}

	//放在这里是不行的，因为已经发起请求结束了，只是没有处理后续逻辑而已
	//time.Sleep(time.Second * 5)
	//if err != nil {
	//	got := status.Code(err)
	//	log.Fatalf("Could not add order: %v", got)
	//}
	//log.Printf("Order ID: %s added successfully", r.GetValue())

	if err != nil {
		errorCode := status.Code(err)
		if errorCode == codes.InvalidArgument {
			log.Printf("Invalid Argument Error : %s", errorCode)
			errorStatus := status.Convert(err)
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				case *epb.BadRequest_FieldViolation:
					log.Printf("Request Field Invalid: %s", info)
				default:
					log.Printf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Printf("Unhandled error : %s ", errorCode)
		}
	} else {
		log.Print("AddOrder Response -> ", r.Value)
	}
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

func UpdateOrders(conn *grpc.ClientConn, ctx context.Context, cancel context.CancelFunc) {
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

	cancel()

	updateRes, err := updateOrderStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateOrderStream, err, nil)
	}
	log.Printf("Update Orders Res : %s", updateRes)
}

func ProcessOrders(conn *grpc.ClientConn, ctx context.Context) {
	c := pb.NewOrderManagementClient(conn)
	streamProcOrder, err := c.ProcessOrders(ctx)
	if err != nil {
		log.Fatalf("%v.ProcessOrders(_) = _, %v", c, err)
	}

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "102"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "102", err)
	}

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "103"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "103", err)
	}

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "104"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "104", err)
	}
	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "105"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "105", err)
	}
	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "106"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "106", err)
	}

	channel := make(chan struct{})
	go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "101"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "101", err)
	}
	if err := streamProcOrder.CloseSend(); err != nil {
		log.Fatal(err)
	}
	channel <- struct{}{}
}
func asncClientBidirectionalRPC(streamProcOrder pb.OrderManagement_ProcessOrdersClient, c chan struct{}) {
	for {
		combinedShipment, errProcOrder := streamProcOrder.Recv()
		if errProcOrder == io.EOF {
			break
		}
		log.Printf("Combined shipment : ", combinedShipment.OrdersList)
	}
	<-c
}
