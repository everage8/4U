package dto

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	Admin AdminPublic `json:"admin"`
}

type AdminPublic struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}
