package author

import (
	"github.com/gofiber/fiber/v3"

	"github.com/HemanthGowda/library-management-go/internal/auth"
	"github.com/HemanthGowda/library-management-go/internal/domain"
	"github.com/HemanthGowda/library-management-go/internal/httpx"
	"github.com/HemanthGowda/library-management-go/internal/user"

	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c fiber.Ctx) error {
	principal, ok := auth.PrincipalFrom(c.Context())
	if !ok {
		return domain.Unauthorized("Authentication is required")
	}
	if !principal.HasAny(user.RoleAdmin, user.RoleLibrarian) {
		return domain.Forbidden()
	}

	var request CreateAuthorRequest
	if err := httpx.Decode(c, &request); err != nil {
		return err
	}

	response, err := h.service.Add(c.Context(), request)
	if err != nil {
		return err
	}
	return httpx.WriteJSON(c, fiber.StatusCreated, response)
}

func (h *Handler) View(c fiber.Ctx) error {
	principal, ok := auth.PrincipalFrom(c.Context())

	if !ok {
		return domain.Unauthorized("Authentication is required")
	}

	if !principal.HasAny(user.RoleAdmin, user.RoleLibrarian, user.RoleMember) {
		return domain.Forbidden()
	}

	id := c.Params("id")
	response, err := h.service.View(c.Context(), id)

	if err != nil {
		return err
	}

	return httpx.WriteJSON(c, fiber.StatusOK, response)
}

func (h *Handler) Search(c fiber.Ctx) error {
	principal, ok := auth.PrincipalFrom(c.Context())
	if !ok {
		return domain.Unauthorized("Authentication is required")
	}

	if !principal.HasAny(user.RoleAdmin, user.RoleLibrarian, user.RoleMember) {
		return domain.Forbidden()
	}

	page := 1
	limit := 20

	name := c.Query("name")

	if value := c.Query("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return domain.BadRequest("Invalid page number")
		}
		page = parsed
	}

	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return domain.BadRequest("Invalid limit")
		}
		limit = parsed
	}

	response, err := h.service.Search(c.Context(), name, page, limit)
	if err != nil {
		return err
	}

	return httpx.WriteJSON(c, fiber.StatusOK, response)
}
