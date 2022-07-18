package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type todo struct {
	Id   int    `bson:"_id"   json:"id"`
	Task string `bson:"task"  json:"task"`
}

type collection struct {
	CollectionName string `json:"collection_name"`
	Todos          []todo `json:"todos"`
}

func main() {
	//connecting to mongoDB cloud database
	URI := "mongodb+srv://TO_DO_APP:sivahari@cluster0.3zhpt.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	db := client.Database("testDatabase")
	morningCollection := db.Collection("morning_tasks")

	//router connection established
	router := gin.Default()

	//route to create a new task
	router.POST("/create", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		var job todo
		if err := json.Unmarshal(data, &job); err != nil {
			panic(err)
		}
		_, err := morningCollection.InsertOne(context.TODO(), bson.D{{Key: "_id", Value: job.Id}, {Key: "task", Value: job.Task}})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Success!!"})
	})

	//route to update a task
	router.POST("/update", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		var job todo
		if err := json.Unmarshal(data, &job); err != nil {
			panic(err)
		}
		filter := bson.D{{Key: "_id", Value: job.Id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "task", Value: job.Task}}}}
		_, err := morningCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Task updated!!"})
	})

	//route to delete a task
	router.POST("/delete", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		var job todo
		if err := json.Unmarshal(data, &job); err != nil {
			panic(err)
		}
		filter := bson.D{{Key: "_id", Value: job.Id}}
		_, err := morningCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted!!"})
	})

	//route to fetch a task
	router.GET("/fetch/:id", func(c *gin.Context) {
		id := c.Param("id")
		Id, _ := strconv.Atoi(id)
		filter := bson.D{{Key: "_id", Value: Id}}
		result := morningCollection.FindOne(context.TODO(), filter)
		if result.Err() != nil {
			panic(result.Err())
		}
		var job todo
		if err := result.Decode(&job); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"id": job.Id, "task": job.Task})
	})

	//route to create a new collection with different tasks
	router.POST("/createCollection", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		var coll collection
		if err := json.Unmarshal(data, &coll); err != nil {
			panic(err)
		}
		newCollection := db.Collection(coll.CollectionName)
		for _, todo := range coll.Todos {
			_, err := newCollection.InsertOne(context.TODO(), bson.D{{Key: "_id", Value: todo.Id}, {Key: "task", Value: todo.Task}})
			if err != nil {
				panic(err)
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "new collection created!!"})
	})

	//route to create a group of todos
	router.POST("/deleteTasks", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		var coll collection
		if err := json.Unmarshal(data, &coll); err != nil {
			panic(err)
		}
		newCollection := db.Collection(coll.CollectionName)
		_, err := newCollection.DeleteMany(context.TODO(), bson.D{})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully!!"})
	})
	router.Run("localhost:8080")
}
