package grpc

import (
	"fmt"
	"log"
	"net"
	"sync"

	"uooobarry/yuuna-danmu/pkg/server/grpc/pb"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedLiveServiceServer
	mu          sync.RWMutex
	subscribers map[chan *pb.LiveEvent]struct{}
	grpcServer  *grpc.Server
}

func New() *GRPCServer {
	return &GRPCServer{
		subscribers: make(map[chan *pb.LiveEvent]struct{}),
	}
}

func (s *GRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterLiveServiceServer(s.grpcServer, s)

	go func() {
		log.Printf("gRPC server listening on %d", port)
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	return nil
}

func (s *GRPCServer) Subscribe(_ *pb.Empty, stream pb.LiveService_SubscribeServer) error {
	ch := make(chan *pb.LiveEvent, 10)

	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.subscribers, ch)
		s.mu.Unlock()
		close(ch)
	}()

	for {
		select {
		case <-stream.Context().Done():
			if err := stream.Context().Err(); err != nil {
				return err
			}
			return nil
		case event := <-ch:
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}

func (s *GRPCServer) Dispatch(event any) {
	log.Printf("Dispatching clean event: %v", event)
	pbEvent := s.mapToProto(event)
	if pbEvent == nil {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for ch := range s.subscribers {
		select {
		case ch <- pbEvent:
		default:
		}
	}
}

func (s *GRPCServer) Stop() error {
	s.grpcServer.GracefulStop()
	return nil
}
