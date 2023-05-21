package api

import (
	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type OrganizationHandler struct {
	s *db.Store
}

func NewOrgHandler(store *db.Store) *OrganizationHandler {
	return &OrganizationHandler{
		s: store,
	}
}

func (h *OrganizationHandler) HandlePost(c *fiber.Ctx) error {
	var params types.CreateOrgParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}

	org, err := types.NewOrgFromParams(params)
	if err != nil {
		return err
	}
	insertedOrg, err := h.s.Org.Insert(org)
	if err != nil {
		return err
	}
	return c.JSON(insertedOrg)
}

func (h *OrganizationHandler) HandleList(c *fiber.Ctx) error {

	orgs, err := h.s.Org.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(orgs)
}

func (h *OrganizationHandler) HandleGet(c *fiber.Ctx) error {
	orgID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	org, err := h.s.Org.GetByID(orgID)
	if err != nil {
		return err
	}

	return c.JSON(org)
}

func (h *OrganizationHandler) HandlePut(c *fiber.Ctx) error {
	var (
		update types.UpdateOrgParams
	)
	orgID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	if err := c.BodyParser(&update); err != nil {
		return err
	}
	if err := h.s.Org.Update(orgID, update); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": orgID})
}

func (h *OrganizationHandler) HandleDelete(c *fiber.Ctx) error {
	orgID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if err := h.s.Org.Delete(orgID); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"deleted": orgID})
}
