package clients

import (
	"context"
	"log"

	pb "github.com/shinemost/grpc-up-client/pbs"
	"google.golang.org/grpc"
)

func P(conn *grpc.ClientConn, ctx context.Context) {
	c := pb.NewProductInfoClient(conn)
	name := "Apple iPhone 12"
	description := "Meet Apple iPhone 12. All-new dual-camera system with Ultra Wide and Night mode."
	price := float32(1699.00)

	r, err := c.AddProduct(ctx, &pb.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	p, err := c.GetProduct(ctx, &pb.ProductID{
		Value: r.Value,
	})

	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}

	log.Printf("Product: %v", p.String())
}
