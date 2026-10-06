package main

import (
	"context"
	"net/http"
	"oficina-os/internal/handler"
	"oficina-os/internal/repository/memory"
	"oficina-os/internal/service"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Customer struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

type Vehicle struct {
	Year           uint32 `json:"year"`
	Model          string `json:"model"`
	VehicleMileage uint32 `json:"vehicleMileage"`
}

type Parts struct {
	PartName  string  `json:"partName"`
	PartValue float64 `json:"partValue"`
	Quantity  uint8   `json:"quantity"`
}

func c() string {
	return uuid.NewString()
}

type Payload struct {
	ID         string   `json:"id"`
	Customer   Customer `json:"customer"`
	Vehicle    Vehicle  `json:"vehicle"`
	Parts      []Parts  `json:"parts"`
	LaborValue float64  `json:"laborValue"`
}

var (
	db   []*Payload
	dbMu sync.RWMutex
)

func createOrder(c *echo.Context) error {
	order := new(Payload)
	if err := c.Bind(order); err != nil {
		return err
	}
	order.ID = uuid.NewString()

	dbMu.Lock()
	db = append(db, order)
	dbMu.Unlock()

	return c.JSON(http.StatusCreated, order)
}

func orderID(c *echo.Context) (string, error) {
	id, err := echo.PathParam[string](c, "id")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "invalid order")
	}
	return id, nil
}

func updateOrder(c *echo.Context) error {
	id, err := orderID(c)
	if err != nil {
		return err
	}

	order := new(Payload)
	if err := c.Bind(order); err != nil {
		return err
	}

	order.ID = id

	dbMu.Lock()
	found := false
	for i, existing := range db {
		if existing.ID == id {
			order.ID = id
			db[i] = order
			found = true
			break
			// return c.JSON(http.StatusOK, order)
		}
	}
	dbMu.Unlock()

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	return c.JSON(http.StatusOK, order)
}

func main() {
	// 	e := echo.New()
	// 	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 		AllowOrigins: []string{"*"},
	// 		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// 	}))
	// 	e.Use(middleware.RequestLogger())
	//
	// 	e.GET("/orders", func(c *echo.Context) error {
	// 		dbMu.RLock()
	// 		orders := make([]*Payload, len(db))
	// 		copy(orders, db)
	// 		dbMu.RUnlock()
	// 		return c.JSON(200, db)
	// 	})
	// 	e.POST("/orders", createOrder)
	// 	e.PUT("/orders/:id", updateOrder)
	//
	// 	sc := echo.StartConfig{Address: ":3339"}
	// 	if err := sc.Start(context.Background(), e); err != nil {
	// 		e.Logger.Error("failed to start server", "error", err)
	// 	}
	repo := memory.NewOrderRepository()
	svc := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(svc)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.RequestLogger())

	orderHandler.Register(e.Group("/orders"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{Address: ":3339"}
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
