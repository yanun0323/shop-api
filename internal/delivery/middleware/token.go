package middleware

import (
	"main/internal/delivery/payload"
	"main/internal/delivery/response"
	"main/internal/domain/usecase"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/yanun0323/pkg/logs"
)

const (
	_authorizationHeaderKey = "Authorization"
	_bearer                 = "Bearer "
)

func Token(tokenUsecase usecase.TokenUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerToken := c.Request().Header.Get(_authorizationHeaderKey)
			if len(bearerToken) == 0 {
				return c.JSON(http.StatusUnauthorized, response.MsgErr("no token provided"))
			}

			token, ok := strings.CutPrefix(bearerToken, _bearer)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.MsgErr("wrong token format"))
			}

			ctx := c.Request().Context()
			claims, err := tokenUsecase.VerifyToken(ctx, token)
			if err != nil {
				if errors.Is(err, usecase.ErrTokenExpired) {
					return c.JSON(http.StatusUnauthorized, response.MsgErr("token expired"))
				}

				return c.JSON(http.StatusUnauthorized, response.Err(err, "invalid token provided"))
			}

			clientDeviceID := payload.GetDeviceID(c.Request())

			logs.Debugf("client device id: %s", clientDeviceID)
			logs.Debugf("%+v", claims)

			if !strings.EqualFold(claims.DeviceID, clientDeviceID) {
				return c.JSON(http.StatusUnauthorized, response.MsgErr("invalid token provided", "mismatch deviceID, client: %s, token: %s", clientDeviceID, claims.DeviceID))
			}

			ctx = payload.SetToken(ctx, token)
			ctx = payload.SetUserID(ctx, claims.UserID)
			ctx = payload.SetDeviceID(ctx, claims.DeviceID)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
