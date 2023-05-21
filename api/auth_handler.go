package api

import (
	"fmt"

	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/middleware"
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	s *db.Store
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		s: store,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.s.User.GetByEmail(authParams.Email)
	if err != nil {
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return fmt.Errorf("invalid credentials")
	}

	jwt, err := middleware.GenerateJWT(user)

	return c.JSON(fiber.Map{
		"user": user,
		"jwt":  jwt,
	})
}

func (h *AuthHandler) HandleRegister(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.s.User.Insert(user)
	if err != nil {
		return err
	}

	jwt, err := middleware.GenerateJWT(insertedUser)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"user": user,
		"jwt":  jwt,
	})
}
