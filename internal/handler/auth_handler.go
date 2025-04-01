package handler

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/ent/user"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RedirectTo string `json:redirectTo`
}

func Login(client *ent.Client, store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserLoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		// ✅ Ambil session dan validasi CSRF token
		cookie, err := c.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "no session"})
		}

		sentToken := c.Request().Header.Get("X-CSRF-Token")
		expectedToken, ok := store.GetCSRF(cookie.Value)
		if !ok || sentToken != expectedToken {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "invalid csrf token"})
		}

		if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "email dan password wajib diisi"})
		}

		ctx := context.Background()

		// Cek user berdasarkan email
		userRecord, err := client.User.
			Query().
			Where(user.EmailEQ(req.Email)).
			Only(ctx)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "email atau password salah"})
		}

		// Bandingkan password
		err = bcrypt.CompareHashAndPassword([]byte(userRecord.Password), []byte(req.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "email atau password salah"})
		}

		// 🔁 Ambil semua role user
		userRoles, err := userRecord.
			QueryUserRoles().
			QueryRole().
			All(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "gagal mengambil roles"})
		}

		roleNames := []string{}
		for _, r := range userRoles {
			roleNames = append(roleNames, r.Name)
		}

		// TODO: Set role in cookies or session or store if it can, which best

		sessionID := uuid.NewString()
		store.Set(sessionID, userRecord.ID)
		store.SetValue(sessionID, "roles", strings.Join(roleNames, ","))

		c.SetCookie(&http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // true untuk HTTPS
		})

		// 🔁 Tentukan redirect
		redirectTo := strings.TrimSpace(req.RedirectTo)

		if redirectTo == "" && len(userRoles) > 0 {
			sort.Slice(userRoles, func(i, j int) bool {
				return userRoles[i].ID < userRoles[j].ID
			})
			redirectTo = userRoles[0].RedirectUrl
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message":    "login berhasil",
			"user_id":    userRecord.ID,
			"username":   userRecord.Username,
			"redirectTo": redirectTo,
		})
	}
}

func Register(client *ent.Client, store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UserRegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		// ✅ Ambil session dan validasi CSRF token
		cookie, err := c.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "no session"})
		}

		sentToken := c.Request().Header.Get("X-CSRF-Token")
		expectedToken, ok := store.GetCSRF(cookie.Value)
		if !ok || sentToken != expectedToken {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "invalid csrf token"})
		}

		// Basic validation
		if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || len(req.Password) < 6 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "username, email harus diisi dan password minimal 6 karakter",
			})
		}

		ctx := context.Background()

		exists, err := client.User.
			Query().
			Where(user.Or(
				user.UsernameEQ(req.Username),
				user.EmailEQ(req.Email),
			)).
			Exist(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
		}
		if exists {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "username atau email sudah terdaftar",
			})
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to hash password"})
		}

		// Simpan ke DB
		user, err := client.User.
			Create().
			SetUsername(req.Username).
			SetEmail(req.Email).
			SetPassword(string(hashedPassword)). // nanti bisa di-hash
			Save(ctx)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"message": "user created",
			"user":    user,
		})

	}
}
