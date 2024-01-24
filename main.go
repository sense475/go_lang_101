package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task struct represents a simple task model.
type Task struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Completed bool      `json:"completed" bson:"completed"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

var collection *mongo.Collection

func init() {
	// Set up the MongoDB connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://username:password@localhost:27019/"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("tasksdb").Collection("tasks")
}

func main() {
	r := gin.Default()

	// Routes
	r.POST("/tasks", createTask)
	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	// Run the server
	port := 8080
	fmt.Printf("Server running on :%d\n", port)
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = "" // MongoDB will generate an ID
	task.CreatedAt = time.Now()

	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func getTasks(c *gin.Context) {
	cursor, err := collection.Find(context.Background(), bson.M{})
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var tasks []Task
	if err := cursor.All(context.Background(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func getTask(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ObjectID"})
		return
	}

	var task Task
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": task,
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
