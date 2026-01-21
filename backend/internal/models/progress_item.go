package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressItem struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ParentSkillID primitive.ObjectID `json:"parent_skill_id" bson:"parent_skill_id"`
	Name          string             `json:"name" bson:"name"`
	Achieved      bool               `json:"achieved" bson:"achieved"`

	// We store the raw weightage (1, 3, 5, 8, 13, 21) only
	Weightage int `json:"weightage" bson:"weightage"`

	Comments string `json:"comments" bson:"comments"`

	// Computed Field (We add this tag so it shows in JSON, but we don't save it to DB)
	WeightPercent float64 `json:"weight_percent,omitempty" bson:"-"`
}

// ValidateWeightage checks if the value is in your allowed Fibonacci list
func (p *ProgressItem) ValidateWeightage() error {
	allowed := map[int]bool{1: true, 3: true, 5: true, 8: true, 13: true, 21: true}
	if !allowed[p.Weightage] {
		return errors.New("weightage must be one of: 1, 3, 5, 8, 13, 21")
	}
	return nil
}
