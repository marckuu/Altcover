package domains

import (
	"github.com/google/uuid"
)

type Cover struct {
	ID          uuid.UUID
	Title       string
	Description string
	ImagesKeys  []string // Нужно заменить на работу с minio
	Status      int      // Должен браться status из файла status

	UserID            uuid.UUID
	DesignerNickname  string
	DesignerAvatarKey string // Заменить на подтягивание из minio

	BookID uuid.UUID
}
