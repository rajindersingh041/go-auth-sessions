package orderproduction

import (
	"context"
	"time"
)

type Production struct {
	OrderID               uint64    `json:"order_id"`
	ProductionID          string    `json:"production_id"`
	ProductionTimestamp   time.Time `json:"production_timestamp"`
	FulfillmentCentreID   uint16    `json:"fulfillement_center_id"`
	FulfillmentCentreName string    `json:"fulfillement_center_name"`
}

type ProductionRepositary interface {
	Create(ctx context.Context, production *Production) error
	GetbyOrderID(ctx context.Context, orderID string) (*Production, error)
	GetbyProductionID(ctx context.Context, productionID string) (*Production, error)

}


type CreateProductionRequest struct {
	OrderID					uint64		`json:"order_id"`
	ProductionID 			string 		`json:"production_id"`
	ProductionTimestamp 	time.Time 	`json:"production_timestamp"`
	FulfillmentCentreID   	uint16    	`json:"fulfillement_center_id"`
	FulfillmentCentreName 	string    	`json:"fulfillement_center_name"`
}


