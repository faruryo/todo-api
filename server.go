package main

import (
	"net/http"
	"os"

	"github.com/labstack/gommon/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/faruryo/toban-api/databases"
	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/repositories"
	"github.com/faruryo/toban-api/resolvers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const defaultPort = "8080"

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

	if debugEcho {
		e.Use(middleware.Logger())
		e.Logger.SetLevel(log.DEBUG)
	}

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	logLevel := databases.Silent
	if debugDb {
		logLevel = databases.Info
	}
	db, err := databases.GetDbByEnv(logLevel)
	if err != nil {
		e.Logger.Fatal("failed to connect database: %v", err)
		return
	}

	gqlEp := "api/graphql"
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
