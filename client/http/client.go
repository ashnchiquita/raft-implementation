package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"tubes.sister/raft/client/http/application"
	"tubes.sister/raft/client/http/cluster"
)

func StartHTTPClient(clientPort, serverPort int) {
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

	log.Printf("Web client started at port %d", clientPort)
	http.ListenAndServe(fmt.Sprintf(":%d", clientPort), router)
}
