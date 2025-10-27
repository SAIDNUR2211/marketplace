package db

import (
	"fmt"
	"marketplace/internal/configs"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitConnection() (*sqlx.DB, error) {
	connectionConfigs := configs.AppSettings.PostgresParams
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", connectionConfigs.Host, connectionConfigs.Port, connectionConfigs.User, os.Getenv("DB_PASSWORD"), connectionConfigs.Database)
	dbConn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
func CloseConnection(db *sqlx.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}
