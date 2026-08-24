package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ Repo = (*mockRepo)(nil)

type mockRepo struct {
	createCalled bool
	createdUser  User
	createErr    error

	userByLoginCalled bool
	userByLoginLogin  string
	userByLogin       User
	userByLoginErr    error

	loginByIDCalled bool
	loginByIDUserID uuid.UUID
	loginByID       string
	loginByIDErr    error
}

func (m *mockRepo) Create(ctx context.Context, userDm User) error {
	m.createCalled = true
	m.createdUser = userDm

	return m.createErr
}

func (m *mockRepo) GetUserByLogin(ctx context.Context, login string) (User, error) {
	m.userByLoginCalled = true
	m.userByLoginLogin = login

	return m.userByLogin, m.userByLoginErr
}

func (m *mockRepo) GetLoginByID(ctx context.Context, userID uuid.UUID) (string, error) {
	m.loginByIDCalled = true
	m.loginByIDUserID = userID

	return m.loginByID, m.loginByIDErr
}

func TestNewUser(t *testing.T) {
	req := SignUpRequest{
		Login:    "alex",
		Password: "pass",
	}

	got := NewUser(req)

	if got.Login != req.Login {
		t.Errorf("Login = %q, want %q", got.Login, req.Login)
	}

	if got.Password != req.Password {
		t.Errorf("Password = %q, want %q", got.Password, req.Password)
	}

	if got.UserID == uuid.Nil {
		t.Error("UserID is nil, want non-nil UUID")
	}
}

func TestRegister(t *testing.T) {
	type want struct {
		createCalled bool
		err          error
		repoErr      error
	}

	login := "alex"
	password := "pass"
	repoErr := errors.New("repo error")

	tests := []struct {
		name string
		req  SignUpRequest
		want want
	}{
		{
			name: "valid registration",
			req: SignUpRequest{
				Login:    login,
				Password: password,
			},
			want: want{
				createCalled: true,
				err:          nil,
			},
		}, {
			name: "empty login",
			req: SignUpRequest{
				Login:    "",
				Password: password,
			},
			want: want{
				createCalled: false,
				err:          ErrEmptyLogpass,
			},
		}, {
			name: "empty password",
			req: SignUpRequest{
				Login:    login,
				Password: "",
			},
			want: want{
				createCalled: false,
				err:          ErrEmptyLogpass,
			},
		}, {
			name: "empty login & password",
			req: SignUpRequest{
				Login:    "",
				Password: "",
			},
			want: want{
				createCalled: false,
				err:          ErrEmptyLogpass,
			},
		}, {
			name: "password too long (len 73)",
			req: SignUpRequest{
				Login:    login,
				Password: strings.Repeat("1", 73),
			},
			want: want{
				createCalled: false,
				err:          ErrPassTooLong,
			},
		}, {
			name: "max password len (len 72)",
			req: SignUpRequest{
				Login:    login,
				Password: strings.Repeat("1", 72),
			},
			want: want{
				createCalled: true,
				err:          nil,
			},
		}, {
			name: "repo error",
			req: SignUpRequest{
				Login:    login,
				Password: password,
			},
			want: want{
				createCalled: true,
				err:          repoErr,
				repoErr:      repoErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				createErr: tt.want.repoErr,
			}
			service := NewService(repo)

			err := service.Register(context.Background(), tt.req)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("Err = %v, want %v", err, tt.want.err)
			}

			if repo.createCalled != tt.want.createCalled {
				t.Errorf("Create called = %v, want %v", repo.createCalled, tt.want.createCalled)
			}

			if tt.want.createCalled {
				if repo.createdUser.UserID == uuid.Nil {
					t.Error("UserID is nil, want non-nil UUID")
				}

				if repo.createdUser.Login != tt.req.Login {
					t.Errorf("Login = %q, want %q", repo.createdUser.Login, tt.req.Login)
				}

				if err := bcrypt.CompareHashAndPassword(
					[]byte(repo.createdUser.Password),
					[]byte(tt.req.Password),
				); err != nil {
					t.Errorf("Compare hash and password failed: (%v)", err)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type want struct {
		userByLoginCalled bool
		userByLoginLogin  string
		user              SafeUser
		err               error
		repoErr           error
	}

	id := uuid.New()

	validLogin := "alex"
	invalidLogin := "bubba"

	validPass := "pass"
	invalidPass := "nopass"

	validHash, err := bcrypt.GenerateFromPassword([]byte(validPass), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	repoErr := errors.New("repo err")

	tests := []struct {
		name     string
		login    string
		password string
		want     want
	}{
		{
			name:     "valid login",
			login:    validLogin,
			password: validPass,
			want: want{
				userByLoginCalled: true,
				userByLoginLogin:  validLogin,
				user: SafeUser{
					UserID: id,
					Login:  validLogin,
				},
			},
		}, {
			name:     "empty login",
			login:    "",
			password: validPass,
			want: want{
				userByLoginCalled: false,
				user:              SafeUser{},
				err:               ErrEmptyLogpass,
			},
		}, {
			name:     "empty password",
			login:    validLogin,
			password: "",
			want: want{
				userByLoginCalled: false,
				user:              SafeUser{},
				err:               ErrEmptyLogpass,
			},
		}, {
			name:     "empty login & password",
			login:    "",
			password: "",
			want: want{
				userByLoginCalled: false,
				user:              SafeUser{},
				err:               ErrEmptyLogpass,
			},
		}, {
			name:     "repo error",
			login:    validLogin,
			password: validPass,
			want: want{
				userByLoginCalled: true,
				userByLoginLogin:  validLogin,
				user:              SafeUser{},
				err:               repoErr,
				repoErr:           repoErr,
			},
		}, {
			name:     "wrong password",
			login:    validLogin,
			password: invalidPass,
			want: want{
				userByLoginCalled: true,
				userByLoginLogin:  validLogin,
				user:              SafeUser{},
				err:               ErrInvalidCredentials,
			},
		}, {
			name:     "wrong login",
			login:    invalidLogin,
			password: validPass,
			want: want{
				userByLoginCalled: true,
				userByLoginLogin:  invalidLogin,
				user:              SafeUser{},
				err:               ErrInvalidCredentials,
				repoErr:           ErrUserNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				userByLogin: User{
					UserID:   id,
					Login:    validLogin,
					Password: string(validHash),
				},
				userByLoginErr: tt.want.repoErr,
			}
			service := NewService(repo)

			user, err := service.Login(context.Background(), tt.login, tt.password)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("Err = %v, want %v", err, tt.want.err)
			}

			if repo.userByLoginCalled != tt.want.userByLoginCalled {
				t.Errorf(
					"GetUserByLogin called = %v, want %v",
					repo.userByLoginCalled,
					tt.want.userByLoginCalled,
				)
			}

			if tt.want.userByLoginCalled {
				if repo.userByLoginLogin != tt.login {
					t.Errorf("Login = %q, want %q", repo.userByLoginLogin, tt.login)
				}
			}

			if user != tt.want.user {
				t.Errorf("User = %q, want %q", user, tt.want.user)
			}
		})
	}
}

func TestGetUserLogin(t *testing.T) {
	type want struct {
		loginByIDCalled bool
		loginByIDUserID uuid.UUID
		err             error
		repoErr         error
	}

	id := uuid.New()
	repoErr := errors.New("repo error")

	tests := []struct {
		name   string
		userID uuid.UUID
		want   want
	}{
		{
			name:   "success",
			userID: id,
			want: want{
				loginByIDCalled: true,
				loginByIDUserID: id,
			},
		}, {
			name:   "repo error",
			userID: id,
			want: want{
				loginByIDCalled: true,
				loginByIDUserID: id,
				err:             repoErr,
				repoErr:         repoErr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				loginByID:    "alex",
				loginByIDErr: tt.want.repoErr,
			}
			service := NewService(repo)

			login, err := service.GetUserLogin(context.Background(), tt.userID)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("Err = %v, want %v", err, tt.want.err)
			}

			if tt.want.loginByIDCalled != repo.loginByIDCalled {
				t.Errorf(
					"GetLoginByID called = %v, want %v",
					repo.loginByIDCalled,
					tt.want.loginByIDCalled,
				)
			}

			if tt.want.loginByIDCalled {
				if repo.loginByIDUserID != tt.userID {
					t.Errorf("UserID = %q, want %q", repo.loginByIDUserID, tt.userID)
				}

				if login != repo.loginByID {
					t.Errorf("Login = %q, want %q", login, repo.loginByID)
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	type want struct {
		loginByIDCalled bool
		loginByIDUserID uuid.UUID
		user            SafeUser
		err             error
		repoErr         error
	}

	id := uuid.New()
	repoErr := errors.New("repo err")

	tests := []struct {
		name   string
		userID uuid.UUID
		want   want
	}{
		{
			name:   "success",
			userID: id,
			want: want{
				loginByIDCalled: true,
				loginByIDUserID: id,
				user: SafeUser{
					UserID: id,
					Login:  "alex",
				},
			},
		}, {
			name:   "repo error",
			userID: id,
			want: want{
				loginByIDCalled: true,
				loginByIDUserID: id,
				user:            SafeUser{},
				err:             repoErr,
				repoErr:         repoErr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				loginByID:    "alex",
				loginByIDErr: tt.want.repoErr,
			}
			service := NewService(repo)

			user, err := service.GetUserByID(context.Background(), tt.userID)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("Err = %v, want %v", err, tt.want.err)
			}

			if repo.loginByIDCalled != tt.want.loginByIDCalled {
				t.Errorf(
					"GetLoginByID called = %v, want %v",
					repo.loginByIDCalled,
					tt.want.loginByIDCalled,
				)
			}

			if tt.want.loginByIDCalled {
				if repo.loginByIDUserID != tt.userID {
					t.Errorf("UserID = %q, want %q", repo.loginByIDUserID, tt.userID)
				}
			}

			if user != tt.want.user {
				t.Errorf("User = %q, want %q", user, tt.want.user)
			}
		})
	}
}
