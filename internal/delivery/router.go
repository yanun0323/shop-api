package delivery

import (
	"context"
	"main/internal/delivery/middleware"
	"main/internal/delivery/restful"
	"main/internal/domain/usecase"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type RegisterRoutersParam struct {
	fx.In

	Auth    restful.AuthController
	Product restful.ProductController

	TokenUsecase usecase.TokenUsecase
}

func RegisterRouters(ctx context.Context, ctr RegisterRoutersParam) *echo.Echo {
	server := echo.New()

	v1 := server.Group("/api/v1")
	v1.Use(
		echomiddleware.Logger(),
		echomiddleware.Recover(),
		echomiddleware.CORS(),
		middleware.RateLimiter(),
	)

	public := v1.Group("")
	login := v1.Group("", middleware.Token(ctr.TokenUsecase))

	public.POST("/auth", ctr.Auth.Register)
	public.POST("/auth/otp", ctr.Auth.SendVerifyCode)
	public.POST("/auth/token", ctr.Auth.Login)

	login.PUT("/auth/token", ctr.Auth.RefreshToken)
	login.DELETE("/auth/token", ctr.Auth.Logout)

	login.GET("/product/recommendation", ctr.Product.ListRecommendation)

	return server
}
