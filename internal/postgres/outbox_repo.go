package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/dbgen"
	"github.com/mbeka02/ticketing-service/internal/outbox"
)

type outboxRepo struct{ store *Store }

func NewOutboxRepository(store *Store) outbox.Store { return &outboxRepo{store: store} }

func (r *outboxRepo) Claim(ctx context.Context, owner string, leaseExpiresAt time.Time, limit int32) ([]outbox.Record, error) {
	rows, err := r.store.ClaimEmailOutbox(ctx, dbgen.ClaimEmailOutboxParams{
		Limit: limit, LeaseOwner: &owner,
		LeaseExpiresAt: pgtype.Timestamptz{Time: leaseExpiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	records := make([]outbox.Record, 0, len(rows))
	for _, row := range rows {
		records = append(records, outboxFromDB(row))
	}
	return records, nil
}

func (r *outboxRepo) MarkPublished(ctx context.Context, id uuid.UUID, owner string) error {
	count, err := r.store.MarkEmailOutboxPublished(ctx, dbgen.MarkEmailOutboxPublishedParams{ID: id, LeaseOwner: &owner})
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("outbox record %s is no longer leased by %s", id, owner)
	}
	return nil
}

func (r *outboxRepo) Release(ctx context.Context, id uuid.UUID, owner string, publishErr error) error {
	message := publishErr.Error()
	count, err := r.store.ReleaseEmailOutbox(ctx, dbgen.ReleaseEmailOutboxParams{ID: id, LeaseOwner: &owner, LastError: &message})
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("outbox record %s is no longer leased by %s", id, owner)
	}
	return nil
}

func outboxFromDB(row dbgen.EmailOutbox) outbox.Record {
	return outbox.Record{
		ID: row.ID, EventType: row.EventType, SchemaVersion: int(row.SchemaVersion),
		AggregateType: row.AggregateType, AggregateID: row.AggregateID, Recipient: row.Recipient,
		Template: row.Template, Payload: row.Payload, AttemptCount: int(row.AttemptCount), CreatedAt: row.CreatedAt,
	}
}
