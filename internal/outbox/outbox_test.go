package outbox

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeStore struct {
	records   []Record
	published []uuid.UUID
	released  []uuid.UUID
}

func (s *fakeStore) Claim(_ context.Context, _ string, _ time.Time, _ int32) ([]Record, error) {
	return s.records, nil
}
func (s *fakeStore) MarkPublished(_ context.Context, id uuid.UUID, _ string) error {
	s.published = append(s.published, id)
	return nil
}
func (s *fakeStore) Release(_ context.Context, id uuid.UUID, _ string, _ error) error {
	s.released = append(s.released, id)
	return nil
}

type fakePublisher struct{ err error }

func (p fakePublisher) Publish(context.Context, Record) error { return p.err }

func TestRelayMarksConfirmedRecordPublished(t *testing.T) {
	record := Record{ID: uuid.New()}
	store := &fakeStore{records: []Record{record}}
	count, err := (Relay{Store: store, Publisher: fakePublisher{}, Owner: "worker-1"}).RunOnce(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, count)
	require.Equal(t, []uuid.UUID{record.ID}, store.published)
}

func TestRelayReleasesFailedPublish(t *testing.T) {
	record := Record{ID: uuid.New()}
	store := &fakeStore{records: []Record{record}}
	_, err := (Relay{Store: store, Publisher: fakePublisher{err: errors.New("broker down")}, Owner: "worker-1"}).RunOnce(context.Background())
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{record.ID}, store.released)
}
