package utils

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tsanny/kenmeri/internal/config"
	"github.com/tsanny/kenmeri/internal/datasources/drivers"
)

func SetupDatabase() (*sqlx.DB, error) {
	// Setup Config Databse
	configDB := drivers.DBConfig{
		DriverName:     config.AppConfig.DBPostgreDriver,
		DataSourceName: config.AppConfig.DBPostgreDsn,
		MaxOpenConns:   100,
		MaxIdleConns:   10,
		MaxLifetime:    15 * time.Minute,
	}

	// Initialize Database driversSQL
	conn, err := configDB.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
