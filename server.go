package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/gommon/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/faruryo/toban-api/graph/generated"
	"github.com/faruryo/toban-api/graph/resolvers"
	"github.com/faruryo/toban-api/repository"
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
	if b, err := strconv.ParseBool(os.Getenv("DEBUG_ECHO")); err != nil {
		debugEcho = b
	}
	debugDb := false
	if b, err := strconv.ParseBool(os.Getenv("DEBUG_DB")); err != nil {
		debugDb = b
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

	logLevel := repository.Silent
	if debugDb {
		logLevel = repository.Info
	}
	db, err := repository.GetDb(makeDBConnectionOptionsByEnv(), logLevel)
	if err != nil {
		e.Logger.Fatal("failed to connect database: %v", err)
		return
	}

	gqlEp := "api/graphql"
	plgEp := "playground"
	e.POST("/"+gqlEp, func(c echo.Context) error {
		repo, err := repository.NewRepository(db)
		if err != nil {
			e.Logger.Fatalf("Failed to create repository : %s", err)
		}

		h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
			Resolvers:  &resolvers.Resolver{Repository: repo},
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

func makeDBConnectionOptionsByEnv() repository.DBConnectionOptions {
	opt := repository.DBConnectionOptions{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
		Host:     os.Getenv("MYSQL_HOST"),
	}
	return opt
}
