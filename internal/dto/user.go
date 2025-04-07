package dto

import (
	"time"

	"github.com/duaminggu/sijiden/ent"
)

type UserResponse struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	LastLoginAt time.Time `json:"last_login_at"`
	LoginsCount int       `json:"logins_count"`
	Roles       []string  `json:"roles,omitempty"`
}

func ToUserResponse(u *ent.User, includeRoles bool) UserResponse {
	res := UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		LastLoginAt: u.LastLoginAt,
		LoginsCount: u.LoginsCount,
	}

	if includeRoles && u.Edges.UserRoles != nil {
		for _, ur := range u.Edges.UserRoles {
			if ur.Edges.Role != nil {
				res.Roles = append(res.Roles, ur.Edges.Role.Name)
			}
		}
	}

	return res
}

func ToUserResponses(users []*ent.User, includeRoles bool) []UserResponse {
	var results []UserResponse
	for _, u := range users {
		results = append(results, ToUserResponse(u, includeRoles))
	}
	return results
}
