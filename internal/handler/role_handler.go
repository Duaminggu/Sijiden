package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/duaminggu/sijiden/ent"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	Client *ent.Client
}

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// POST /roles
func (h *RoleHandler) Create(c echo.Context) error {
	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	r, err := h.Client.Role.
		Create().
		SetName(req.Name).
		SetDescription(req.Description).
		Save(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, r)
}

// GET /roles
func (h *RoleHandler) List(c echo.Context) error {
	roles, err := h.Client.Role.Query().All(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, roles)
}

// GET /roles/:id
func (h *RoleHandler) Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	r, err := h.Client.Role.Get(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Role not found"})
	}
	return c.JSON(http.StatusOK, r)
}

// PUT /roles/:id
func (h *RoleHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	r, err := h.Client.Role.
		UpdateOneID(id).
		SetName(req.Name).
		SetDescription(req.Description).
		Save(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, r)
}

// DELETE /roles/:id
func (h *RoleHandler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Client.Role.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
