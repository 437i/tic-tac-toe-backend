package auth

import "apg105/auth"

func (j JWTRequest) toDomain() auth.JWTRequest {
	return auth.JWTRequest{
		Login:    j.Login,
		Password: j.Password,
	}
}

func toWebResponse(resp auth.JWTResponse) JWTResponse {
	return JWTResponse{
		Type:         resp.Type,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}
}
