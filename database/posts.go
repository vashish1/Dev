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
	Dislikes int
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

//UpdateComments updates the Post info
func UpdateComments(c *mongo.Collection,id int,cmt string)bool{
	filter := bson.D{
		{"id", id},
	}
	update := bson.D{
		{
			"$push",bson.D{{"comments",cmt}},
		},
	}
	updateResult, err := c.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return true
}

func FindComment(c *mongo.Collection,email string,id int)Post{
	filter := bson.D{
		{"id", id},
		{"email",email},
	}
	var result Post
	err := c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result
	}
	return result
}

//UpdateLikes ........
func UpdateLikes(c *mongo.Collection,email string,id int)bool{
	filter := bson.D{
		{"id", id},
		{"email",email},
	}
	update:=bson.D{
		{
			"$inc",bson.D{{"likes",1}},},
	}
	updateResult, err := c.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return true
}

//UpdateDisLikes ........
func UpdateDisLikes(c *mongo.Collection,email string,id int)bool{
	filter := bson.D{
		{"id", id},
		{"email",email},
	}
	update:=bson.D{
		{
			"$inc",bson.D{{"dislikes",1}},},
	}
	updateResult, err := c.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return true

}