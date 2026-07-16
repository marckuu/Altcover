package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Database interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func CreateConnection(ctx context.Context, dbConnPath string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dbConnPath)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
