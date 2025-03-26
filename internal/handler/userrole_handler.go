package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/duaminggu/sijiden/ent"
	"github.com/labstack/echo/v4"
)

type UserRoleHandler struct {
	Client *ent.Client
}

type CreateUserRoleRequest struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

// POST /user-roles
func (h *UserRoleHandler) Create(c echo.Context) error {
	var req CreateUserRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	userRole, err := h.Client.UserRole.
		Create().
		SetUserID(req.UserID).
		SetRoleID(req.RoleID).
		Save(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, userRole)
}

// GET /user-roles
func (h *UserRoleHandler) List(c echo.Context) error {
	userRoles, err := h.Client.UserRole.
		Query().
		WithUser().
		WithRole().
		All(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, userRoles)
}

// GET /user-roles/:id
// func (h *UserRoleHandler) Get(c echo.Context) error {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
// 	}

// 	userRole, err := h.Client.UserRole.
// 		Query().
// 		Where(ent.UserRole.IDEQ(id)).
// 		WithUser().
// 		WithRole().
// 		Only(context.Background())
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, echo.Map{"error": "UserRole not found"})
// 	}

// 	return c.JSON(http.StatusOK, userRole)
// }

// DELETE /user-roles/:id
func (h *UserRoleHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	err = h.Client.UserRole.
		DeleteOneID(id).
		Exec(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
