package recommendation

import (
	"context"
	"errors"
	"fmt"
	"math"
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

func (svc *Service) Get(ctx context.Context, tripStart time.Time, tripEnd time.Time, location string, budget *money.Money) (*Recommendation, error) {
	if tripStart.IsZero() {
		return nil, errors.New("trip start can't be empty")
	}

	if tripEnd.IsZero() {
		return nil, errors.New("trip end can't be empty")
	}

	if location == "" {
		return nil, errors.New("location can't be empty")
	}

	opts, err := svc.availability.GetAvailability(ctx, tripStart, tripEnd, location)
	if err != nil {
		return nil, fmt.Errorf("error getting availability: %w", err)
	}

	tripDuration := math.Round(tripEnd.Sub(tripStart).Hours() / 24)
	lowestPrice := money.NewFromFloat(999999999, "USD")

	var cheapestTrip *Option
	for _, option := range opts {
		price := option.PricePerNight.Multiply(int64(tripDuration))
		ok, _ := price.GreaterThan(budget)
		if ok {
			continue
		}

		ok, _ = price.LessThan(lowestPrice)
		if ok {
			lowestPrice = price
			cheapestTrip = &option
		}
	}

	if cheapestTrip == nil {
		return nil, errors.New("no trips within budget")
	}

	return &Recommendation{
		TripStart: tripStart,
		TripEnd:   tripEnd,
		HotelName: cheapestTrip.HotelName,
		Location:  cheapestTrip.Location,
		TripPrice: *lowestPrice,
	}, nil
}
