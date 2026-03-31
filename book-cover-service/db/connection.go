package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:12345@localhost:5432/postgres")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
