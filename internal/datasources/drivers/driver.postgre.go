package drivers

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tsanny/kenmeri/internal/constants"
	"github.com/tsanny/kenmeri/pkg/logger"
)

type DBConfig struct {
	DriverName     string
	DataSourceName string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetime    time.Duration
}

func (config *DBConfig) InitializeDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// set maximum number of open connections to database
	logger.Info(fmt.Sprintf("Setting maximum number of open connections to %d", config.MaxOpenConns), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
	db.SetMaxOpenConns(config.MaxOpenConns)

	// set maximum number of idle connections in the pool
	logger.Info(fmt.Sprintf("Setting maximum number of idle connections to %d", config.MaxIdleConns), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
	db.SetMaxIdleConns(config.MaxIdleConns)

	// set maximum time to wait for new connection
	logger.Info(fmt.Sprintf("Setting maximum lifetime for a connection to %s", config.MaxLifetime), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
	db.SetConnMaxLifetime(config.MaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}
