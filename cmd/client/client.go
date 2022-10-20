package main

import (
	"context"
	"fmt"
	"github.com/henriquelucasdf/fc3-grpc/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	// defer: closes the connection when it's not used
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// Unary request
	//AddUser(client)

	//Server stream response
	//AddUserVerbose(client)

	// Client stream request
	//AddUsers(client)

	// Bidirectional stream
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive message: %v", err)
		}

		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "h1",
			Name:  "Henrique1",
			Email: "h1@h.com",
		},
		&pb.User{
			Id:    "h2",
			Name:  "Henrique2",
			Email: "h2@h.com",
		},
		&pb.User{
			Id:    "h3",
			Name:  "Henrique3",
			Email: "h3@h.com",
		},
		&pb.User{
			Id:    "h4",
			Name:  "Henrique4",
			Email: "h4@h.com",
		},
		&pb.User{
			Id:    "h5",
			Name:  "Henrique5",
			Email: "h5@h.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	reqs := []*pb.User{
		&pb.User{
			Id:    "h1",
			Name:  "Henrique1",
			Email: "h1@h.com",
		},
		&pb.User{
			Id:    "h2",
			Name:  "Henrique2",
			Email: "h2@h.com",
		},
		&pb.User{
			Id:    "h3",
			Name:  "Henrique3",
			Email: "h3@h.com",
		},
		&pb.User{
			Id:    "h4",
			Name:  "Henrique4",
			Email: "h4@h.com",
		},
		&pb.User{
			Id:    "h5",
			Name:  "Henrique5",
			Email: "h5@h.com",
		},
	}

	// go routings: async threads - one to send and one to receive
	// the bellow anonymous func will run in background

	// channel: to hold the app until all the responses are received
	wait := make(chan int)
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 3)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
