package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/faruryo/toban-api/databases"
	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/repositories"
	"github.com/faruryo/toban-api/resolvers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const defaultPort = "8080"
const dsn = "toban:toban@tcp(toban-mysql:3306)/toban?charset=utf8mb4&parseTime=true"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	debugEcho := false
	if os.Getenv("DEBUG_ECHO") != "" {
		debugEcho = true
	}
	debugDb := false
	if os.Getenv("DEBUG_DB") != "" {
		debugDb = true
	}

	e := echo.New()

	e.Use(middleware.Recover())
	if debugEcho {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Gzip())

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	logLevel := databases.Silent
	if debugDb {
		logLevel = databases.Info
	}
	db, err := databases.GetDbByDsn(dsn, logLevel)
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		return
	}

	gqlEp := "graphql"
	plgEp := "playground"
	e.POST("/"+gqlEp, func(c echo.Context) error {
		h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
			Resolvers:  &resolvers.Resolver{TobanRepository: repositories.NewTobanRepository(db)},
			Directives: generated.DirectiveRoot{},
			Complexity: generated.ComplexityRoot{},
		}))
		h.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	e.GET("/"+plgEp, func(c echo.Context) error {
		h := playground.Handler("GraphQL playground", "/"+gqlEp)
		h.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	e.HideBanner = true
	e.Logger.Infof("connect to http://localhost:%s/%s for GraphQL playground", port, plgEp)
	e.Logger.Fatal(e.Start(":" + port))
}
