package user

import (
	dmu "apg105/domain/user"
)

func (s SignUpRequest) toDomain() dmu.SignUpRequest {
	return dmu.SignUpRequest{
		Login:    s.Login,
		Password: s.Password,
	}
}

func toMe(user dmu.SafeUser) MeResponse {
	return MeResponse{
		UserID: user.UserID,
		Login:  user.Login,
	}
}
