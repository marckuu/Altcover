package repositories

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var coonPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(os.Getenv("DB_CONN_PATH_TEST_EXT"))
	if err != nil {
		fmt.Println("не удалось создать конфиг:", err)
		os.Exit(1)
	}
	config.MaxConns = 6

	coonPool, err = pgxpool.New(ctx, config.ConnString())
	if err != nil {
		fmt.Println("не удалось создать connection pool:", err)
		os.Exit(1)
	}

	if err = coonPool.Ping(ctx); err != nil {
		fmt.Println("connection pool недоступен:", err)
		coonPool.Close()
		os.Exit(1)
	}

	code := m.Run()

	coonPool.Close()
	os.Exit(code)
}
