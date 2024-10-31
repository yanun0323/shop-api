package restful

import (
	"main/internal/delivery/payload"
	"main/internal/delivery/response"
	"main/internal/domain/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productUsecase usecase.ProductUsecase
}

func NewProductController(productUsecase usecase.ProductUsecase) ProductController {
	return ProductController{
		productUsecase: productUsecase,
	}
}

func (ctr *ProductController) ListRecommendation(c echo.Context) error {
	ctx := c.Request().Context()
	userID := payload.GetUserID(ctx)
	if userID == 0 {
		return c.JSON(
			http.StatusBadRequest,
			response.MsgErr("invalid request", "empty userID"),
		)
	}

	pagination := payload.GetPage(c.Request())
	products, page, err := ctr.productUsecase.ListUserRecommended(ctx, userID, pagination)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.Err(err, "list product", "userID: %d", userID),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.PagedData(products, page, "list product"),
	)
}
