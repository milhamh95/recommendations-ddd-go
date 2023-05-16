package recommendation

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) (*Handler, error) {
	if svc == (Service{}) {
		return nil, errors.New("svc can't be empty")
	}

	return &Handler{svc: svc}, nil
}

type GetRecommendationResponse struct {
	HotelName string    `json:"hotelName"`
	TotalCost TotalCost `json:"totalCost"`
}

type TotalCost struct {
	Cost     int64  `json:"cost"`
	Currency string `json:"currenct"`
}

func (h Handler) GetRecommendation(c echo.Context) error {
	location := c.QueryParam("location")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	budget := c.QueryParam("budget")

	const expectedFormat = "2006-01-02"

	formattedStart, err := time.Parse(expectedFormat, from)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid from date")
	}

	formattedEnd, err := time.Parse(expectedFormat, to)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid to date")
	}

	b, err := strconv.ParseInt(budget, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid budget")
	}

	budgetMon := money.New(b, "USD")
	rec, err := h.svc.Get(c.Request().Context(), formattedStart, formattedEnd, location, budgetMon)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error getting recommendation")
	}

	return c.JSON(http.StatusOK, GetRecommendationResponse{
		HotelName: rec.HotelName,
		TotalCost: TotalCost{
			Cost:     rec.TripPrice.Amount(),
			Currency: "USD",
		},
	})
}
