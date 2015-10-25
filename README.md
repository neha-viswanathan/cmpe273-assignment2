# cmpe273-assignment2

Location and Trip Planner Service in Go using httprouter handler and Google Maps API

CRUD Operations with sample response:

1. Create Location - POST
Request: localhost:1111/locations
{
    "Name":"AT&T Park",
    "Address":"24 Willie Mays Plaza",
    "City":"San Francisco",
    "State":"CA",
    "Zip":"94107"
}
Response: 201 Created
{
  "_id": "562c525be7024724c440210f",
  "Name": "AT&T Park",
  "Address": "24 Willie Mays Plaza",
  "City": "San Francisco",
  "State": "CA",
  "Zip": "94107",
  "Coordinates": {
    "Lattitude": "37.7781747",
    "Longitude": "-122.3907248"
  }
}

2. Retrieve Location - GET
Request: localhost:1111/locations/562c52ffe7024723488f2b30

Response: 200 OK
{
  "_id": "562c52ffe7024723488f2b30",
  "Name": "AT&T Park",
  "Address": "24 Willie Mays Plaza",
  "City": "San Francisco",
  "State": "CA",
  "Zip": "94107",
  "Coordinates": {
    "Lattitude": "37.7781747",
    "Longitude": "-122.3907248"
  }
}

3. Update Location - PUT
Request: localhost:1111/locations/562c52ffe7024723488f2b30
{
    "Address":"900 North Point St #52",
    "City":"San Francisco",
    "State":"CA",
    "Zip":"94109"
}
Response: 201 Created
{
  "_id": "562c52ffe7024723488f2b30",
  "Name": "AT&T Park",
  "Address": "900 North Point St #52",
  "City": "San Francisco",
  "State": "CA",
  "Zip": "94109",
  "Coordinates": {
    "Lattitude": "37.8055762",
    "Longitude": "-122.4229471"
  }
}

4. Delete Location - DELETE
Request: localhost:1111/locations/562c52ffe7024723488f2b30

Response: 200 OK
