package api

import (
	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type OrganizationHandler struct {
	orgStore db.OrgStore
}

func NewOrgHandler(orgStore db.OrgStore) *OrganizationHandler {
	return &OrganizationHandler{
		orgStore: orgStore,
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
	insertedOrg, err := h.orgStore.Insert(org)
	if err != nil {
		return err
	}
	return c.JSON(insertedOrg)
}

func (h *OrganizationHandler) HandleList(c *fiber.Ctx) error {

	orgs, err := h.orgStore.GetAll()
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

	org, err := h.orgStore.GetByID(orgID)
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
	if err := h.orgStore.Update(orgID, update); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": orgID})
}

func (h *OrganizationHandler) HandleDelete(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if err := h.orgStore.Delete(userID); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"deleted": userID})
}