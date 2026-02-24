package main

import (
	"log"
	"net/http"
	"time"

	"connectrpc.com/grpcreflect"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	chatv1connect "github.com/zkxjzmswkwl/go-grpc-testing/gen/go/proto/chat/v1/chatv1connect"
	"github.com/zkxjzmswkwl/go-grpc-testing/internal/app"
	"github.com/zkxjzmswkwl/go-grpc-testing/internal/store/memory"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	store := memory.NewStore()
	svc := app.NewChatService(store)

	// (csmith): can still use chi middleware etc
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	path, handler := chatv1connect.NewChatServiceHandler(svc)
	r.Mount(path, handler)

	reflector := grpcreflect.NewStaticReflector("chat.v1.ChatService")
	r.Handle(grpcreflect.NewHandlerV1(reflector))
	r.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, h2c.NewHandler(r, &http2.Server{})))
}
