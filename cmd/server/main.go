package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	"github.com/99designs/gqlgen/handler"

	"github.com/rodrwan/medium-example/cmd/server/logger"
	"github.com/rodrwan/medium-example/database"
	"github.com/rodrwan/medium-example/graphql"
	"github.com/rodrwan/medium-example/service"

	_ "github.com/lib/pq"
)

func main() {
	postgresDSN := flag.String("postgres-dsn", "postgres://mediumexample:me1234@postgres_medium:5432/example?sslmode=disable", "Postgres domain service name")
	port := flag.Int("port", 8080, "server port")

	flag.Parse()

	db, err := database.NewPostgres(*postgresDSN)
	if err != nil {
		panic(err)
	}

	svc := service.NewService(db)
	logs := func(writer io.Writer, params handlers.LogFormatterParams) {
		l := &logger.Logger{
			StatusCode: params.StatusCode,
			Size:       params.Size,
			Method:     params.Request.Method,
			TimeStamp:  time.Now(),
			URL:        params.URL.String(),
		}

		json.NewEncoder(writer).Encode(l)
	}
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	rootHandler := handler.GraphQL(
		graphql.NewExecutableSchema(
			graphql.Config{
				Resolvers: &graphql.Resolver{
					Service: svc,
				},
			},
		),
	)
	http.Handle("/query", handlers.CustomLoggingHandler(os.Stdout, rootHandler, logs))

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
