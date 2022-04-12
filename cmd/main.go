package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/joelgarciajr84/go-grpc-stream-server/pkg/pb"
	"google.golang.org/grpc"
)

type server struct{}

func (s server) FetchResponse(in *pb.Request, srv pb.StreamService_FetchResponseServer) error {

	log.Printf("RECEIVED REQUEST ID : %d", in.Id)

	var wg sync.WaitGroup
	for i := 0; i < 42; i++ {
		wg.Add(1)
		go func(count int64) {
			defer wg.Done()

			time.Sleep(time.Duration(count) * time.Second)
			resp := pb.Response{Result: fmt.Sprintf("RESPONSE #%d FOR REQUEST ID:%d", count, in.Id)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("Sending response to request number : %d", count)
		}(int64(i))
	}

	wg.Wait()
	return nil
}

func main() {

	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, server{})

	log.Println("Starting Server")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed: %v", err)
	}

}
