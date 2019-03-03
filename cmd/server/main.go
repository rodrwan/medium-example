package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	"github.com/99designs/gqlgen/handler"

	"github.com/rodrwan/medium-example/database"
	"github.com/rodrwan/medium-example/graphql"
	"github.com/rodrwan/medium-example/service"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func main() {
	postgresDSN := flag.String("postgres-dsn", "postgres://mediumexample:me1234@localhost:5432/example?sslmode=disable", "Postgres domain service name")
	port := flag.Int("port", 8080, "server port")

	flag.Parse()

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	db, err := database.NewPostgres(logger, *postgresDSN)
	if err != nil {
		panic(err)
	}

	svc := service.NewService(db)

	logs := func(logger *log.Logger) func(writer io.Writer, params handlers.LogFormatterParams) {
		return func(writer io.Writer, params handlers.LogFormatterParams) {
			logger.WithFields(log.Fields{
				"server":      "GraphQL",
				"status_code": params.StatusCode,
				"size":        params.Size,
				"method":      params.Request.Method,
				"timestamp":   time.Now(),
				"url":         params.URL.String(),
			}).Info(http.StatusText(params.StatusCode))
		}
	}

	rootHandler := handler.GraphQL(
		graphql.NewExecutableSchema(
			graphql.Config{
				Resolvers: &graphql.Resolver{
					Service: svc,
				},
			},
		),
	)

	programLogger := logger.WithFields(log.Fields{
		"server": "GraphQL",
	})
	mux := http.NewServeMux()
	mux.Handle("/", handler.Playground("GraphQL playground", "/query"))
	mux.Handle("/query", handlers.CustomLoggingHandler(os.Stdout, rootHandler, logs(logger)))

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	programLogger.Infof("connect to http://localhost:%d/ for GraphQL playground", *port)
	programLogger.Fatal(server.ListenAndServe())
}
