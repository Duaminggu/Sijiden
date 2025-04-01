package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserRole holds the schema definition for the UserRole entity.
type UserRole struct {
	ent.Schema
}

func (UserRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sijiden_user_roles"},
	}
}

// Fields of the UserRole.
func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("role_id"),
	}
}

// Edges of the UserRole.
func (UserRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("user_roles").
			Field("user_id").
			Unique().
			Required(),
		edge.From("role", Role.Type).
			Ref("user_roles").
			Field("role_id").
			Unique().
			Required(),
	}
}
