package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Skill struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`

	// Pointer allows this to be null
	ParentID *primitive.ObjectID `json:"parent_id" bson:"parent_id"`
}

// SkillPopulated is used only when sending data to the frontend
type SkillPopulated struct {
	Skill      `bson:",inline"` // Includes all fields from Skill struct above
	ParentData *Skill           `json:"parent_data,omitempty"`
}
