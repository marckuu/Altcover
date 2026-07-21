package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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

func CreateConnPool(ctx context.Context, connPath string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connPath)
	if err != nil {
		fmt.Println("не удалось создать конфиг:", err)
		return nil, err
	}
	config.MaxConns = 6

	coonPool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		fmt.Println("не удалось создать connection pool:", err)
		return nil, err
	}

	if err = coonPool.Ping(ctx); err != nil {
		fmt.Println("connection pool недоступен:", err)
		coonPool.Close()
		return nil, err
	}

	return coonPool, nil
}
