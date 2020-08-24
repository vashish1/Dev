package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DbURL=os.Getenv("DbURL")
//Createdb creates a database
func Createdb() (*mongo.Collection, *mongo.Collection, *mongo.Collection, *mongo.Client) {

	clientOptions := options.Client().ApplyURI(DbURL)

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
	filter := bson.D{primitive.E{Key: "uuid", Value: e}}
	var result Profile

	err := c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result
	}
	return result
}

//Finddb finds the required database
func Finddb(c *mongo.Collection, s string) User {
	filter := bson.D{primitive.E{Key: "uuid", Value: s}}
	var result User

	err := c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result
	}
	return result
}

func FindDevelopers(c *mongo.Collection)[]Profile{

	findOptions := options.Find()
	var result []Profile

	cur, err := c.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Profile
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

func UpdateUserPostId(c *mongo.Collection,email string,id int)bool{
	filter := bson.D{{"email", email}}
	update :=bson.M{
		"$push":bson.M{"postId":id},
	}
	updateResult, err := c.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return true
}