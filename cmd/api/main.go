package main

import (
	"fmt"
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
			constants.LoggerCategory: constants.LoggerCategoryInit,
		})
	}
	logger.Info("configuration loaded", logrus.Fields{
		constants.LoggerCategory: constants.LoggerCategoryInit,
	})
}

func main() {
	numCPU := runtime.NumCPU()
	logger.Info(fmt.Sprintf("The project is running on %d CPU(s)", numCPU),
		logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDevice})

	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(numCPU / 2)
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{
			constants.LoggerCategory: constants.LoggerCategoryInit,
		})
	}
	if err := app.Run(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{
			constants.LoggerCategory: constants.LoggerCategoryClose,
		})
	}
}
