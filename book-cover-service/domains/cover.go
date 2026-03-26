package domains

type Cover struct {
	id          int64
	title       string
	description string
	images      []string // Нужно заменить на работу с minio
	status      int      // Должен браться status из файла status

	designerId       int64
	designerNickname string
	designerAvatar   string // Заменить на подтягивание из minio

	ratingId int64
}
