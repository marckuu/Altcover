package domains

import (
	"github.com/google/uuid"
)

type Cover struct {
	ID          uuid.UUID
	Title       string
	Description string
	Likes       int
	ImagesKeys  []string // Нужно заменить на работу с minio
	Status      int      // Должен браться status из файла status

	DesignerID        uuid.UUID
	DesignerNickname  string
	DesignerAvatarKey string // Заменить на подтягивание из minio

	BookID uuid.UUID
}
