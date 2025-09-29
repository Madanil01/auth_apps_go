package main

import (
		"github.com/gofiber/fiber/v2";
		"apps_v1/database";
		"apps_v1/routes";
		"github.com/joho/godotenv"
		"github.com/gofiber/fiber/v2/middleware/cors";
		"log"
		"os"
		
	)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var ALLOW_ORIGIN string = os.Getenv("ALLOW_ORIGIN")
	database.ConnectDB()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     ALLOW_ORIGIN,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Authorization, Content-Type",
		AllowCredentials: true,
	}))
	app.Static("/uploads", "./uploads")
	// routes.WebSocketRoute(app)
	routes.Setup(app)
	routes.AuthRoute(app)

	app.Listen("0.0.0.0:8002")
	 
}

