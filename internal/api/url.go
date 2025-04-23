package api

import (
	"context"
	"net/http"

	"github.com/bobbybaiOuO/BShortUrl/internal/model"
	"github.com/labstack/echo/v4"
)

// URLService .
type URLService interface {
	CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error)

	GetURL(ctx context.Context, shortCode string) (string, error)
}

// URLHandler .
type URLHandler struct {
	urlService URLService
}

// NewURLHandler .
func NewURLHandler(urlService URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}


// CreateURL .
func (h *URLHandler) CreateURL(c echo.Context) error {
	// c.Request().Header.Add("Content-Type", "application/json")
	// 提取数据
	var req model.CreateURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// 验证数据格式
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// 调用业务函数
	resp, err := h.urlService.CreateURL(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusCreated, resp)
}


// RedirectURL .
func (h *URLHandler) RedirectURL(c echo.Context) error {
	// c.Request().Header.Add("Content-Type", "application/json")
	// 取出code
	shortCode := c.Param("code")

	// 根据shortcode查询url
	OriginalURL, err := h.urlService.GetURL(c.Request().Context(), shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusPermanentRedirect, OriginalURL)
}