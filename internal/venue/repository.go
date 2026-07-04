package venue

import "context"

// Repository defines the data access contract for the venue domain.
type Repository interface {
	Create(ctx context.Context, req CreateVenueRequest) (*Venue, error)
	GetByID(ctx context.Context, id int32) (*Venue, error)
	List(ctx context.Context) ([]Venue, error)
}
