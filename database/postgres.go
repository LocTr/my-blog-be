package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func DBConnect() (*bun.DB, error) {

	user := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	database := viper.GetString("POSTGRES_DB")
	port := viper.GetString("POSTGRES_PORT")

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr("localhost:"+port),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(database),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := checkConn(db); err != nil {
		return nil, err
	}
	return db, nil
}

func checkConn(db *bun.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
