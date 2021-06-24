package db

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/InfoBlogFriends/server/config"
	"go.uber.org/zap"
)

func NewDB(logger *zap.Logger, dbConf config.DB) (*sqlx.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = dbConf.Net
	cfg.Addr = dbConf.Host
	cfg.User = dbConf.Username
	cfg.Passwd = dbConf.Password
	cfg.DBName = dbConf.DBName
	cfg.ParseTime = true
	cfg.Timeout = 2 * time.Second

	dsn := cfg.FormatDSN()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.After(time.Duration(dbConf.Timeout) * time.Second)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("db connection timeout, after %d second", dbConf.Timeout)
		case <-ticker.C:
			db, err := connectDB(dbConf.Driver, dsn)
			if err == nil {
				return db, nil
			}
			logger.Error("failed to connect to database", zap.String("dsn", dsn), zap.Error(err))
		}
	}
}

func connectDB(driverName string, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return db, nil
}
