package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/joelgarciajr84/go-grpc-stream-server/pkg/pb"
)

type Server struct {
	pb.UnimplementedStreamServiceServer
	NumResponses int
}

func (s *Server) FetchResponse(in *pb.Request, srv pb.StreamService_FetchResponseServer) error {
	log.Printf("received request id=%d, sending %d responses", in.Id, s.NumResponses)

	ctx := srv.Context()
	results := make(chan string, s.NumResponses)

	var wg sync.WaitGroup
	for i := 0; i < s.NumResponses; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(count) * time.Second):
			}
			results <- fmt.Sprintf("RESPONSE #%d FOR REQUEST ID:%d", count, in.Id)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if err := srv.Send(&pb.Response{Result: result}); err != nil {
			return err
		}
		log.Printf("sent: %s", result)
	}

	return ctx.Err()
}
