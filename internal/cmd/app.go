package cmd

import (
	"context"
	"fmt"
	"main/config"

	"github.com/labstack/echo/v4"
	"github.com/yanun0323/pkg/logs"
	"go.uber.org/fx"
)

func Run() {
	logs.SetDefaultLevel(logs.NewLevel(config.Load().Log.Level))

	fx.New(
		injectInfra(),
		injectRepository(),
		injectUsecase(),
		injectDelivery(),
		fx.Invoke(start),
	).Run()
}

func start(ctx context.Context, conf config.Config, server *echo.Echo) error {
	return server.Start(fmt.Sprintf(":%s", conf.Http.Port))
}
