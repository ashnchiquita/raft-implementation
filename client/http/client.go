package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tubes.sister/raft/client/http/handler"
	_ "tubes.sister/raft/docs"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type HTTPClient struct {
	conn   *grpc.ClientConn
	router *chi.Mux
	port   int
	client handler.GRPCClient
}

// @title Web Client API
// @version 1.0
// @description This is the API documentation for the web client of the application.
func NewHTTPClient(clientPort int, serverAddr string) *HTTPClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "DELETE", "PUT", "OPTIONS", "PATCH"},
	}))

	client := gRPC.NewCmdExecutorClient(conn)
	hc := &HTTPClient{conn: conn, router: router, port: clientPort, client: *handler.NewGRPCClient(&client)}

	hc.router.Get("/swagger/*", httpSwagger.WrapHandler)

	hc.router.Route("/app", func(r chi.Router) {
		r.Get("/{key}", hc.client.Get)
		r.Put("/", hc.client.Set)
		r.Get("/{key}/strlen", hc.client.Strlen)
		r.Delete("/{key}", hc.client.Delete)
		r.Patch("/", hc.client.Append)
		r.Get("/", hc.client.GetAll)
		r.Delete("/", hc.client.DelAll)
	})
	hc.router.Route("/cluster", func(r chi.Router) {
		r.Get("/ping", hc.client.Ping)
		r.Get("/log", hc.client.RequestLog)
	})

	return hc
}

func (hc *HTTPClient) Start() {
	log.Printf("Web client started at port %d", hc.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", hc.port), hc.router))
}

func (hc *HTTPClient) Stop() {
	if err := hc.conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}
