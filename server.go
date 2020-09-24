package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/repositories"
	"github.com/faruryo/toban-api/resolvers"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultPort = "8080"
const dsn = "toban:toban@tcp(toban-mysql:3306)/toban?charset=utf8mb4&parseTime=true"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.Debug()

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers:  &resolvers.Resolver{TobanRepository: repositories.NewTobanRepository(db)},
			Directives: generated.DirectiveRoot{},
			Complexity: generated.ComplexityRoot{},
		}),
	)

	srv.Use(extension.FixedComplexityLimit(300))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
