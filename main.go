package main

import (
	"fmt"
	"log"
	"main/server"
	"main/server/db"
	"main/server/services/user"

	"os"

	"github.com/joho/godotenv"
)

// @title Gin Demo App
// @version 1.0
// @description This is a demo version of Gin app.
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("auth token: ", os.Getenv("TWILIO_AUTH_TOKEN"))
	user.TwilioInit(os.Getenv("TWILIO_AUTH_TOKEN"))

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
