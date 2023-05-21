package api

import (
	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	s *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		s: store,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
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
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleListUsers(c *fiber.Ctx) error {

	users, err := h.s.User.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	user, err := h.s.User.GetByID(userID)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		update types.UpdateUserParams
	)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	if err := c.BodyParser(&update); err != nil {
		return err
	}
	if err := h.s.User.Update(userID, update); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": userID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if err := h.s.User.Delete(userID); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"deleted": userID})
}
