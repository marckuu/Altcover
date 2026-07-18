package repositories

import (
	"auth-service/core/db"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	conn, err = db.CreateConnection(ctx, os.Getenv("DB_CONN_PATH_TEST_EXT"))
	if err != nil {
		fmt.Println("не удалось подключиться к бд:", err)
		os.Exit(1)
	}

	code := m.Run()

	conn.Close(ctx)
	os.Exit(code)
}
