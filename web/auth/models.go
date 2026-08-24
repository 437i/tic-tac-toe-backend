package auth

type JWTRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshJWTRequest struct {
	RefreshToken string `json:"refreshToken"`
}
