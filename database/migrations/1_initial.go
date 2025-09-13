package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

const postTable = `
CREATE TABLE IF NOT EXISTS posts (
id SERIAL PRIMARY KEY,
title VARCHAR(255) NOT NULL,
content TEXT NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

func init() {
	up := []string{
		postTable,
	}

	down := []string{
		`DROP TABLE IF EXISTS posts`,
	}

	Migrations.MustRegister(func(context context.Context, db *bun.DB) error {
		fmt.Println("creating initial tables")
		for _, query := range up {
			_, err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("dropping initial tables")
		for _, query := range down {
			_, err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
