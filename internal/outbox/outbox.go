package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Record is a durable event that still needs publication to the broker.
type Record struct {
	ID            uuid.UUID
	EventType     string
	SchemaVersion int
	AggregateType string
	AggregateID   string
	Recipient     string
	Template      string
	Payload       []byte
	AttemptCount  int
	CreatedAt     time.Time
}

type Store interface {
	Claim(context.Context, string, time.Time, int32) ([]Record, error)
	MarkPublished(context.Context, uuid.UUID, string) error
	Release(context.Context, uuid.UUID, string, error) error
}

type Publisher interface {
	Publish(context.Context, Record) error
}

type Relay struct {
	Store     Store
	Publisher Publisher
	Owner     string
	LeaseFor  time.Duration
	BatchSize int32
}

// RunOnce publishes one claimed batch. A caller schedules it at a fixed poll interval.
func (r Relay) RunOnce(ctx context.Context) (int, error) {
	if r.Store == nil || r.Publisher == nil {
		return 0, fmt.Errorf("outbox store and publisher are required")
	}
	if r.Owner == "" {
		return 0, fmt.Errorf("outbox relay owner is required")
	}
	leaseFor := r.LeaseFor
	if leaseFor <= 0 {
		leaseFor = time.Minute
	}
	batchSize := r.BatchSize
	if batchSize <= 0 {
		batchSize = 25
	}
	records, err := r.Store.Claim(ctx, r.Owner, time.Now().Add(leaseFor), batchSize)
	if err != nil {
		return 0, fmt.Errorf("claim outbox records: %w", err)
	}
	for _, record := range records {
		if err := r.Publisher.Publish(ctx, record); err != nil {
			if releaseErr := r.Store.Release(ctx, record.ID, r.Owner, err); releaseErr != nil {
				return 0, fmt.Errorf("publish %s: %w (release: %v)", record.ID, err, releaseErr)
			}
			continue
		}
		if err := r.Store.MarkPublished(ctx, record.ID, r.Owner); err != nil {
			return 0, fmt.Errorf("mark outbox record %s published: %w", record.ID, err)
		}
	}
	return len(records), nil
}
