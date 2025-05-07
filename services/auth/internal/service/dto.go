package service

type RegisterRequestDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequestDto struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserInfoResponseDto struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RefreshRequestDto struct {
	RefreshToken string `json:"refreshToken"`
}

type JwtTokenRespnseDto struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	AccessExpiresIn  int64  `json:"accessExpiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}
