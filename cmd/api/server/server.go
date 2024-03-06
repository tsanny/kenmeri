package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsanny/kenmeri/internal/config"
	"github.com/tsanny/kenmeri/internal/constants"
	"github.com/tsanny/kenmeri/internal/datasources/drivers"
	"github.com/tsanny/kenmeri/internal/http/middlewares"
	"github.com/tsanny/kenmeri/internal/http/routes"
	"github.com/tsanny/kenmeri/pkg/logger"
	"gorm.io/gorm"
)

type App struct {
	HttpServer *http.Server
}

func NewApp() (*App, error) {
	// setup databases
	_, err := setupDatabse()
	if err != nil {
		return nil, err
	}

	// setup router
	router := setupRouter()

	// // jwt service
	// jwtService := jwt.NewJWTService()

	// // cache
	// redisCache := caches.NewRedisCache(config.AppConfig.REDISHost, 0, config.AppConfig.REDISPassword, time.Duration(config.AppConfig.REDISExpired))
	// ristrettoCache, err := caches.NewRistrettoCache()
	// if err != nil {
	// 	panic(err)
	// }

	// // user middleware
	// authMiddleware := middlewares.NewAuthMiddleware(jwtService, false)
	// // admin middleware
	// authAdminMiddleware := middlewares.NewAuthMiddleware(jwtService, true)

	// Routes
	router.GET("/", routes.RootHandler)
	// routes.NewUsersRoute(conn, jwtService, redisCache, ristrettoCache, router, authMiddleware).UsersRoute()

	// setup http server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &App{
		HttpServer: server,
	}, nil
}

func (a *App) Run() (err error) {
	// Gracefull Shutdown
	go func() {
		logger.Info(fmt.Sprintf("success to listen and serve on :%d", config.AppConfig.Port), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryInit})
		if err := a.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	<-quit
	logger.Info("Shutdown server ...", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryClose})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = a.HttpServer.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Error during server shutdown: %+v", err)
	}

	return nil
}

func setupDatabse() (*gorm.DB, error) {
	// Setup Config Databse
	configDB := drivers.ConfigPostgreSQL{
		DB_Username: config.AppConfig.DBUsername,
		DB_Password: config.AppConfig.DBPassword,
		DB_Host:     config.AppConfig.DBHost,
		DB_Port:     config.AppConfig.DBPort,
		DB_Database: config.AppConfig.DBDatabase,
		DB_DSN:      config.AppConfig.DBDsn,
	}

	// Initialize Database driversSQL
	conn, err := configDB.InitializeDatabasePostgreSQL()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func setupRouter() *gin.Engine {
	// set the runtime mode
	var mode = gin.ReleaseMode
	// if config.AppConfig.Debug {
	// 	mode = gin.DebugMode
	// }
	gin.SetMode(mode)

	// create a new router instance
	router := gin.New()

	// set up middlewares
	router.Use(middlewares.CORSMiddleware())
	if mode == gin.DebugMode {
		router.Use(gin.LoggerWithFormatter(logger.CustomLogFormatter))
	}
	router.Use(gin.Recovery())

	return router
}
