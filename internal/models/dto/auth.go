package dto

type AuthResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=30"`
	Password    string `json:"password" binding:"required,min=6"`
	Email       string `json:"email" binding:"omitempty,email"`
	TelegramID  int64  `json:"telegram_id" binding:"omitempty"`
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	Fingerprint  string `json:"fingerprint" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	Fingerprint  string `json:"fingerprint" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
