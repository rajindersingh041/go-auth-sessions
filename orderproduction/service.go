package orderproduction

import (
	"context"
	"fmt"
	"time"
)

type ProductionService interface {
	CreateProduction(ctx context.Context, req CreateProductionRequest) (*Production, error)
	// GetProductionByOrderID(ctx context.Context, orderID string) (*Production, error)
	// GetProductionByProductionID(ctx context.Context, productionID string) (*Production, error)
}

type productionService struct {
	repo ProductionRepositary
}


func NewProductionService(repo ProductionRepositary) ProductionService  {
	return &productionService{repo:repo}
}

func (p *productionService) CreateProduction(ctx context.Context, req CreateProductionRequest) (*Production, error) {
	if req.OrderID == 0 {
		return nil, fmt.Errorf("order id cannot be empty")
	}

	//create production
	orderproduction := &Production{
		OrderID: req.OrderID,
		ProductionID: req.ProductionID,
		ProductionTimestamp: time.Now(),
		FulfillmentCentreID: req.FulfillmentCentreID,
		FulfillmentCentreName: req.FulfillmentCentreName,
	}

	if err := p.repo.Create(ctx, orderproduction); err != nil {
		return nil, err
	}

	return orderproduction, nil

}
