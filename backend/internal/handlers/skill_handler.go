package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// Replace with your actual module name
	"backend/internal/models"
)

type SkillHandler struct {
	skillCollection    *mongo.Collection
	metadataCollection *mongo.Collection
}

func NewSkillHandler(skillCol *mongo.Collection, metaCol *mongo.Collection) *SkillHandler {
	return &SkillHandler{
		skillCollection:    skillCol,
		metadataCollection: metaCol,
	}
}

// Helper: Check if category is valid
func (h *SkillHandler) isValidCategory(category string) bool {
	if category == "" {
		return true
	} // Allow empty category if not strict
	filter := bson.M{"type": "SKILL_CATEGORY", "value": category, "is_active": true}
	count, _ := h.metadataCollection.CountDocuments(context.Background(), filter)
	return count > 0
}

// --- CRUD OPERATIONS ---

// CreateSkill handles both Parent and Child skills
func (h *SkillHandler) CreateSkill(c *fiber.Ctx) error {
	skill := new(models.Skill)
	if err := c.BodyParser(skill); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if !h.isValidCategory(skill.Category) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Category"})
	}

	skill.ID = primitive.NewObjectID()
	skill.CreatedAt = time.Now()

	_, err := h.skillCollection.InsertOne(context.Background(), skill)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create skill"})
	}

	return c.Status(201).JSON(skill)
}

func (h *SkillHandler) GetSkills(c *fiber.Ctx) error {
	cursor, err := h.skillCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch skills"})
	}
	var skills []models.Skill
	if err = cursor.All(context.Background(), &skills); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse skills"})
	}

	// Return empty list if nil
	if skills == nil {
		skills = []models.Skill{}
	}
	return c.JSON(skills)
}

// GetSkill fetches a single skill and populates its Parent data
func (h *SkillHandler) GetSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: objID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "skills"},
			{Key: "localField", Value: "parent_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "parent_data"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$parent_data"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := h.skillCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Query failed"})
	}

	var results []models.SkillPopulated
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Parsing failed"})
	}

	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Skill not found"})
	}

	return c.JSON(results[0])
}

func (h *SkillHandler) UpdateSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	updateData := new(models.Skill)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if !h.isValidCategory(updateData.Category) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Category"})
	}

	update := bson.M{
		"$set": bson.M{
			"name":        updateData.Name,
			"category":    updateData.Category,
			"description": updateData.Description,
			"parent_id":   updateData.ParentID, // Can change parent to move skill
		},
	}

	_, err := h.skillCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}

	return c.JSON(fiber.Map{"message": "Skill updated"})
}

func (h *SkillHandler) DeleteSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	// Pro Tip: Check if this skill is a parent to others before deleting!
	// For now, we just delete.
	_, err := h.skillCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Delete failed"})
	}

	return c.JSON(fiber.Map{"message": "Skill deleted"})
}
