package orderproduction

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) ProductionRepositary {
	return &PostgresRepository{db: db}
}

// ensureProductionsTable creates the productions table if it does not exist
func (r *PostgresRepository) ensureProductionsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS productions (
	    id SERIAL PRIMARY KEY,
	    order_id BIGINT NOT NULL,
	    production_id TEXT NOT NULL,
	    production_timestamp TIMESTAMP NOT NULL,
	    fulfillement_center_id INT NOT NULL,
	    fulfillement_center_name TEXT NOT NULL
	)`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, production *Production) error {
       if err := r.ensureProductionsTable(ctx); err != nil {
	       return err
       }
       query := `INSERT INTO productions (order_id, production_id, production_timestamp, fulfillement_center_id, fulfillement_center_name) VALUES ($1, $2, $3, $4, $5)`
       _, err := r.db.ExecContext(ctx, query,
	       production.OrderID,
	       production.ProductionID,
	       production.ProductionTimestamp,
	       production.FulfillmentCentreID,
	       production.FulfillmentCentreName,
       )
       return err
}

func (r *PostgresRepository) GetbyOrderID(ctx context.Context, orderID string) (*Production, error) {
	// Minimal stub, not implemented
	return nil, nil
}

func (r *PostgresRepository) GetbyProductionID(ctx context.Context, productionID string) (*Production, error) {
	// Minimal stub, not implemented
	return nil, nil
}
