package dto

type CreteDesignerProfileRequest struct {
	Nickname  string `json:"nickname" example:"Tom"`
	AvatarKey string `json:"avatarKey" example:"designer_profile/avatars/1542"` // Нужно будет настроить работу с minio
}
