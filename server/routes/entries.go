package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LOVESH1808/go-react-calorie-tracker/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

var entryCollection *mongo.Collection = OpenCollection(Client, "calories")

func AddEntry(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry

	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	// what does validation error do?
	validationErr := validate.Struct(entry)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	entry.ID = primitive.NewObjectID()
	result, InsertErr := entryCollection.InsertOne(ctx, entry)
	if InsertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": InsertErr.Error()})
		fmt.Println(InsertErr)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}

func GetEntries(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func GetEntriesByIngredient(c *gin.Context) {
	ingredientId := c.Params.ByName("id")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M

	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredientId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println()
	c.JSON(http.StatusOK, entries)
}

func GetEntryById(c *gin.Context) {
	EntryId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(EntryId)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(entry)
		return
	}
	defer cancel()
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)
}

func UpdateIngredient(c *gin.Context) {
	entryId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}

	var ingredient Ingredient
	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docId},
		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.MatchedCount)
}

func UpdateEntry(c *gin.Context) {
	entryId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var entry models.Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validationErr := validate.Struct(entry)

	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	result, err := entryCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docId},
		bson.M{
			"dish":        entry.Dish,
			"fat":         entry.Fat,
			"ingredients": entry.Ingredients,
			"calories":    entry.Calories,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}

func DeleteEntry(c *gin.Context) {
	entryId := c.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
	}

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
