package main

import (
	"context"
	"log"
	"net/http"
	"time"

	authpb "micro-blog/internal/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var authClient authpb.AuthServiceClient

func main() {
	conn, err := grpc.NewClient(
		"localhost:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	authClient = authpb.NewAuthServiceClient(conn)

	router := http.NewServeMux()
	router.HandleFunc("/checkAuth", postsHandler)

	srv := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	log.Printf("auth service started on :8082")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if !checkAuth(token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("posts list"))
}

func checkAuth(token string) bool {
	if authClient == nil {
		log.Println("authClient is nil")
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		log.Printf("rpc token validation error: %v", err)
		return false
	}

	if resp == nil {
		log.Println("auth response is nil")
		return false
	}

	return resp.Valid
}
