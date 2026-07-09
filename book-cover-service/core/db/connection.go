package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DB_CONN_PATH"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
