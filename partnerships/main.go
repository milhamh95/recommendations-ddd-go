package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

type Res struct {
	AvailableHotels []AvailableHotel `json:"availableHotels"`
}

type AvailableHotel struct {
	Name               string `json:"name"`
	PriceInUSDPerNight int    `json:"priceInUSDPerNight"`
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 1
	max := 10

	sampleRes := Res{
		AvailableHotels: []AvailableHotel{
			{
				Name:               "some hotel",
				PriceInUSDPerNight: 300,
			},
			{
				Name:               "some other hotel",
				PriceInUSDPerNight: 30,
			},
			{
				Name:               "some third hotel",
				PriceInUSDPerNight: 90,
			},
			{
				Name:               "some fourth hotel",
				PriceInUSDPerNight: 80,
			},
		},
	}

	e := echo.New()
	e.GET("/parnterships", func(c echo.Context) error {
		ran := r.Intn(max - min + 1)
		if ran > 7 {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return c.JSON(http.StatusOK, sampleRes)
	})

	log.Println("running")

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
