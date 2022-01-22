package sqlite

import (
	"context"
	"log"

	"github.com/andrwkng/hudumaapp/model"
)

type SubscriptionService struct {
	db *DB
}

func NewSubscriptionService(db *DB) *SubscriptionService {
	return &SubscriptionService{db}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscription *model.Subscription) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createSubscription(ctx, tx, subscription); err != nil {
		return err
	}
	return tx.Commit()
}

func createSubscription(ctx context.Context, tx *Tx, subscription *model.Subscription) error {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO subscription (
			plan_id,
			status,
			billing_cycle,
			start_at,
			expire_by,
		) VALUES (?,?,?,?)
		`,
		subscription.PlanID,
		subscription.Status,
		subscription.StartAt,
		true,
	); err != nil {
		log.Println("failed inserting subscription into db:", err)
		return err
	}

	return nil
}
