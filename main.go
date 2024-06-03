package main

import (
	"log"
	"time"
	"usermanagement/controller"
	"usermanagement/model"
	"usermanagement/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func main() {

	// godotenv package
	WEB_PORT := utils.DotEnvVariable(utils.WEB_PORT)
	model.ConnectDb()

	// Create a new engine
	engine := html.New("./views", ".html")

	// Or from an embedded system
	// See github.com/gofiber/embed for examples
	// engine := html.NewFileSystem(http.Dir("./views", ".html"))

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	/* Sessions Config */
	utils.Store = session.New(session.Config{
		CookieHTTPOnly: true,
		// CookieSecure: true, for https
		Expiration: time.Hour * 1,
	})

	// Serve static files (HTML templates and stylesheets).
	app.Static("/", "./assets")

	/* Views */
	app.Get("/401", controller.GET401ViewCreatePage)

	app.Get("/404", controller.GET404ViewCreatePage)

	app.Get("/500", controller.GET500ViewCreatePage)

	app.Get("/login", controller.GETLoginViewCreatePage)

	app.Post("/login", controller.POSTLoginViewCreatePage)

	app.Get("/register", controller.GETRegisterViewCreatePage)

	app.Post("/register", controller.POSTRegisterViewCreatePage)

	app.Get("/", controller.GETIndexViewCreatePage)

	/* Views protected with session middleware */
	adminApp := app.Group("/", controller.AuthMiddleware)
	/* todoApp.Get("/list", HandleViewList)
	todoApp.Get("/create", HandleViewCreatePage)
	todoApp.Post("/create", HandleViewCreatePage)
	todoApp.Get("/edit/:id", HandleViewEditPage)
	todoApp.Post("/edit/:id", HandleViewEditPage) */

	adminApp.Get("/dashboard", controller.GETDashboardViewCreatePage)

	adminApp.Get("/profile", controller.GETProfileViewCreatePage)

	adminApp.Get("/all-users", controller.GETAllUsersViewCreatePage)

	adminApp.Get("/delete", controller.GETDeleteUser)

	adminApp.Get("/logout", controller.GETLogoutPage)

	// a custom 404 handler at router tail
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Status(404).Render("404", fiber.Map{})
		return nil
	})

	log.Fatal(app.Listen(":" + WEB_PORT))
}
