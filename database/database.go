package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type DBInfo struct {
	Host     string
	Username string
	Password string
}

type IDBHandler interface {
	pgxtype.Querier
}

func InitDB(info *DBInfo) (*pgxpool.Pool, error) {
	log.Debug("initializing database")

	// initialize configuration
	connectionString := fmt.Sprintf("host=%s user=%s password=%s", info.Host, info.Username, info.Password)
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// initialize connection pool
	connectionAttempts := 0
	var db *pgxpool.Pool
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		db, err = pgxpool.ConnectConfig(ctx, poolConfig)
		if err == nil {
			break
		}
		connectionAttempts++
		if connectionAttempts >= 10 {
			return nil, err
		}
		cancel()
		// retry db
		log.WithError(err).Error("database connection attempt failed, waiting 5s then retrying")
		time.Sleep(time.Second * 5)
	}

	log.Debug("database initialized")
	return db, nil
}
