package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/ent/user"
	"github.com/duaminggu/sijiden/ent/userrole"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Client *ent.Client
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

// POST /users
func (h *UserHandler) Create(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	user, err := h.Client.User.
		Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPassword(req.Password).
		SetFirstName(req.FirstName).
		SetLastName(req.LastName).
		SetPhoneNumber(req.PhoneNumber).
		Save(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// GET /users
func GetUsers(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := client.User.Query().All(context.Background())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, users)
	}
}

// GET /users/:id
func (h *UserHandler) Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u, err := h.Client.User.Get(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, u)
}

// PUT /users/:id
func (h *UserHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	u, err := h.Client.User.UpdateOneID(id).
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPassword(req.Password).
		SetFirstName(req.FirstName).
		SetLastName(req.LastName).
		SetPhoneNumber(req.PhoneNumber).
		Save(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, u)
}

// DELETE /users/:id
func (h *UserHandler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Client.User.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// GET /users/:id/roles
func (h *UserHandler) GetUserRoles(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	roles, err := h.Client.UserRole.
		Query().
		Where(userrole.HasUserWith(user.IDEQ(id))).
		WithRole().
		All(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	var result []any
	for _, ur := range roles {
		result = append(result, ur.Edges.Role)
	}

	return c.JSON(http.StatusOK, result)
}
