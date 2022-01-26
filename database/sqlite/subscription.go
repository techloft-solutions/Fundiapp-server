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
		INSERT INTO subscriptions (
			subscription_id,
			client_id,
			payment_id,
			plan_id,
			auto_renew,
			status,
			billing_cycles,
			next_billing_at,
			activated_at,
			cancelled_at,
			starts_at,
			expires_at,
			ends_at
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
		`,
		subscription.SubscriptionID,
		subscription.ClientID,
		subscription.PaymentMethodID,
		subscription.PlanID,
		subscription.AutoRenew,
		subscription.Status,
		subscription.BillingCycles,
		subscription.NextBilling,
		subscription.ActivatedAt,
		subscription.CancelledAt,
		subscription.StartsAt,
		subscription.ExpireBy,
		subscription.EndsAt,
	); err != nil {
		log.Println("failed inserting subscription into db:", err)
		return err
	}

	return nil
}
