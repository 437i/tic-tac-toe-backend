package user

import dmu "apg105/domain/user"

func toRepoUser(user dmu.User) User {
	return User{
		UserID:   user.UserID,
		Login:    user.Login,
		Password: user.Password,
	}
}

func (u User) toDomain() dmu.User {
	return dmu.User{
		UserID:   u.UserID,
		Login:    u.Login,
		Password: u.Password,
	}
}
