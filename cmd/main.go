package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joelgarciajr84/go-grpc-stream-server/internal/server"
	pb "github.com/joelgarciajr84/go-grpc-stream-server/pkg/pb"
	"google.golang.org/grpc"
)

func main() {
	port := envOr("PORT", "50005")
	numResponses := envIntOr("RESPONSE_COUNT", 5)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, &server.Server{NumResponses: numResponses})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("shutting down gracefully...")
		s.GracefulStop()
	}()

	log.Printf("server listening on :%s (RESPONSE_COUNT=%d)", port, numResponses)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
