package domains

type DesignerProfile struct {
	ID        int64
	AvatarKey string // Нужно будет настроить работу с minio

	UserID int64
}
