package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").NotEmpty().Unique(),
		field.String("email").NotEmpty().Unique(),
		field.String("password").NotEmpty(),
		field.String("first_name"),
		field.String("last_name"),
		field.String("phone_number"),
		field.Bool("email_verified").Default(false),
		field.Bool("phone_verified").Default(false),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_roles", UserRole.Type),
	}
}
