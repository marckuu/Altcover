package domains

type Cover struct {
	ID          int64
	Title       string
	Description string
	Likes       int
	ImagesKeys  []string // Нужно заменить на работу с minio
	Status      int      // Должен браться status из файла status

	DesignerID        int64
	DesignerNickname  string
	DesignerAvatarKey string // Заменить на подтягивание из minio

	BookID int64
}
