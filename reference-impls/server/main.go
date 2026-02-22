package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/dshearer/modelhawk/gen/go/v1"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedInfoServiceServer
	pb.UnimplementedNotifyServiceServer
	pb.UnimplementedPermissionServiceServer
}

func makePointer[T any](v T) *T {
	return &v
}

func (s *server) GiveToolInfo(ctx context.Context, req *pb.GiveToolInfoRequest) (*pb.ServiceStatusResponse, error) {
	log.Printf("GiveToolInfo: app=%q tool=%q desc=%q args=%v", req.GetApp().GetValue(), req.GetName(), req.GetDesc(), req.GetArgs())
	return &pb.ServiceStatusResponse{Result: makePointer(pb.ServiceStatusResponse_RESULT_OK), Msg: makePointer("ok")}, nil
}

func (s *server) WillCallTool(ctx context.Context, req *pb.WillCallToolRequest) (*pb.ServiceStatusResponse, error) {
	log.Printf("WillCallTool: app=%q tool=%q args=%v", req.GetApp().GetValue(), req.GetToolName(), req.GetArgs())
	return &pb.ServiceStatusResponse{Result: makePointer(pb.ServiceStatusResponse_RESULT_OK), Msg: makePointer("ok")}, nil
}

func (s *server) DidCallTool(ctx context.Context, req *pb.DidCallToolRequest) (*pb.ServiceStatusResponse, error) {
	log.Printf("DidCallTool: app=%q tool=%q args=%v result=%q", req.GetApp().GetValue(), req.GetToolName(), req.GetArgs(), req.GetResult())
	return &pb.ServiceStatusResponse{Result: makePointer(pb.ServiceStatusResponse_RESULT_OK), Msg: makePointer("ok")}, nil
}

func (s *server) WantsToCallTool(ctx context.Context, req *pb.WantsToCallToolRequest) (*pb.WantsToCallToolResponse, error) {
	log.Printf("WantsToCallTool: app=%q tool=%q args=%v msgs=%v", req.GetApp().GetValue(), req.GetToolName(), req.GetArgs(), req.GetLastMessages())
	return &pb.WantsToCallToolResponse{
		Permitted: makePointer(false),
		Details:   makePointer("denied by default"),
	}, nil
}

func main() {
	port := flag.Int("port", 50051, "port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	svc := &server{}
	pb.RegisterInfoServiceServer(s, svc)
	pb.RegisterNotifyServiceServer(s, svc)
	pb.RegisterPermissionServiceServer(s, svc)

	log.Printf("modelhawk server listening on :%d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
