package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"

	pb "github.com/dshearer/modelhawk/gen/go/v0"
	"github.com/dshearer/modelhawk/reference-impls/server/ollama"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedNotifyServiceServer
	pb.UnimplementedPermissionServiceServer 
	ollama *ollama.Client
}

func makePointer[T any](v T) *T {
	return &v
}

func (s *server) WillCallTool(ctx context.Context, req *pb.WillCallToolRequest) (*pb.ServiceStatusResponse, error) {
	// log.Printf("WillCallTool: app=%q tool=%q args=%v", req.GetApp().GetValue(), req.GetToolName(), req.GetArgs())
	return &pb.ServiceStatusResponse{Result: makePointer(pb.ServiceStatusResponse_RESULT_OK), Msg: makePointer("ok")}, nil
}

func (s *server) DidCallTool(ctx context.Context, req *pb.DidCallToolRequest) (*pb.ServiceStatusResponse, error) {
	// log.Printf("DidCallTool: app=%q tool=%q args=%v result=%q", req.GetApp().GetValue(), req.GetToolName(), req.GetArgs(), req.GetResult())
	return &pb.ServiceStatusResponse{Result: makePointer(pb.ServiceStatusResponse_RESULT_OK), Msg: makePointer("ok")}, nil
}

func (s *server) WantsToCallTool(ctx context.Context, req *pb.WantsToCallToolRequest) (*pb.WantsToCallToolResponse, error) {
	// log.Printf("WantsToCallTool: app=%q tool=%q args=%v msgs=%v", req.GetApp().GetValue(), req.GetToolName(), req.GetArgs(), req.GetLastMessages())

	// ask Ollama
	result, err := s.ollama.WantsToCallTool(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("Ollama says: %v", result)
	return &pb.WantsToCallToolResponse{
		Permitted: makePointer(!result.Malicious),
		Details:   makePointer(result.Why),
	}, nil
}

func main() {
	port := flag.Int("port", 50051, "port to listen on")
	ollamaURL := flag.String("ollama-url", "http://localhost:11434", "Ollama server URL")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize Ollama client
	parsedOllamaURL, err := url.Parse(*ollamaURL)
	if err != nil {
		log.Fatalf("failed to pasrse Ollama URL: %v", err)
	}
	ollamaClient, err := ollama.NewClient(parsedOllamaURL, "server-log")
	if err != nil {
		log.Fatalf("failed to start logging: %v", err)
	}

	svc := &server{
		ollama: ollamaClient,
	}

	// Check Ollama connection
	if err := svc.ollama.CheckConnection(context.Background()); err != nil {
		log.Printf("Warning: Could not connect to Ollama: %v", err)
		log.Printf("Continuing without Ollama integration...")
	}

	s := grpc.NewServer()
	pb.RegisterNotifyServiceServer(s, svc)
	pb.RegisterPermissionServiceServer(s, svc)

	log.Printf("modelhawk server listening on :%d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
