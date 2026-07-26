package dto

type LoginRequest struct {
	Nickname string `json:"nickname" example:"Paul"`
	Password string `json:"password" example:"1234#678"`
}
