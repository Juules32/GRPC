package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	gRPC "github.com/Juules32/GRPC/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var server gRPC.TemplateClient  //the server
var ServerConn *grpc.ClientConn //the server connection

func conReady(s gRPC.TemplateClient) bool {
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

	server = gRPC.NewTemplateClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())

	reader := bufio.NewReader(os.Stdin)

	//Infinite loop to listen for clients input.
	for {
		fmt.Print("-> ")

		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input

		if !conReady(server) {
			log.Printf("Client %s: something was wrong with the connection to the server :5400")
			continue
		}

		message := &gRPC.Greeting{
			ClientName: "benja",
			Message:    "hejsa",
		}

		stream, err := server.SayHi(context.Background())

		stream.Send(message)
		farewell, err := stream.CloseAndRecv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("server says: ", farewell)
	}
}
