package cmd

import (
	"context"
	"main/config"
	"main/internal/delivery"
	"main/internal/delivery/restful"
	"main/internal/repository"
	"main/internal/repository/conn"
	"main/internal/usecase"

	"go.uber.org/fx"
)

func injectInfra() fx.Option {
	return fx.Provide(
		context.Background,
		config.Load,
		conn.NewGormDB,
		conn.NewRedisClient,
		conn.NewDao,
	)
}

func injectRepository() fx.Option {
	return fx.Provide(
		repository.NewTokenRepository,
		repository.NewNotificationRepository,
		repository.NewOTPRepository,
		repository.NewUserRepository,
		repository.NewProductRepository,
		repository.NewUserCategoryPreferenceRepository,
	)
}

func injectUsecase() fx.Option {
	return fx.Provide(
		usecase.NewTokenUsecase,
		usecase.NewOTPUsecase,
		usecase.NewAuthUsecase,
		usecase.NewProductUsecase,
	)
}

func injectDelivery() fx.Option {
	return fx.Provide(
		restful.NewUserController,
		restful.NewProductController,
		delivery.RegisterRouters,
	)
}
