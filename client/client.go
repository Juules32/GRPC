package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"time"

	gRPC "github.com/Juules32/GRPC/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var server gRPC.StreamingServiceClient //the server
var ServerConn *grpc.ClientConn        //the server connection

func conReady(s gRPC.StreamingServiceClient) bool {
	return ServerConn.GetState().String() == "READY"
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(":5400", opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}

	server = gRPC.NewStreamingServiceClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())
	defer ServerConn.Close()

	stream, err := server.StreamData(context.Background())

	go sendMessages(stream)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // End of the stream
		}
		if err != nil {
			log.Fatalf("Error receiving response: %v", err)
		}
		log.Printf("Received response: %s", resp.Message)
	}
}

func sendMessages(stream gRPC.StreamingService_StreamDataClient) {
	reader := bufio.NewReader(os.Stdin)

	// Send data to the server
	for {

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		greeting := &gRPC.DataRequest{ClientName: "benja", Message: input}
		if err := stream.Send(greeting); err != nil {
			log.Fatalf("Error sending data: %v", err)
		}
		time.Sleep(time.Second)
	}
}
