package transport

import (
	"recommendations-ddd-go/recommendations/internal/recommendation"

	"github.com/labstack/echo/v4"
)

func NewEcho(e *echo.Echo, recHandler recommendation.Handler) {
	e.GET("/recommendation", recHandler.GetRecommendation)
}
