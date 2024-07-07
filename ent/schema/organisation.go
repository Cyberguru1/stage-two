package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Organisation holds the schema definition for the Organisation entity.
type Organisation struct {
	ent.Schema
}

// Fields of the Organisation.
func (Organisation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("orgid", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("name").NotEmpty().Unique(),
		field.String("description").Optional(),
	}
}

// Edges of the Organisation.
func (Organisation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
