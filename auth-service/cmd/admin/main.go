package main

import (
	"auth-service/core/db"
	"auth-service/core/enums"
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	nickname := flag.String("nickname", "", "Admin nickname")
	flag.Parse()

	if *nickname == "" {
		fmt.Println("не указан никнейм админа")
		return
	}

	ctx := context.Background()

	connPool, err := db.CreateConnPool(ctx, os.Getenv("DB_CONN_PATH"))
	if err != nil {
		fmt.Println("не удалось создать connection poll:", err)
		return
	}

	query := `
	UPDATE users
	SET role = $1
	WHERE nickname = $2;
`

	tag, err := connPool.Exec(ctx, query, enums.Admin, nickname)
	if err != nil {
		fmt.Println("не удалось изменить роль на админа:", err)
		return
	}

	if tag.RowsAffected() == 0 {
		fmt.Println("пользователь не найден:", err)
		return
	}

	fmt.Println("администратор успешно установлен")
}
