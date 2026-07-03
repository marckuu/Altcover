package tools

import (
	"fmt"
	"strconv"
)

func GetOffsetAndLimitFromQuery(offset string, limit string) (int, int, error) {
	offst, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка получения offset из query параметра %w", err)
	}
	lmt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка получения limit из query параметра %w", err)
	}

	return int(offst), int(lmt), nil
}
