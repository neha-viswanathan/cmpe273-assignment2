/*
* Golang REST API - CRUD operations using Google Maps API
* Neha Viswanathan
* 010029097
*/
package main

//import statements
import (
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/julienschmidt/httprouter"
	"github.com/neha-viswanathan/cmpe273-assignment2/geoLocator"
)

//main function
func main() {

	//Instatiating a new router
	router := httprouter.New()

	//Creating a NewGeoLocator instance
	geoLoc := geoLocator.NewGeoLocator(getMongoSession())

	//Create a New location
	router.POST("/locations", geoLoc.CreateLocation)

	//Retrieve a location
	router.GET("/locations/:location_id", geoLoc.GetLocation)

	//Modify/Update a location
	router.PUT("/locations/:location_id", geoLoc.UpdateLocation)

	//Delete a location
	router.DELETE("/locations/:location_id", geoLoc.DeleteLocation)

	http.ListenAndServe("localhost:1111", router)
}

//Create and return a new session
func getMongoSession() *mgo.Session {
	//Connect to MongoDB
	sess, err := mgo.Dial("mongodb://admin:root123@ds043324.mongolab.com:43324/cmpe273")

	//panic if connection failed
	if err != nil {
		panic(err)
	}

	//to make data reading consistent across sequential quesries in same session
	sess.SetMode(mgo.Monotonic, true)

	return sess
}
