// internal/handler/order.go
package handler

import (
	"context"
	"errors"
	"net/http"

	"oficina-os/internal/domain"

	"github.com/labstack/echo/v5"
)

type OrderService interface {
	Create(ctx context.Context, o domain.Order) (domain.Order, error)
	List(ctx context.Context) ([]domain.Order, error)
	Get(ctx context.Context, id string) (domain.Order, error)
	Update(ctx context.Context, id string, o domain.Order) (domain.Order, error)
	Delete(ctx context.Context, id string) error
}

type OrderHandler struct {
	svc OrderService
}

func NewOrderHandler(svc OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Register(g *echo.Group) {
	g.GET("", h.List)
	g.POST("", h.Create)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *OrderHandler) Create(c *echo.Context) error {
	var o domain.Order
	if err := c.Bind(&o); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	created, err := h.svc.Create(c.Request().Context(), o)
	if err != nil {
		return mapError(err)
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *OrderHandler) List(c *echo.Context) error {
	orders, err := h.svc.List(c.Request().Context())
	if err != nil {
		return mapError(err)
	}
	return c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) Get(c *echo.Context) error {
	id, err := pathID(c)
	if err != nil {
		return err
	}
	o, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return mapError(err)
	}
	return c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) Update(c *echo.Context) error {
	id, err := pathID(c)
	if err != nil {
		return err
	}
	var o domain.Order
	if err := c.Bind(&o); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	updated, err := h.svc.Update(c.Request().Context(), id, o)
	if err != nil {
		return mapError(err)
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *OrderHandler) Delete(c *echo.Context) error {
	id, err := pathID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return mapError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func pathID(c *echo.Context) (string, error) {
	id, err := echo.PathParam[string](c, "id")
	if err != nil || id == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "invalid order id")
	}
	return id, nil
}

// mapError turns domain errors into HTTP errors. Only this layer knows about status codes.
func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidOrder):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		return err // echo responds with 500
	}
}
