package db

import (
	"PullRequestService/pkg/config"
	"PullRequestService/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DataBase struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func New() *DataBase {
	log := logger.New()
	err := config.InitConfig()
	if err != nil {
		log.Error(err)
	}
	conn, err := pgxpool.New(context.Background(), config.GetString("DATABASE_URL"))
	if err != nil {
		log.Errorf("Connection to the database failed: %v\n", err)
		return nil
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Errorf("Database ping failed: %v\n", err)
		return nil
	}
	log.Info("Succesfully connect to the database\n")
	return &DataBase{db: conn, log: log}
}

func (d *DataBase) Close() {
	if d.db != nil {
		d.db.Close()
		d.log.Info("Database connection closed")
	}
}
