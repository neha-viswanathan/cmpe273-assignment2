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
		Name string `json:"name"`
		Address string `json:"address"`
		City string	`json:"city"`
		State string	`json:"state"`
		Zip string	`json:"zip"`
}

type Output struct {
		Id bson.ObjectId `json:"_id" bson:"_id,omitempty"`
		Name string `json:"name"`
		Address string `json:"address"`
		City string	`json:"city" `
		State string `json:"state"`
		Zip string	`json:"zip"`

		Coordinates struct{
			Latitude string `json:"lat"`
			Longitude string `json:"lng"`
		}
	}

type GeoResponse struct {
	Location []GeoLocations
}

type GeoLocations struct {
	Address string `json:"formatted_address"`
	AddressParts []GoogleAddress `json:"address_components"`
	Geometry Geometry
	Types []string
}

type GoogleAddress struct {
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

//Reference to GeoLocator with MongoDB session
func NewGeoLocator(ms *mgo.Session) *GeoLocator {
	return &GeoLocator{ms}
}

//Accessing Google Maps API to retrieve co-ordinates
func getGeoLocation(addr string) Output{
	client := &http.Client{}
		fmt.Println(addr)
	req_URL := "http://maps.google.com/maps/api/geocode/json?address="
	req_URL += url.QueryEscape(addr)
	req_URL += "&sensor=false";
	fmt.Println("Request URL :: ", req_URL)

	request, err := http.NewRequest("GET", req_URL , nil)
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while calling Google Maps API :: ", err);
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading response :: ", err);
	}

	var geoRes GeoResponse
	err = json.Unmarshal(body, &geoRes)
	if err != nil {
		fmt.Println("Error while unmarshalling response :: ", err);
	}

	var val Output
	val.Coordinates.Latitude = strconv.FormatFloat(geoRes.Location[0].Geometry.Location.Lat,'f',7,64)
	val.Coordinates.Longitude = strconv.FormatFloat(geoRes.Location[0].Geometry.Location.Lng,'f',7,64)

	return val;
}

//Function to retrieve location
func (geo GeoLocator) GetLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	loc_id := p.ByName("location_id")
	if !bson.IsObjectIdHex(loc_id) {
        rw.WriteHeader(http.StatusNotFound)
        return
    }

    obj_id := bson.ObjectIdHex(loc_id)
	var out Output
	if err := geo.mSess.DB("cmpe273").C("Locations").FindId(obj_id).One(&out); err != nil {
        rw.WriteHeader(http.StatusNotFound)
        return
    }

	jm, _ := json.Marshal(out)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "%s", jm)
}

//Function to create new location
func (geo GeoLocator) CreateLocation(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var in Input
	var out Output

	json.NewDecoder(req.Body).Decode(&in)
	addr := in.Address + "+" + in.City + "+" + in.State + "+" + in.Zip
	//fmt.Println("address :: ", in.Address + "+" + in.City + "+" + in.State + "+" + in.Zip)
	geoLocation := getGeoLocation(addr)
  fmt.Println("Geo Co-ordinates :: ", geoLocation.Coordinates.Latitude, geoLocation.Coordinates.Longitude);

	out.Id = bson.NewObjectId()
	out.Name = in.Name
	out.Address = in.Address
	out.City= in.City
	out.State= in.State
	out.Zip = in.Zip
	out.Coordinates.Latitude = geoLocation.Coordinates.Latitude
	out.Coordinates.Longitude = geoLocation.Coordinates.Longitude

	geo.mSess.DB("cmpe273").C("Locations").Insert(out)

	jm, _ := json.Marshal(out)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "%s", jm)
}

//Function to delete location
func (geo GeoLocator) DeleteLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	loc_id := p.ByName("location_id")

	if !bson.IsObjectIdHex(loc_id) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	obj_id := bson.ObjectIdHex(loc_id)

	if err := geo.mSess.DB("cmpe273").C("Locations").RemoveId(obj_id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

//Function to update/modify location
func (geo GeoLocator) UpdateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var in Input
	var out Output

	loc_id := p.ByName("location_id")
	if !bson.IsObjectIdHex(loc_id) {
        rw.WriteHeader(http.StatusNotFound)
        return
    }
    obj_id := bson.ObjectIdHex(loc_id)

	if err := geo.mSess.DB("cmpe273").C("Locations").FindId(obj_id).One(&out); err != nil {
        rw.WriteHeader(http.StatusNotFound)
        return
    }

	json.NewDecoder(req.Body).Decode(&in)
	geoLocation := getGeoLocation(in.Address + "+" + in.City + "+" + in.State + "+" + in.Zip);
    fmt.Println("Geo Co-ordinates :: ", geoLocation.Coordinates.Latitude, geoLocation.Coordinates.Longitude);


	out.Address = in.Address
	out.City = in.City
	out.State = in.State
	out.Zip = in.Zip
	out.Coordinates.Latitude = geoLocation.Coordinates.Latitude
	out.Coordinates.Longitude = geoLocation.Coordinates.Longitude

	c := geo.mSess.DB("cmpe273").C("Locations")

	id := bson.M{"_id": obj_id}
	err := c.Update(id, out)
	if err != nil {
		panic(err)
	}
	jm, _ := json.Marshal(out)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "%s", jm)
}
