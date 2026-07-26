package dto

type UpdateDesignerProfileRequest struct {
	Nickname  string `json:"nickname" example:"Dan"`
	AvatarKey string `json:"avatarKey" example:"designer_profiles/avatars/1512"` // Нужно будет настроить работу с minio
}
