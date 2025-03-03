package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	chatpb "server/proto"

	"google.golang.org/grpc"
)

type ChatServer struct {
	chatpb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[chatpb.ChatService_ChatServer]bool
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[chatpb.ChatService_ChatServer]bool),
	}
}

func (s *ChatServer) Chat(stream chatpb.ChatService_ChatServer) error {
	s.addClient(stream)
	defer s.removeClient(stream)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return err
		}
		log.Printf("Received from %s: %s", msg.User, msg.Message)
		s.broadcast(stream, msg)
	}
}

func (s *ChatServer) addClient(stream chatpb.ChatService_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[stream] = true
}

func (s *ChatServer) removeClient(stream chatpb.ChatService_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, stream)
}

func (s *ChatServer) broadcast(sender chatpb.ChatService_ChatServer, msg *chatpb.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for client := range s.clients {
		if client == sender {
			continue
		}
		if err := client.Send(msg); err != nil {
			log.Printf("Error sending message to a client: %v", err)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatpb.RegisterChatServiceServer(grpcServer, NewChatServer())

	fmt.Printf("Server listening at %v\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
