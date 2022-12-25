package main

import (
	"context"
	"log"
	"time"

	pb "github.com/shinemost/grpc-up-client/pbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Apple iPhone 12"
	description := "Meet Apple iPhone 12. All-new dual-camera system with Ultra Wide and Night mode."
	price := float32(1699.00)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

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
