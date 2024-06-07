package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tubes.sister/raft/client/http/application"
	"tubes.sister/raft/client/http/cluster"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type HTTPClient struct {
	conn   *grpc.ClientConn
	router *chi.Mux
	port   int
}

func NewHTTPClient(clientPort int, serverAddr string) *HTTPClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(serverAddr, opts...)
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

	router.Route("/app", func(r chi.Router) {
		r.Get("/{key}", func(w http.ResponseWriter, r *http.Request) {
			application.Get(client, w, r)
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			application.Set(client, w, r)
		})
		r.Get("/{key}/strlen", func(w http.ResponseWriter, r *http.Request) {
			application.Strlen(client, w, r)
		})
		r.Delete("/{key}", func(w http.ResponseWriter, r *http.Request) {
			application.Delete(client, w, r)
		})
		r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
			application.Append(client, w, r)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			application.GetAll(client, w, r)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			application.DelAll(client, w, r)
		})
	})
	router.Route("/cluster", func(r chi.Router) {
		r.Get("/ping", cluster.Ping)
	})

	return &HTTPClient{conn: conn, router: router, port: clientPort}
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
