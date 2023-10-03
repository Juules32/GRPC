package main

import (
	"io"
	"log"
	"net"
	"sync"

	gRPC "github.com/Juules32/GRPC/proto"

	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedStreamingServiceServer        // You need this line if you have a server
	name                                     string // Not required but useful if you want to name your server
	port                                     string // Not required but useful if your server needs to know what port it's listening to

	mutex sync.Mutex // used to lock the server to avoid race conditions.
}

func (s *Server) StreamData(msgStream gRPC.StreamingService_StreamDataServer) error {
	for {
		// get the next message from the stream
		msg, err := msgStream.Recv()

		// the stream is closed so we can exit the loop
		if err == io.EOF {
			return nil
		}
		// some other error
		if err != nil {
			log.Fatalf("Error receiving data: %v", err)
			return err
		}
		// log the message
		log.Printf("Received message from %s: %s", msg.ClientName, msg.Message)
		err = msgStream.Send(&gRPC.DataResponse{Message: msg.Message})
		if err != nil {
			log.Fatalf("Error sending response: %v", err)
			return err
		}
	}
}

func main() {
	list, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		return
	}

	grpcServer := grpc.NewServer()

	server := &Server{
		name: "localhost",
		port: "5400",
	}

	gRPC.RegisterStreamingServiceServer(grpcServer, server)

	log.Printf("Server %s: Listening at %v\n", "localhost", list.Addr())

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
