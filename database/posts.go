package database

import (
	"fmt"
	"context"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Post ...........
type Post struct{
	Id int
	UserName string
	Email string
	Text string
	Comments []string
	Likes int 
}
//PostId ..........
func PostId()int{
	rand.Seed(time.Now().UnixNano())
    min := 0
    max := 10000
    return(rand.Intn(max - min + 1) + min)
}

//InsertPost inserts the Post data into the database
func InsertPost(collection *mongo.Collection,data Post) bool {

	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return true
}


func FindPost(c *mongo.Collection)[]Post {
	findOptions := options.Find()
	var result []Post

	cur, err := c.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		fmt.Println(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}

		result = append(result, elem)
	}
	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", result)
	return result
}
