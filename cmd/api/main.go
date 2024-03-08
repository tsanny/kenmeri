package main

import (
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/tsanny/kenmeri/cmd/api/server"
	"github.com/tsanny/kenmeri/internal/config"
	"github.com/tsanny/kenmeri/internal/constants"
	"github.com/tsanny/kenmeri/pkg/logger"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{
			constants.LoggerCategory: constants.LoggerCategoryConfig,
		})
	}
	logger.Info("configuration loaded", logrus.Fields{
		constants.LoggerCategory: constants.LoggerCategoryConfig,
	})
}

func main() {
	numCPU := runtime.NumCPU()
	logger.InfoF(
		"The project is running on %d CPU(s)",
		logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig},
		numCPU,
	)
	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(numCPU / 2)
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{
			constants.LoggerCategory: constants.LoggerCategoryServer,
		})
	}
	if err := app.Run(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{
			constants.LoggerCategory: constants.LoggerCategoryServer,
		})
	}
}
