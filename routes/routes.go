package routes

import (
	"apps_v1/controllers"
	"apps_v1/middleware"
	"github.com/gofiber/fiber/v2"
	"os"
)
func Setup(app *fiber.App) {
	env := os.Getenv("APP_ENV") // misal APP_ENV=local / dev / prod
	prefix := ""
	if env == "local" || env == "development" {
		prefix = "/auth"
	}
	app.Get(prefix + "/", controllers.Hello)
}
func AuthRoute(app *fiber.App) {
	env := os.Getenv("APP_ENV") // misal APP_ENV=local / dev / prod
	prefix := ""
	if env == "local" || env == "development" {
		prefix = "/auth"
	}
	app.Post(prefix+"/api/register", controllers.Register)
	app.Post(prefix+"/api/login", controllers.Login)
	app.Post(prefix+"/api/logout", controllers.Logout)
	app.Get(prefix+"/api/confirm/:email", controllers.ConfirmEmail)
	app.Post(prefix+"/api/refresh-token", controllers.RefreshToken)
	app.Get(prefix+"/profile/:user_id", middleware.AuthRequired, controllers.Profile)
}