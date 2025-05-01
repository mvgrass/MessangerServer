package auth

type RegisterRequestDto struct {
	Name  string `json: "name"`
	Email string `json: "email"`
	// password
}

type RegisterResponseDto struct {
	AccessToken    string `json: "accessToken"`
	RefreshToken   string `json: "refreshToken"`
	TimeExparation uint64 `json: "accessToken"`
	// password
}
