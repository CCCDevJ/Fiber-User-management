package controller

import (
	"usermanagement/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
)

// Authentication Middleware
func AuthMiddleware(c *fiber.Ctx) error {
	fm := fiber.Map{}

	session, err := utils.Store.Get(c)
	if err != nil {
		fm["error"] = true
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("login")
	}

	if session.Get(utils.EMAIL) == nil {
		fm["error"] = true
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("login")
	}

	userId := session.Get(utils.ID)
	if userId == nil {
		fm["error"] = true
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("login")
	}

	// user, err := models.GetUserById(fmt.Sprint(userId.(uint64)))
	// if err != nil {
	// 	fm["message"] = "You are not authorized"

	// 	return flash.WithError(c, fm).Redirect("/login")
	// }

	return c.Next()
}
