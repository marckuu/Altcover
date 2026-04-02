package domains

import "time"

type Order struct {
	ID        int64
	Status    int // enum из файла status
	CreatedAt time.Time

	UserID          int64
	OrdersHistoryID int64
}
