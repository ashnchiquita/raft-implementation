package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	app_http "tubes.sister/raft/client/http/application"
	raft_http "tubes.sister/raft/client/http/raft"
	"tubes.sister/raft/hello"
	"tubes.sister/raft/server"
)

var (
	serverr    = flag.Bool("is_server", true, "GRPC Server or Client")
	port       = flag.Int("port", 50051, "GRPC Server Port")
	serverAddr = flag.String("server_address", "localhost:50051", "GRPC Server Address")
	clientPort = flag.Int("client_port", 3000, "HTTP Client Port")
)

type ClientRequest struct {
	Message string `json:"message"`
}

type ClientRes struct {
	Result string `json:"result"`
}

func main() {
	flag.Parse()

	if *serverr {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)

		helloServer := hello.NewHelloServerImpl(server.Address{IP: "localhost", Port: fmt.Sprintf("%d", *port)})

		hello.RegisterHelloServer(grpcServer, helloServer)

		log.Printf("Server started at port %d", *port)
		grpcServer.Serve(lis)
	} else {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		conn, err := grpc.Dial(*serverAddr, opts...)
		if err != nil {
			log.Fatalf("Failed to dial server: %v", err)
		}

		defer conn.Close()

		client := hello.NewHelloClient(conn)

		router := chi.NewRouter()
		router.Use(middleware.Logger)
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"POST"},
		}))

		router.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req ClientRequest

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			response, err := client.SayHello(context.Background(), &hello.HelloMsg{Name: req.Message})
			if err != nil {
				log.Printf("Failed to call SayHello: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			res := &ClientRes{Result: response.Message}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		router.Post("/addNode", func(w http.ResponseWriter, r *http.Request) {
			var req hello.AddNodeRequest

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Kirim request AddNode ke server
			_, err = client.AddNode(context.Background(), &req)
			if err != nil {
				log.Printf("Failed to call AddNode: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		})

		// Application Router
		router.Route("/app", func(r chi.Router) {
			r.Get("/ping", app_http.Ping)
			r.Get("/{key}", app_http.Get)
			r.Put("/", app_http.Set)
			r.Get("/{key}/strlen", app_http.Strlen)
			r.Delete("/{key}", app_http.Del)
			r.Patch("/", app_http.Append)

			r.Get("/", app_http.GetAll)
			r.Delete("/", app_http.DelAll)
		})

		// Raft Router
		router.Route("/raft", func(r chi.Router) {
			r.Get("/get-statuses", raft_http.GetStatuses)
		})

		log.Printf("Client started at port %d", *clientPort)
		http.ListenAndServe(fmt.Sprintf(":%d", *clientPort), router)
	}
}
