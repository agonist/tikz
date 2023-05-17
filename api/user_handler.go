package api

import (
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleListUsers(c *fiber.Ctx) error {
	user  := types.User{

		FirstName: "John",
		LastName: "Bob",
	}
	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("James")
}