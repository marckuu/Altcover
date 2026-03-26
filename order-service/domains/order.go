package domains

import "time"

type Order struct {
	id        int64
	status    int // enum из файла status
	createdAt time.Time

	userId   int64
	coversId []int64
}
