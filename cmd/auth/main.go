package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	authpb "micro-blog/proto/auth"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
}

var kafkaWriter *kafka.Writer

func main() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "user.created",
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	go startGRPCServer()
	startHTTPServer()
}

func startGRPCServer() {
	list, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("grpc server start fail: %v", err)
	}

	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, &AuthServer{})

	log.Println("auth grpc server started on :8081")

	if err := server.Serve(list); err != nil {
		log.Fatal(err)
	}
}

func startHTTPServer() {
	router := http.NewServeMux()
	router.HandleFunc("/register", registerHandler)

	log.Println("auth http server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("http server start fail: %v", err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}

	event := `{"user_id": "1"}`

	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	err := kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte("user-1"),
		Value: []byte(event),
	})
	if err != nil {
		log.Printf("kafka message send fail: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to register user"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user registered"))
}

func (s *AuthServer) ValidateToken(
	ctx context.Context,
	req *authpb.ValidateTokenRequest,
) (*authpb.ValidateTokenResponse, error) {
	if req.Token == "valid-token" {
		return &authpb.ValidateTokenResponse{
			Valid: true,
		}, nil
	}

	return &authpb.ValidateTokenResponse{
		Valid: false,
	}, nil
}
