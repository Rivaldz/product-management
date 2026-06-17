package main

import (
	"log"

	"technical_test/internal/app"
)

// @title Product & Service Management API
// @version 1.0
// @description Backend API for managing products and services.

// @BasePath /
// @schemes http https
func main() {
	log.Println("Starting application...")
	app.Run()
}
