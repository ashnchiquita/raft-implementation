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
)

type HTTPClient struct {
	conn   *grpc.ClientConn
	router *chi.Mux
	port   int
}

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

	// Application Router
	router.Route("/app", func(r chi.Router) {
		r.Get("/{key}", application.Get)
		r.Put("/", application.Set)
		r.Get("/{key}/strlen", application.Strlen)
		r.Delete("/{key}", application.Del)
		r.Patch("/", application.Append)

		r.Get("/", application.GetAll)
		r.Delete("/", application.DelAll)
	})

	// Cluster Router
	router.Route("/cluster", func(r chi.Router) {
		r.Get("/ping", cluster.Ping)
	})

	return &HTTPClient{conn: conn, router: router, port: clientPort}
}

func (hc *HTTPClient) Start() {
	log.Printf("Web client started at port %d", hc.port)
	http.ListenAndServe(fmt.Sprintf(":%d", hc.port), hc.router)
}

func (hc *HTTPClient) Stop() {
	hc.conn.Close()
}
