package di

import (
	"apg105/auth"
	"apg105/datasource"
	gameRepo "apg105/datasource/game"
	userRepo "apg105/datasource/user"
	gameDomain "apg105/domain/game"
	userDomain "apg105/domain/user"
	authWeb "apg105/web/auth"
	gameWeb "apg105/web/game"
	"apg105/web/router"
	userWeb "apg105/web/user"

	"go.uber.org/fx"
)

var GameModule = fx.Module("Game",
	fx.Provide(
		fx.Annotate(
			gameRepo.NewPGRepo,
			fx.As(new(gameDomain.Repo)),
		),
		fx.Annotate(
			gameDomain.NewMinimax,
			fx.As(new(gameDomain.AI)),
		),
		fx.Annotate(
			gameDomain.NewService,
			fx.As(new(gameWeb.Service)),
		),
		gameWeb.NewHandler,
	),
)

var UserModule = fx.Module("User",
	fx.Provide(
		fx.Annotate(
			userRepo.NewPGRepo,
			fx.As(new(userDomain.Repo)),
		),
		fx.Annotate(
			userDomain.NewService,
			fx.As(new(userWeb.Service)),
			fx.As(new(auth.UserService)),
		),
		userWeb.NewHandler,
	),
)

var AuthModule = fx.Module("Auth",
	fx.Provide(
		fx.Annotate(
			auth.NewJwtProvider,
			fx.As(new(auth.TokenManager)),
			fx.As(new(auth.TokenParser)),
		),
		fx.Annotate(
			auth.NewAuthService,
			fx.As(new(authWeb.AuthService)),
		),
		authWeb.NewHandler,
		fx.Annotate(
			authWeb.NewAuthenticator,
			fx.As(new(router.AuthMiddleware)),
		),
	),
)

var RouterModule = fx.Module("Router",
	fx.Provide(
		router.NewRouter,
		router.NewServer,
	),
	fx.Invoke(router.StartServer),
)

var CommonModule = fx.Module("xoCommon",
	fx.Provide(datasource.InitDBPool),
	GameModule,
	UserModule,
	AuthModule,
	RouterModule,
)
