package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

//Here is our waitgroup
var wg sync.WaitGroup

//Here are our environment examples
var TEST_VAR_ENV42069 string
var MONGO_URI string

/* TEMPLATE DEFINITION BEGINNING */
var template1 *template.Template

//Define function maps
var funcMap = template.FuncMap{
	"upperCase": strings.ToUpper, //upperCase is a key we can call inside of the template html file
}

//Parse our templates
func init() {
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
	setEnvVariables()
}

//Set our environment variables
func setEnvVariables() {
	//Look up test URI
	_, ok := os.LookupEnv("TEST_VAR_ENV42069")
	if !ok {
		message := "This env variable is not present: " + "TEST_VAR_ENV42069"
		fmt.Println(message)
		logWriter(message)
	} else {
		TEST_VAR_ENV42069 = os.Getenv("TEST_VAR_ENV42069")
		fmt.Printf("DEBUG: Environment test host is: %v\n", TEST_VAR_ENV42069)
	}
	//Look up MongoString URI
	_, ok2 := os.LookupEnv("MONGO_URI")
	if !ok2 {
		message := "This env variable is not present: " + "MONGO_URI"
		fmt.Println(message)
		logWriter(message)
	} else {
		MONGO_URI = os.Getenv("MONGO_URI")
		fmt.Printf("DEBUG: Environment Mongo host is: %v\n", MONGO_URI)
	}
}

//Writes to the log; called from most anywhere in this program!
func logWriter(logMessage string) {
	//Logging info

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "logging.txt")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
}

// Handle Errors passing templates
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	//Write to logger that we are handling requests
	debugMessage := "\n\nDEBUG: We are now handling requests"
	//fmt.Println(debugMessage)
	logWriter(debugMessage)
	fmt.Println(debugMessage)
	//Favicon and page spots
	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", index)
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//Handle Requests
	handleRequests()
}
