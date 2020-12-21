package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	port := viper.GetString("port")
	if port == "" {
		port = defaultPort
	}
	debugEcho := false
	if b, err := strconv.ParseBool(viper.GetString("debug.echo")); err != nil {
		debugEcho = b
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

	db, err := connectDB()
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

func connectDB() (*gorm.DB, error) {
	debugDb := false
	if b, err := strconv.ParseBool(viper.GetString("debug.db")); err != nil {
		debugDb = b
	}
	logLevel := logger.Silent
	if debugDb {
		logLevel = logger.Info
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.database"),
	)

	var retryDuration time.Duration
	maxRetryNumber := 4
	var db *gorm.DB
	var err error
	for i := 0; i < maxRetryNumber; i++ {
		db, err = gorm.Open(
			mysql.Open(dsn),
			&gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
				Logger:                 logger.Default.LogMode(logLevel),
			},
		)
		if err == nil {
			break
		}
		log.Print(err)
		retryDuration = time.Duration(i*2) * time.Second
		log.Printf("issue connecting to database, retrying. retryNumber:%d, retryDuration:%s", i, retryDuration)
		if i != maxRetryNumber-1 {
			time.Sleep(retryDuration)
		}
	}
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		log.Printf("dialector: %s", dsn)
		return nil, err
	}

	return db, nil
}
