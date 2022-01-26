package sqlite

import (
	"context"
	"log"

	"github.com/andrwkng/hudumaapp/model"
)

type PlanService struct {
	db *DB
}

func NewPlanService(db *DB) *PlanService {
	return &PlanService{db}
}

func (s *PlanService) CreatePlan(ctx context.Context, plan *model.Plan) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createPlan(ctx, tx, plan); err != nil {
		return err
	}
	return tx.Commit()
}

func createPlan(ctx context.Context, tx *Tx, plan *model.Plan) error {
	slugName := plan.Name
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO plan (
			name,
			description,
			code,
			period,
			interval,
			billing_cycles,
			amount,
			currency
		) VALUES (?,?,?,?,?,?,?,?)
		`,
		plan.Name,
		plan.Description,
		slugName,
		plan.Period,
		plan.Interval,
		plan.BillingCycles,
		plan.Price,
		plan.Currency,
	); err != nil {
		log.Println("Failed inserting plan into db:", err)
		return err
	}

	return nil
}
