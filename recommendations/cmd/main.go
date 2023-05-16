package main

import (
	"log"
	"recommendations-ddd-go/recommendations/internal/recommendation"
	"recommendations-ddd-go/recommendations/internal/transport"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo/v4"
)

func main() {
	c := retryablehttp.NewClient()
	c.RetryMax = 10

	partnerAdaptor, err := recommendation.NewPartnershipAdaptor(
		c.StandardClient(),
		"http://localhost:3031",
	)
	if err != nil {
		log.Fatal("failed to create a partnerAdaptor: ", err)
	}

	svc, err := recommendation.NewService(partnerAdaptor)
	if err != nil {
		log.Fatal("failed to create a service: ", err)
	}

	handler, err := recommendation.NewHandler(*svc)
	if err != nil {
		log.Fatal("failed to create a handler: ", err)
	}

	e := echo.New()
	transport.NewEcho(e, *handler)

	err = e.Start(":1324")
	if err != nil {
		log.Fatal("server error: ", err)
	}
}
