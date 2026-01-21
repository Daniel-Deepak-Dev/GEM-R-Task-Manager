package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// REPLACE THIS with your actual module name
	"backend/internal/models"
)

func main() {
	// 1. SETUP: Connect to DB directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, assuming defaults")
	}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	defer client.Disconnect(context.Background())
	db := client.Database("backend_test") // Use a TEST database, not production!

	// Clean up previous test runs
	db.Collection("skills").Drop(context.Background())
	db.Collection("progress_items").Drop(context.Background())

	fmt.Println("--- 1. CREATING SKILLS ---")

	// 2. CREATE PARENT SKILL
	parentID := primitive.NewObjectID()
	parentSkill := models.Skill{
		ID:       parentID,
		Name:     "Full Stack Development",
		Category: "Engineering",
	}
	_, err := db.Collection("skills").InsertOne(context.Background(), parentSkill)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created Parent: %s (ID: %s)\n", parentSkill.Name, parentID.Hex())

	// 3. CREATE CHILD SKILL
	childID := primitive.NewObjectID()
	childSkill := models.Skill{
		ID:       childID,
		Name:     "Go Backend",
		ParentID: &parentID, // Link to Parent
	}
	_, err = db.Collection("skills").InsertOne(context.Background(), childSkill)
	fmt.Printf("Created Child: %s (Linked to Parent)\n", childSkill.Name)

	fmt.Println("\n--- 2. ADDING PROGRESS ITEMS (FIBONACCI) ---")

	// 4. ADD ITEMS WITH FIBONACCI WEIGHTS
	// Scenario:
	// Item A: 5
	// Item B: 3
	// Item C: 13
	// Total Weight = 21
	weights := []int{5, 3, 13}

	for i, w := range weights {
		item := models.ProgressItem{
			ID:            primitive.NewObjectID(),
			ParentSkillID: childID, // Attached to the Child Skill
			Name:          fmt.Sprintf("Task %d", i+1),
			Weightage:     w,
			Achieved:      false,
		}

		// Run the Validation Logic manually
		if err := item.ValidateWeightage(); err != nil {
			log.Fatalf("Validation Failed for weight %d: %v", w, err)
		}

		_, err := db.Collection("progress_items").InsertOne(context.Background(), item)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Added Item: %s | Weight: %d\n", item.Name, item.Weightage)
	}

	fmt.Println("\n--- 3. RUNNING CALCULATION PIPELINE ---")

	// 5. RUN THE AGGREGATION (The "Math" Part)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "parent_skill_id", Value: childID}}}},
		{{Key: "$setWindowFields", Value: bson.D{
			{Key: "partitionBy", Value: nil},
			{Key: "output", Value: bson.D{
				{Key: "total_weight", Value: bson.D{{Key: "$sum", Value: "$weightage"}}},
			}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "weight_percent", Value: bson.D{
				{Key: "$multiply", Value: bson.A{
					bson.D{{Key: "$divide", Value: bson.A{"$weightage", "$total_weight"}}},
					100,
				}},
			}},
		}}},
	}

	cursor, _ := db.Collection("progress_items").Aggregate(context.Background(), pipeline)

	var results []bson.M // Use bson.M just to see the raw output easily
	cursor.All(context.Background(), &results)

	fmt.Printf("\n%-10s | %-10s | %-10s | %-10s\n", "NAME", "WEIGHT", "TOTAL", "CALC %")
	fmt.Println("-------------------------------------------------------")

	for _, res := range results {
		name := res["name"]
		weight := res["weightage"]
		total := res["total_weight"]
		// Format percent nicely
		percent := res["weight_percent"].(float64)

		fmt.Printf("%-10s | %-10v | %-10v | %.2f%%\n", name, weight, total, percent)
	}
}
