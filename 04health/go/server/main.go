package main

import (
	"context"
	"log"
	"sync"

	pb "github.com/clin211/grpc/health/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedHealthServer
	mu        sync.RWMutex
	shutdown  bool
	statusMap map[string]pb.HealthCheckResponse_ServingStatus
	updates   map[string]map[pb.Health_WatchServer]chan pb.HealthCheckResponse_ServingStatus
}

func NewServer() *Server {
	return &Server{
		statusMap: map[string]pb.HealthCheckResponse_ServingStatus{
			"": pb.HealthCheckResponse_NOT_SERVING,
		},
		updates: make(map[string]map[pb.Health_WatchServer]chan pb.HealthCheckResponse_ServingStatus),
	}
}

func (s *Server) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if servingStatus, ok := s.statusMap[in.Service]; ok {
		return &pb.HealthCheckResponse{
			Status: servingStatus,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "service not found")
}

func (s *Server) Watch(in *pb.HealthCheckRequest, stream pb.Health_WatchServer) error {
	service := in.Service

	update := make(chan pb.HealthCheckResponse_ServingStatus, 1)

	s.mu.Lock()

	if servingStatus, ok := s.statusMap[service]; ok {
		update <- servingStatus
	} else {
		update <- pb.HealthCheckResponse_SERVICE_UNKNOWN
	}

	if _, ok := s.updates[service]; !ok {
		s.updates[service] = make(map[pb.Health_WatchServer]chan pb.HealthCheckResponse_ServingStatus)
	}

	s.updates[service][stream] = update
	defer func() {
		s.mu.Lock()
		delete(s.updates[service], stream)
		s.mu.Unlock()
	}()
	s.mu.Unlock()

	var lastSentStatus pb.HealthCheckResponse_ServingStatus = -1
	for {
		select {
		case servingStatus := <-update:
			if lastSentStatus == servingStatus {
				continue
			}

			lastSentStatus = servingStatus
			err := stream.Send(&pb.HealthCheckResponse{
				Status: servingStatus,
			})
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		}
	}
}

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *Server) SetServingStatus(service string, servingStatus pb.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.shutdown {
		log.Printf("health: status changing for %s to %v is ignored because health service is shutdown", service, servingStatus)
		return
	}

	s.setServingStatusLocked(service, servingStatus)
}

func (s *Server) setServingStatusLocked(service string, servingStatus pb.HealthCheckResponse_ServingStatus) {
	s.statusMap[service] = servingStatus
	for _, update := range s.updates[service] {
		// Clears previous updates, that are not sent to the client, from the channel.
		// This can happen if the client is not reading and the server gets flow control limited.
		select {
		case <-update:
		default:
		}
		// Puts the most recent update to the channel.
		update <- servingStatus
	}
}

// Shutdown sets all serving status to NOT_SERVING, and configures the server to
// ignore all future status changes.
func (s *Server) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shutdown = true
	for service := range s.statusMap {
		s.setServingStatusLocked(service, pb.HealthCheckResponse_NOT_SERVING)
	}
}

// Resume sets all serving status to SERVING, and configures the server to
// accept all future status changes.
func (s *Server) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shutdown = false
	for service := range s.statusMap {
		s.setServingStatusLocked(service, pb.HealthCheckResponse_SERVING)
	}
}

func main() {

}
