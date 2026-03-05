package dto

import "github.com/google/uuid"

type JwtRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtResponse struct {
	Token        string `json:"type" binding:"required"`
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UserInfoResponse struct {
	ID    uuid.UUID `json:"id"`
	Login string    `json:"login"`
}
