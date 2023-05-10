package main

import (
	"fmt"
	"log"
	"main/server"
	"main/server/db"
	"main/server/services/twilio"

	"os"

	"github.com/joho/godotenv"
)

// @title Ecommerce
// @version 1.0.0
// @description API Documentation for Ecommerce-Application
// @host localhost:3000
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("auth token: ", os.Getenv("TWILIO_AUTH_TOKEN"))
	twilio.TwilioInit(os.Getenv("TWILIO_AUTH_TOKEN"))

	connection := db.InitDB()
	db.Transfer(connection)

	// defer func() {

	// 	if err := connection.DB().Close(); err != nil {
	// 		log.Print(err)
	// 	}
	// }()

	app := server.NewServer(connection)
	server.ConfigureRoutes(app)

	if err := app.Run(os.Getenv("PORT")); err != nil {
		log.Print(err)
	}
}
