package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Createdb creates a database
func Createdb() (*mongo.Collection, *mongo.Collection, *mongo.Collection, *mongo.Client) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	usercollection := client.Database("Dev").Collection("User")
	profilecollection :=client.Database("Dev").Collection("Profile")
	Post:=client.Database("Dev").Collection("Post")
	return usercollection,profilecollection,Post, client
}

//Insertintouserdb inserts the data into the database
func Insertintouserdb(usercollection *mongo.Collection, u User) bool {

	fmt.Println(u.Name)
	insertResult, err := usercollection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return true
}



//Insertprofile inserts the data into the database
func Insertprofile(usercollection *mongo.Collection, p Profile) bool {

	fmt.Println(p.Email)
	insertResult, err := usercollection.InsertOne(context.TODO(), p)
	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return true
}
//InsertPost inserts the Post data into the database
func InsertPost(collection *mongo.Collection,name,email,text string) bool {

	post:=Post{
		UserName:name,
		Email:email,
		Text:text,
		Comments:[]string{},
	}
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return true
}
//Findfromuserdb finds the required data
func Findfromuserdb(usercollection *mongo.Collection, st string, p string) bool {
	filter := bson.D{primitive.E{Key: "email", Value: st}}
	var result User

	err := usercollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if result.PasswordHash != SHA256ofstring(p) {
		return false
	}
	return true
}

//Findprofile .............
func Findprofile(c *mongo.Collection,e string) Profile{
	filter := bson.D{primitive.E{Key: "email", Value: e}}
	var result Profile

	err := c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result
	}
	return result
}

//Finddb finds the required database
func Finddb(c *mongo.Collection, s string) User {
	filter := bson.D{primitive.E{Key: "email", Value: s}}
	var result User

	err := c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result
	}
	return result
}

func FindPost(c *mongo.Collection)[]Post {
	findOptions := options.Find()
	var result []Post

	cur, err := c.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", result)
	return result
}