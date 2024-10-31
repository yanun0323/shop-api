package restful

import (
	"main/internal/delivery/payload"
	"main/internal/delivery/response"
	"main/internal/domain/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
}

func NewUserController(auth usecase.AuthUsecase) AuthController {
	return AuthController{
		authUsecase: auth,
	}
}

func (ctr *AuthController) Register(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     int    `json:"code"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.Err(err, "register", "request: %+v", req),
		)
	}

	ctx := c.Request().Context()
	err := ctr.authUsecase.Register(ctx, req.Email, req.Password, strconv.Itoa(req.Code))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.Err(err, "register", "request: %+v", req),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.Msg("register success"),
	)
}

func (ctr *AuthController) SendVerifyCode(c echo.Context) error {
	req := struct {
		Email string `json:"email"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.Err(err, "invalidate", "request: %+v", req),
		)
	}

	ctx := c.Request().Context()
	deviceID := payload.GetDeviceID(c.Request())
	if len(deviceID) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			response.MsgErr("invalidate", "empty deviceID"),
		)
	}

	if err := ctr.authUsecase.SendVerifyCode(ctx, req.Email); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.Err(err, "send verify code", "request: %+v", req),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.Msg("send verify code success"),
	)
}

func (ctr *AuthController) Login(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     int    `json:"code"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.Err(err, "invalid request", "request: %+v", req),
		)
	}

	ctx := c.Request().Context()

	deviceID := payload.GetDeviceID(c.Request())
	if len(deviceID) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			response.MsgErr("invalid header parameter", "empty deviceID"),
		)
	}

	authToken, err := ctr.authUsecase.Login(ctx, usecase.LoginParam{
		Email:    req.Email,
		Password: req.Password,
		Code:     strconv.Itoa(req.Code),
		DeviceID: deviceID,
	})
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.Err(err, "login", "request: %+v", req),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.Data(authToken, "login success"),
	)
}

func (ctr *AuthController) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()
	userID := payload.GetUserID(ctx)
	if userID == 0 {
		return c.JSON(
			http.StatusBadRequest,
			response.MsgErr("invalid request", "empty userID"),
		)
	}

	deviceID := payload.GetDeviceID(c.Request())
	if len(deviceID) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			response.MsgErr("invalid request", "empty deviceID"),
		)
	}

	authToken, err := ctr.authUsecase.RefreshToken(ctx, userID, deviceID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.Err(err, "refresh token", "userID: %d", userID),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.Data(authToken, "refresh token success"),
	)
}

func (ctr *AuthController) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	userID := payload.GetUserID(ctx)
	if userID == 0 {
		return c.JSON(
			http.StatusBadRequest,
			response.MsgErr("invalid request", "empty userID"),
		)
	}

	deviceID := payload.GetDeviceID(c.Request())
	if len(deviceID) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			response.MsgErr("invalid request", "empty deviceID"),
		)
	}

	if err := ctr.authUsecase.Logout(ctx, userID, deviceID); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.Err(err, "logout", "userID: %d", userID),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.Msg("logout success"),
	)
}
