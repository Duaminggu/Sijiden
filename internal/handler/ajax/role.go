package ajax

import (
	"context"
	"net/http"
	"strconv"

	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/ent/role"
	"github.com/duaminggu/sijiden/ent/userrole"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	Client *ent.Client
	Store  *session.SessionStore
}

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RedirectUrl string `json:"redirectUrl"`
}

// POST /roles
func (h *RoleHandler) Create(c echo.Context) error {
	csrfToken := c.Request().Header.Get("X-CSRF-Token")
	sessionID, err := c.Cookie("session_id")
	if err != nil || !h.Store.ValidateCSRF(sessionID.Value, csrfToken) {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Invalid CSRF token"})
	}

	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	r, err := h.Client.Role.
		Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetRedirectUrl(req.RedirectUrl).
		Save(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, r)
}

// GET /roles
func (h *RoleHandler) List(c echo.Context) error {
	roles, err := h.Client.Role.
		Query().
		WithUserRoles(). // include relasi user_roles
		All(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Siapkan hasil dengan jumlah user
	result := make([]map[string]interface{}, 0, len(roles))
	for _, r := range roles {
		result = append(result, map[string]interface{}{
			"id":          r.ID,
			"name":        r.Name,
			"description": r.Description,
			"redirectUrl": r.RedirectUrl,
			"userCount":   len(r.Edges.UserRoles), // hitung jumlah user_roles
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *RoleHandler) Detail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	role, err := h.Client.Role.
		Query().
		Where(role.IDEQ(id)).
		WithUserRoles(func(q *ent.UserRoleQuery) {
			q.WithUser() // include user di relasi
		}).
		Only(context.Background())
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Role not found"})
	}

	// Format respons
	users := []map[string]interface{}{}
	for _, ur := range role.Edges.UserRoles {
		if ur.Edges.User != nil {
			users = append(users, map[string]interface{}{
				"username": ur.Edges.User.Username,
				"email":    ur.Edges.User.Email,
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":          role.ID,
		"name":        role.Name,
		"description": role.Description,
		"redirectUrl": role.RedirectUrl,
		"created_at":  role.CreatedAt,
		"updated_at":  role.UpdatedAt,
		"users":       users,
	})
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
	// --- Validasi CSRF Token ---
	csrfToken := c.Request().Header.Get("X-CSRF-Token")
	sessionCookie, err := c.Cookie("session_id")
	if err != nil || !h.Store.ValidateCSRF(sessionCookie.Value, csrfToken) {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Invalid CSRF token"})
	}

	// --- Ambil ID dan validasi ---
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	// --- Cek apakah role adalah "admin" ---
	role, err := h.Client.Role.Get(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Role not found"})
	}
	if role.Name == "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Cannot update 'admin' role"})
	}

	// --- Binding Data ---
	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// --- Update ---
	r, err := h.Client.Role.
		UpdateOneID(id).
		SetName(req.Name).
		SetDescription(req.Description).
		SetRedirectUrl(req.RedirectUrl).
		Save(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, r)
}

// DELETE /roles/:id
func (h *RoleHandler) Delete(c echo.Context) error {
	// --- Validasi CSRF Token ---
	csrfToken := c.Request().Header.Get("X-CSRF-Token")
	sessionCookie, err := c.Cookie("session_id")
	if err != nil || !h.Store.ValidateCSRF(sessionCookie.Value, csrfToken) {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Invalid CSRF token"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	// --- Cek apakah role adalah "admin" ---
	role, err := h.Client.Role.Get(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Role not found"})
	}
	if role.Name == "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Cannot delete 'admin' role"})
	}

	// --- Cek apakah role sedang dipakai oleh user ---
	count, err := h.Client.UserRole.
		Query().
		Where(userrole.RoleID(id)).
		Count(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check user usage"})
	}
	if count > 0 {
		return c.JSON(http.StatusConflict, echo.Map{
			"error":   "Role is still used by users",
			"details": count,
		})
	}

	// --- Hapus role jika aman ---
	err = h.Client.Role.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
