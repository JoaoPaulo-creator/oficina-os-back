// internal/domain/order.go
package domain

import (
	"errors"
	"fmt"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order")
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

type Part struct {
	PartName  string  `json:"partName"`
	PartValue float64 `json:"partValue"`
	Quantity  uint8   `json:"quantity"`
}

type Order struct {
	ID         string   `json:"id"`
	Customer   Customer `json:"customer"`
	Vehicle    Vehicle  `json:"vehicle"`
	Parts      []Part   `json:"parts"`
	LaborValue float64  `json:"laborValue"`
}

func (o Order) Validate() error {
	if o.Customer.Name == "" {
		return fmt.Errorf("%w: customer name is required", ErrInvalidOrder)
	}
	if o.LaborValue < 0 {
		return fmt.Errorf("%w: labor value cannot be negative", ErrInvalidOrder)
	}
	return nil
}

func (o Order) Total() float64 {
	total := o.LaborValue
	for _, p := range o.Parts {
		total += p.PartValue * float64(p.Quantity)
	}
	return total
}
