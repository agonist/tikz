package api

import (
	"github.com/agonist/hotel-reservation/db"
	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	s *db.Store
}

func NewEventHandler(store *db.Store) *EventHandler {
	return &EventHandler{
		s: store,
	}
}

func (h *EventHandler) HandlePost(c *fiber.Ctx) error {
	var params types.CreateEventParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}

	event, err := types.NewEventFromParams(params)
	if err != nil {
		return err
	}
	insertedEvent, err := h.s.Event.Insert(event)
	if err != nil {
		return err
	}
	return c.JSON(insertedEvent)
}

func (h *EventHandler) HandleList(c *fiber.Ctx) error {

	events, err := h.s.Event.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(events)
}

func (h *EventHandler) HandleGet(c *fiber.Ctx) error {
	eventID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	event, err := h.s.Event.GetByID(eventID)
	if err != nil {
		return err
	}

	return c.JSON(event)
}

func (h *EventHandler) HandlePut(c *fiber.Ctx) error {
	var (
		update types.UpdateEventParams
	)
	eventID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	if err := c.BodyParser(&update); err != nil {
		return err
	}
	if err := h.s.Event.Update(eventID, update); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": eventID})
}

func (h *EventHandler) HandleDelete(c *fiber.Ctx) error {
	eventID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if err := h.s.Event.Delete(eventID); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"deleted": eventID})
}
