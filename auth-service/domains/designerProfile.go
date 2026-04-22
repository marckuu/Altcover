package domains

import "github.com/jackc/pgx/v5/pgtype"

type DesignerProfile struct {
	ID        pgtype.UUID
	AvatarKey string // Нужно будет настроить работу с minio

	UserID int64
}
