package recommendation

import (
	"context"
	"errors"
	"time"

	"github.com/Rhymond/go-money"
)

type Recommendation struct {
	TripStart time.Time
	TripEnd   time.Time
	HotelName string
	Location  string
	TripPrice money.Money
}

type Option struct {
	HotelName     string
	Location      string
	PricePerNight money.Money
}

type AvailabilityGetter interface {
	GetAvailability(ctx context.Context, tripStart time.Time, tripEnd time.Time, location string) ([]Option, error)
}

type Service struct {
	availability AvailabilityGetter
}

func NewService(availability AvailabilityGetter) (*Service, error) {
	if availability == nil {
		return nil, errors.New("availability must not be nil")
	}

	return &Service{
		availability: availability,
	}, nil
}
