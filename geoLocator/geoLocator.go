/*
* Golang REST API - CRUD operations using Google Maps API
* Neha Viswanathan
* 010029097
*/
package geoLocator

//import statements
import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
  "strconv"
  "io/ioutil"
)

//struct definitions
type GeoLocator struct {
		mSess *mgo.Session
	}

type Input struct {
		Name string `json:"Name"`
		Address string `json:"Address"`
		City string	`json:"City"`
		State string `json:"State"`
		Zip string `json:"Zip"`
}

type Output struct {
		Id bson.ObjectId `json:"_id" bson:"_id,omitempty"`
		Name string `json:"Name"`
		Address string `json:"Address"`
		City string	`json:"City" `
		State string `json:"State"`
		Zip string	`json:"Zip"`

		Coordinates struct{
			Latitude string `json:"Lattitude"`
			Longitude string `json:"Longitude"`
		}
	}

//Google Maps Response struct -- start
type GoogleResponse struct {
	Results []GoogleResult
}

type GoogleResult struct {
	Address string `json:"formatted_address"`
	AddressParts []GoogleAddressPart `json:"address_components"`
	Geometry Geometry
	Types []string
}

type GoogleAddressPart struct {
	Name string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types []string
}

type Geometry struct {
	Bounds Bounds
	Location Point
	Type string
	Viewport Bounds
}

type Bounds struct {
	NorthEast, SouthWest Point
}

type Point struct {
	Lat float64
	Lng float64
}
//Google Maps Response struct -- end

//Reference to GeoLocator with MongoDB session
func NewGeoLocator(ms *mgo.Session) *GeoLocator {
	return &GeoLocator{ms}
}

//Accessing Google Maps API to retrieve co-ordinates
func getGeoLocation(addr string) Output{
	//create a client to get data from Google Maps
	client := &http.Client{}
	//fmt.Println(addr)
	req_URL := "http://maps.google.com/maps/api/geocode/json?address="
	req_URL += url.QueryEscape(addr)
	//fmt.Println("req_URL", req_URL)
	req_URL += "&sensor=false";
	fmt.Println("Request URL :: ", req_URL)

	request, err := http.NewRequest("GET", req_URL , nil)
	//sending a http request
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while calling Google Maps API :: ", err);
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading response :: ", err);
	}

	//parsing the json response
	var geoRes GoogleResponse
	err = json.Unmarshal(body, &geoRes)
	if err != nil {
		fmt.Println("Error while unmarshalling response :: ", err);
	}

	//retrieve latitude and longitude from GoogleResponse
	var val Output
	val.Coordinates.Latitude = strconv.FormatFloat(geoRes.Results[0].Geometry.Location.Lat,'f',7,64)
	val.Coordinates.Longitude = strconv.FormatFloat(geoRes.Results[0].Geometry.Location.Lng,'f',7,64)
	fmt.Println("Retrieved co-ordinates :: ", val.Coordinates)
	return val;
}

//Function to retrieve location
func (geo GeoLocator) GetLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	loc_id := p.ByName("location_id")
	//Checking for valid hex representation of location_id
	if !bson.IsObjectIdHex(loc_id) {
        rw.WriteHeader(http.StatusNotFound)
        return
  }

	//retrieve object ID of location id
  obj_id := bson.ObjectIdHex(loc_id)
	var out Output
	//verify if object id from "GeoLocations" collection in "cmpe273" in MongoDB, retriev it
	if err := geo.mSess.DB("cmpe273").C("GeoLocations").FindId(obj_id).One(&out); err != nil {
    rw.WriteHeader(http.StatusNotFound)
    return
  }
	fmt.Println("Location found in MongoDB")

	//encoding the output in json format
	jm, _ := json.Marshal(out)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "%s", jm)
}

//Function to create new location
func (geo GeoLocator) CreateLocation(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var in Input
	var out Output

	//decode the JSON value from input
	json.NewDecoder(req.Body).Decode(&in)
	addr := in.Address + "+" + in.City + "+" + in.State + "+" + in.Zip
	//fmt.Println("address :: ", in.Address + "+" + in.City + "+" + in.State + "+" + in.Zip)
	//retrieve co-ordinates of the address
	geoLocation := getGeoLocation(addr)
  fmt.Println("Geo Co-ordinates :: ", geoLocation.Coordinates.Latitude, geoLocation.Coordinates.Longitude);

	//create a new object id for output
	out.Id = bson.NewObjectId()

	//setting output values
	out.Name = in.Name
	out.Address = in.Address
	out.City= in.City
	out.State= in.State
	out.Zip = in.Zip
	out.Coordinates.Latitude = geoLocation.Coordinates.Latitude
	out.Coordinates.Longitude = geoLocation.Coordinates.Longitude

	//store the output in Locations collection in cmpe273 database
	geo.mSess.DB("cmpe273").C("GeoLocations").Insert(out)
	fmt.Println("Location created in MongoDB")

	//encoding the output in json format
	jm, _ := json.Marshal(out)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "%s", jm)
}

//Function to update/modify location
func (geo GeoLocator) UpdateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var in Input
	var out Output

	loc_id := p.ByName("location_id")
	//Checking for valid hex representation of location_id
	if !bson.IsObjectIdHex(loc_id) {
        rw.WriteHeader(http.StatusNotFound)
        return
  }

	//retrieve object ID of location id
	obj_id := bson.ObjectIdHex(loc_id)

	//verify if object id from "GeoLocations" collection in "cmpe273" in MongoDB, retrieve it
	if err := geo.mSess.DB("cmpe273").C("GeoLocations").FindId(obj_id).One(&out); err != nil {
        rw.WriteHeader(http.StatusNotFound)
        return
    }

	//decode the JSON value from input and get co-ordinates of the input
	json.NewDecoder(req.Body).Decode(&in)
	geoLocation := getGeoLocation(in.Address + "+" + in.City + "+" + in.State + "+" + in.Zip);
  fmt.Println("Geo Co-ordinates :: ", geoLocation.Coordinates.Latitude, geoLocation.Coordinates.Longitude);

	//setting output values
	out.Address = in.Address
	out.City = in.City
	out.State = in.State
	out.Zip = in.Zip
	out.Coordinates.Latitude = geoLocation.Coordinates.Latitude
	out.Coordinates.Longitude = geoLocation.Coordinates.Longitude

	//update the GeoLocations collection with new values
	c := geo.mSess.DB("cmpe273").C("GeoLocations")

	id := bson.M{"_id": obj_id}
	err := c.Update(id, out)
	if err != nil {
		panic(err)
	}
	fmt.Println("Location updated in MongoDB")

	//encoding the output in json format
	jm, _ := json.Marshal(out)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "%s", jm)
}

//Function to delete location
func (geo GeoLocator) DeleteLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	loc_id := p.ByName("location_id")

	//Checking for valid hex representation of location_id
	if !bson.IsObjectIdHex(loc_id) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	//retrieve object ID of location id
	obj_id := bson.ObjectIdHex(loc_id)

	//verify if object id from "GeoLocations" collection in "cmpe273" in MongoDB, then remove it from DB
	if err := geo.mSess.DB("cmpe273").C("GeoLocations").RemoveId(obj_id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("Location deleted in MongoDB")

	//response
	rw.WriteHeader(http.StatusOK)
}
