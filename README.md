# nearestPlaces

REST API server that finds nearest places using your location. It uses Elasticsearch to store and index data.

## Authentication

To get a token, send a query to http://127.0.0.1:8888/api/get_token. You'll get a JSON with the token field:

```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNjAxOTc1ODI5LCJuYW1lIjoiTmlrb2xheSJ9.FqsRe0t9YhvEC3hK1pCWumGvrJgz9k9WvhJgO8HsIa8"
}
```

To use the token, specify `Authorization: Bearer <your_token>` HTTP header. Unauthorized requests to /api/recommend endpoint will get HTTP 401 error.

## Description

Elasticsearch is a full text search engine built on top of [Lucene](https://en.wikipedia.org/wiki/Apache_Lucene). It provides an HTTP API that we will be using in this task.

When initializing, storage is populated with dataset of restaurants (taken from an Open Data portal) consists of more than 13 thousands of restaurants in the area of Moscow, Russia (you can put together another similar dataset for any other location you want, see datasets/ folder). 
Every entry has:

- ID
- Name
- Address
- Phone
- Longitude
- Latitude

Before uploading all entries into the database, an index and a mapping is created (explicitly specifying data types). Without it Elasticsearch will try to guess field types based on documents provided, and sometimes it won't recognize geopoints. 

Mapping is based on schema. `schema.json` looks like this:

```
{
  "properties": {
    "name": {
        "type":  "text"
    },
    "address": {
        "type":  "text"
    },
    "phone": {
        "type":  "text"
    },
    "location": {
      "type": "geo_point"
    }
  }
}
```

<h3>Simplest Interface</h3>

You can see restaurants added to the database using a web browser. Just enter "http://127.0.0.1:8888/?page=2" in the search box.

<h3>REST API</h3>

If you want to get a list of restaurants, you can use /api/places endpoint with the `page` query parameter.

In case 'page' param is specified with a wrong value (outside [0..last_page] or not numeric) API responds with a HTTP 400 error and similar JSON:

```
{
    "error": "Invalid 'page' value: 'foo'"
}
```

<h3>Closest Restaurants</h3>

Search for three closest restaurants. Send a GET query to /api/recommend specifying `lat` and `lon` query parameters.

`lat` and `lon` can be your current coordinates. So, for an URL http://127.0.0.1:8888/api/recommend?lat=55.674&lon=37.666 application returns JSON like this:

```
{
  "name": "Recommendation",
  "places": [
    {
      "id": 30,
      "name": "Ryba i mjaso na ugljah",
      "address": "gorod Moskva, prospekt Andropova, dom 35A",
      "phone": "(499) 612-82-69",
      "location": {
        "lat": 55.67396575768212,
        "lon": 37.66626689310591
      }
    },
    {
      "id": 3348,
      "name": "Pizzamento",
      "address": "gorod Moskva, prospekt Andropova, dom 37",
      "phone": "(499) 612-33-88",
      "location": {
        "lat": 55.673075576456,
        "lon": 37.664533747576
      }
    },
    {
      "id": 3347,
      "name": "KOFEJNJa «KAPUChINOFF»",
      "address": "gorod Moskva, prospekt Andropova, dom 37",
      "phone": "(499) 612-33-88",
      "location": {
        "lat": 55.672865251005106,
        "lon": 37.6645689561318
      }
    }
  ]
}
```
