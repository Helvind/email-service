package main

import (
	"context"
	"log"

	pb "github.com/Helvind/email-service/proto"
	"google.golang.org/grpc"
)

// const host = "localhost:50051"
const host = "35.235.39.42:50051"

func main() {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed dialing host, %q", err.Error())
	}

	client := pb.NewEmailSenderClient(conn)

	// TODO(helvind): implemnt flags for cli usage
	email := &pb.SendEmailRequest{
		Email: &pb.Email{
			To:      "jakob.helvind@gmail.com",
			From:    "someone@gmail.com",
			Subject: "Hello World",
			Body:    "Testing SparkPost",
		},
	}
	resp, err := client.SendEmail(context.Background(), email)
	if err != nil {
		log.Fatalf("failed sending e-mail, %q", err.Error())
	}

	log.Printf("send e-mail, status %d", resp.GetStatus())
}
