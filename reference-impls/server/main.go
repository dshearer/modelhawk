package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/dshearer/modelhawk/reference-impls/server/modelhawk/v1"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedInfoServiceServer
	pb.UnimplementedNotifyServiceServer
	pb.UnimplementedPermissionServiceServer
}

func (s *server) GiveToolInfo(ctx context.Context, req *pb.GiveToolInfoRequest) (*pb.ServiceStatusResponse, error) {
	log.Printf("GiveToolInfo: app=%q tool=%q desc=%q args=%v", req.GetApp().GetValue(), req.GetName(), req.GetDesc(), req.GetArgs())
	return &pb.ServiceStatusResponse{Result: 0, Msg: "ok"}, nil
}

func (s *server) CalledTool(ctx context.Context, req *pb.CalledToolRequest) (*pb.ServiceStatusResponse, error) {
	log.Printf("CalledTool: app=%q tool=%q args=%v", req.GetApp().GetValue(), req.GetName(), req.GetArgs())
	return &pb.ServiceStatusResponse{Result: 0, Msg: "ok"}, nil
}

func (s *server) GotToolResponse(ctx context.Context, req *pb.GotToolResponseRequest) (*pb.ServiceStatusResponse, error) {
	log.Printf("GotToolResponse: app=%q tool=%q resp=%q", req.GetApp().GetValue(), req.GetName(), req.GetResp())
	return &pb.ServiceStatusResponse{Result: 0, Msg: "ok"}, nil
}

func (s *server) WantsToCallTool(ctx context.Context, req *pb.WantsToCallToolRequest) (*pb.WantsToCallToolResponse, error) {
	log.Printf("WantsToCallTool: app=%q tool=%q args=%v", req.GetApp().GetValue(), req.GetName(), req.GetArgs())
	return &pb.WantsToCallToolResponse{
		Result:  pb.WantsToCallToolResponse_RESULT_PERMITTED,
		Details: "permitted by default",
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
