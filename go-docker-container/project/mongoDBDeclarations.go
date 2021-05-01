package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client //Our Mongo Client we work with
var theContext context.Context
var mongoURI string //Connection string loaded

func connectDB() *mongo.Client {
	//Setup Mongo connection to Atlas Cluster
	theClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Errored getting mongo client: %v\n", err)
		log.Fatal(err)
	}
	theContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = theClient.Connect(theContext)
	if err != nil {
		fmt.Printf("Errored getting mongo client context: %v\n", err)
		log.Fatal(err)
	}
	//Double check to see if we've connected to the database
	err = theClient.Ping(theContext, readpref.Primary())
	if err != nil {
		fmt.Printf("Errored pinging MongoDB: %v\n", err)
		log.Fatal(err)
	}
	//List all available databases
	/*
		databases, err := theClient.ListDatabaseNames(theContext, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(databases)
	*/

	return theClient
}

//Our Test Document to see how many times we clicked a button
type TheClicker struct {
	SpecialID int    `json:"SpecialID"`
	ClickNums int    `json:"ClickNums"`
	Name      string `json:"Name"`
}

//Gets how many times this button has been clicked
func howManyTimesClicked() TheClicker {
	var theClickerReturned TheClicker //Initialize User to be returned after Mongo query
	returnedErr := ""                 //Declare Error to be returned

	//Query for the User, given the userID for the User
	ic_collection := mongoClient.Database("thegame").Collection("timesclicked") //Here's our collection
	theFilter := bson.M{
		"name": bson.M{
			"$eq": "Name", // check if bool field has value of 'false'
		},
	}
	findOptions := options.FindOne()
	findUser := ic_collection.FindOne(theContext, theFilter, findOptions)
	if findUser.Err() != nil {
		if strings.Contains(findUser.Err().Error(), "no documents in result") {
			stringNum := strconv.Itoa(1111) //int to string conversion
			returnedErr = "For " + stringNum + ", no User was returned: " + findUser.Err().Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
		} else {
			stringNum := strconv.Itoa(1111) //int to string conversion
			returnedErr = "For " + stringNum + ", there was a Mongo Error: " + findUser.Err().Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
		}
	} else {
		err := findUser.Decode(&theClickerReturned)
		if err != nil {
			stringNum := strconv.Itoa(1111) //int to string conversion
			returnedErr = "For " + stringNum +
				", there was an error decoding document from Mongo: " + err.Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
		} else {
			stringNum := strconv.Itoa(1111) //int to string conversion
			returnedErr = "For " + stringNum +
				", Clicker should be successfully decoded."
			fmt.Println(returnedErr)
			logWriter(returnedErr)
		}
	}
	return theClickerReturned
}

func addClick(w http.ResponseWriter, r *http.Request) {
	//Initialize struct for taking messages
	type ButtonAdd struct {
		AddAmount int `json:"AddAmount"`
	}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var postedButtonAdd ButtonAdd
	json.Unmarshal(bs, &postedButtonAdd)

	fmt.Printf("Here is our amount already: %v\n", postedButtonAdd.AddAmount)

	//Declare return data and inform Ajax
	type ReturnData struct {
		SuccessMsg  string `json:"SuccessMsg"`
		SuccessBool bool   `json:"SuccessBool"`
		SuccessInt  int    `json:"SuccessInt"`
	}
	theReturnData := ReturnData{
		SuccessMsg:  "You updated the clicker",
		SuccessBool: true,
		SuccessInt:  0,
	}

	//Add the clicks
	theCurrentClicker := TheClicker{
		SpecialID: 1111,
		ClickNums: postedButtonAdd.AddAmount,
	}
	fmt.Printf("Here is our clicker: %v\n", theCurrentClicker.ClickNums)

	//Update this Document
	message_collection := mongoClient.Database("thegame").Collection("timesclicked") //Here's our collection
	theFilter := bson.M{
		"specialid": bson.M{
			"$eq": theCurrentClicker.SpecialID, // check if test value is present for reply Message
		},
	}
	updatedDocument := bson.M{
		"$set": bson.M{
			"specialid": theCurrentClicker.SpecialID,
			"clicknums": theCurrentClicker.ClickNums,
		},
	}

	theResults, err2 := message_collection.UpdateOne(theContext, theFilter, updatedDocument)
	if err2 != nil {
		fmt.Printf("We got an error updating this document: %v\n", err2.Error())
		theReturnData.SuccessBool = false
		theReturnData.SuccessInt = 1
		theReturnData.SuccessMsg = "Errored when updating this document: " + err.Error()
	} else {
		theMessage := "Message updated"
		logWriter(theMessage)
		theReturnData.SuccessBool = true
		theReturnData.SuccessInt = 0
		theReturnData.SuccessMsg = "Good update"
		fmt.Printf("DEBUG: Successful update\n")
		fmt.Printf("DEBUG: hERE are the results: %v, %v, %v\n", theResults.MatchedCount, theResults.ModifiedCount, theResults.UpsertedCount)
	}

	//Return JSON
	dataJSON, err := json.Marshal(theReturnData)
	if err != nil {
		fmt.Println("There's an error marshalling this data")
	}
	fmt.Fprintf(w, string(dataJSON))
}

func testInsertButtonClick() {
	theClicker := TheClicker{
		SpecialID: 1111,
		ClickNums: 0,
		Name:      "Name",
	}

	returnedErr := "" //Declare Error to be returned

	//Query for the User, given the userID for the User
	ic_collection := mongoClient.Database("thegame").Collection("timesclicked") //Here's our collection
	theFilter := bson.M{
		"name": bson.M{
			"$eq": "Name", // check if bool field has value of 'false'
		},
	}
	findOptions := options.FindOne()
	findUser := ic_collection.FindOne(theContext, theFilter, findOptions)
	if findUser.Err() != nil {
		if strings.Contains(findUser.Err().Error(), "no documents in result") {
			stringNum := strconv.Itoa(1111) //int to string conversion
			returnedErr = "For " + stringNum + ", no User was returned: " + findUser.Err().Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
			//Create this test document
			collectedStuff := []interface{}{theClicker}
			//Insert Our Data
			_, err := ic_collection.InsertMany(context.TODO(), collectedStuff)
			if err != nil {
				fmt.Printf("Error inserting results: \n%v\n", err)
				log.Fatal(err)
			} else {
				message := "inserted multiple documents"
				fmt.Printf("Inserted the numberClicker isntead.\n")
				logWriter(message)
			}
		} else {
			stringNum := strconv.Itoa(1111) //int to string conversion
			returnedErr = "For " + stringNum + ", there was a Mongo Error: " + findUser.Err().Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
		}
	} else {
		fmt.Printf("There a clicker document entered, hooray.")
	}
}
