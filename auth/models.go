package auth

type JWTRequest struct {
	Login    string
	Password string
}

type JWTResponse struct {
	Type         string
	AccessToken  string
	RefreshToken string
}
