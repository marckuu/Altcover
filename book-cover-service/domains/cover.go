package domains

type Cover struct {
	ID          string
	Title       string
	Description string
	Likes       int
	ImagesKeys  []string // Нужно заменить на работу с minio
	Status      int      // Должен браться status из файла status

	DesignerID        string
	DesignerNickname  string
	DesignerAvatarKey string // Заменить на подтягивание из minio

	BookID string
}
