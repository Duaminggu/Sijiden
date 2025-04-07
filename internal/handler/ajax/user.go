package ajax

import (
	"context"
	"net/http"
	"strconv"

	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/ent/user"
	"github.com/duaminggu/sijiden/ent/userrole"
	"github.com/duaminggu/sijiden/internal/dto"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Client *ent.Client
	Store  *session.SessionStore
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	RoleIDs     []int  `json:"role_ids"`
}

func (h *UserHandler) CountUsers(c echo.Context) error {
	count, err := h.Client.User.Query().Count(context.Background())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"count": count,
	})
}

// POST /users
func (h *UserHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	// Validasi wajib
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username, email, and password are required"})
	}
	if len(req.RoleIDs) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "At least one role must be selected"})
	}

	// Cek duplikat username
	exists, err := h.Client.User.
		Query().
		Where(user.UsernameEQ(req.Username)).
		Exist(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check username"})
	}
	if exists {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username already in use"})
	}

	// Cek duplikat email
	exists, err = h.Client.User.
		Query().
		Where(user.EmailEQ(req.Email)).
		Exist(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check email"})
	}
	if exists {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email already in use"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
	}

	// Simpan user
	user, err := h.Client.User.
		Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPassword(string(hashedPassword)).
		SetFirstName(req.FirstName).
		SetLastName(req.LastName).
		SetPhoneNumber(req.PhoneNumber).
		Save(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user", "details": err.Error()})
	}

	// Assign roles ke user
	for _, roleID := range req.RoleIDs {
		_, err := h.Client.UserRole.
			Create().
			SetUserID(user.ID).
			SetRoleID(roleID).
			Save(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to assign role", "details": err.Error()})
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfully",
		"user_id": user.ID,
	})

}

func (h *UserHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	// Default pagination
	limit := 20
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// Hitung total user (untuk meta)
	totalCount, err := h.Client.User.Query().Count(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to count users", "details": err.Error()})
	}

	// Query dengan relasi user_roles
	users, err := h.Client.User.
		Query().
		WithUserRoles(func(urq *ent.UserRoleQuery) {
			urq.WithRole()
		}).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch users", "details": err.Error()})
	}

	userResponses := dto.ToUserResponses(users, true)

	// Response JSON
	return c.JSON(http.StatusOK, echo.Map{
		"data": userResponses,
		"meta": echo.Map{
			"total":  totalCount,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GET /users/:id
func (h *UserHandler) Detail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u, err := h.Client.User.Get(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, u)
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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to hash password"})
	}

	u, err := h.Client.User.UpdateOneID(id).
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPassword(string(hashedPassword)).
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
