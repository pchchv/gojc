package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// handlePing checks that the server is up and running.
func handlePing(c *echo.Context) error {
	message := "JSON collector API. Version 0.0.1"
	return c.String(http.StatusOK, message)
}

// The declaration of all routes comes from it.
func routes(e *echo.Echo) {
	e.GET("/", handlePing)
	e.GET("/ping", handlePing)
}

func server() {
	e := echo.New()
	routes(e)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1000)))
	log.Fatal(e.Start(getEnvValue("HOST") + ":" + getEnvValue("PORT")).Error())
}
