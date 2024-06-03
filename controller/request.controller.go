package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"usermanagement/model"
	"usermanagement/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
)

// Render 401
func GET401ViewCreatePage(c *fiber.Ctx) error {
	return c.Render("401", fiber.Map{"Title": "401"})
}

// Render 404
func GET404ViewCreatePage(c *fiber.Ctx) error {
	return c.Render("404", fiber.Map{"Title": "404"})
}

// Render 500
func GET500ViewCreatePage(c *fiber.Ctx) error {
	return c.Render("500", fiber.Map{"Title": "500"})
}

// Render login
func GETLoginViewCreatePage(c *fiber.Ctx) error {
	session, err := utils.Store.Get(c)
	if err == nil {
		if session.Get(utils.EMAIL) != nil {
			return flash.WithSuccess(c, fiber.Map{"Title": "dashboard"}).Redirect("dashboard")
		}
	}
	return c.Render("login", fiber.Map{"Title": "login"})
}

// Render POST login
func POSTLoginViewCreatePage(c *fiber.Ctx) error {
	fm := fiber.Map{"Title": "login"}

	email := c.FormValue("email")
	password := utils.GetEncryptionSHA512(c.FormValue("password"))

	var (
		user model.User
		err  error
	)

	if user, err = model.GetUserByEmail(email); err != nil {
		fm["message"] = "There is no user with that email"
		fm["error"] = true

		return c.Render("login", fm)
	}

	if user.IsDelete == 1 {
		fm["message"] = "Something went wrong! raise your concern with Admin."
		fm["error"] = true

		return c.Render("login", fm)
	}

	if user.IsActive == 0 {
		fm["message"] = "Your account is not Active yet."
		fm["error"] = true

		return c.Render("login", fm)
	}

	if strings.Compare(user.Password, password) != 0 {
		fm["message"] = "Incorrect password"
		fm["error"] = true

		return c.Render("login", fm)
	}

	session, err := utils.Store.Get(c)
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)
		fm["error"] = true

		return c.Render("login", fm)
	}

	session.Set(utils.ID, user.ID)
	session.Set(utils.EMAIL, user.Email)
	session.Set(utils.NAME, user.Name)
	session.Set(utils.ROLE, user.Role)

	err = session.Save()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return c.Render("login", fm)
	}

	return flash.WithSuccess(c, fm).Redirect("dashboard")
}

// Render register
func GETRegisterViewCreatePage(c *fiber.Ctx) error {
	session, err := utils.Store.Get(c)
	if err == nil {
		if session.Get(utils.EMAIL) != nil {
			return c.Render("dashboard", fiber.Map{"Title": "dashboard"})
		}
	}
	return c.Render("register", fiber.Map{"Title": "register"})
}

// Render POST register
func POSTRegisterViewCreatePage(c *fiber.Ctx) error {
	fm := fiber.Map{
		"Title": "register",
	}

	user := model.User{
		Name:     c.FormValue("fullname"),
		Email:    c.FormValue("email"),
		Password: utils.GetEncryptionSHA512(c.FormValue("password")),
	}

	_, err := model.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			err = errors.New("the email is already in use")
		}

		fmt.Println(err)
		fm["error"] = 1
		fm["message"] = "Email already taken. Please try again with other Email."
	} else {
		fm["success"] = 1
		fm["message"] = "Registration successful! wait for Admin approval."
	}

	return c.Render("register", fm)
}

// Render index
func GETIndexViewCreatePage(c *fiber.Ctx) error {
	fm := fiber.Map{"Title": "index"}

	users, err := model.GetAllUsers()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)
		fm["error"] = true
	} else {
		fm["error"] = false
		fm["users"] = users
	}
	return c.Render("index", fm)
}

// Render dashboard
func GETDashboardViewCreatePage(c *fiber.Ctx) error {
	fm := fiber.Map{
		"Title": "dashboard",
	}

	session, err := utils.Store.Get(c)
	if err != nil {
		GETLogoutPage(c)
	}

	if session.Get(utils.NAME) != nil {
		fm["Name"] = session.Get(utils.NAME)
	}

	return c.Render("dashboard", fm)
}

// Render tables
func GETProfileViewCreatePage(c *fiber.Ctx) error {
	fm := fiber.Map{"Title": "profile"}
	session, err := utils.Store.Get(c)
	if err != nil {
		return flash.WithError(c, fm).Redirect("login")
	}

	if session.Get(utils.NAME) != nil {
		fm["Name"] = session.Get(utils.NAME)
	}
	return c.Render("profile", fm)
}

// Render tables
func GETAllUsersViewCreatePage(c *fiber.Ctx) error {
	fm := fiber.Map{"Title": "users"}

	session, err := utils.Store.Get(c)
	if err != nil {
		return flash.WithError(c, fm).Redirect("login")
	}

	if session.Get(utils.ID) != nil {
		users, err := model.GetAllUsersForAdmin(fmt.Sprint(session.Get(utils.ID)))
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)
			fm["error"] = true
		} else {
			fm["error"] = false
			fm["users"] = users
		}
	}

	if session.Get(utils.NAME) != nil {
		fm["Name"] = session.Get(utils.NAME)
	}

	return c.Render("all-users", fm)
}

// Delete User
func GETDeleteUser(c *fiber.Ctx) error {
	fm := fiber.Map{}

	session, err := utils.Store.Get(c)
	if err != nil {
		return flash.WithError(c, fm).Redirect("login")
	}

	if session.Get(utils.ID) == nil {
		return flash.WithError(c, fm).Redirect("login")
	}

	idParams, _ := strconv.Atoi(c.Query("id"))
	userID := uint64(idParams)

	// fmt.Println(userID)

	_, err = model.DeleteUserByID(userID)

	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("all-users")
	}

	return flash.WithSuccess(c, fm).Redirect("all-users")
}

// Handle Logout request
func GETLogoutPage(c *fiber.Ctx) error {
	fm := fiber.Map{
		"Title": "login",
	}

	session, err := utils.Store.Get(c)
	if err != nil {
		fm["message"] = "logged out (no session)"

		return flash.WithError(c, fm).Redirect("login")
	}

	err = session.Destroy()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("login")
	}

	fm["success"] = true
	fm["message"] = "logged out successfully."

	return flash.WithSuccess(c, fm).Redirect("login")
}
