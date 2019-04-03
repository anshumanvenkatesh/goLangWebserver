# goLangWebserver
Web Server written in Go. MongoDB for serving data. Uses gin-gonic for webserver heavylifting

### Endpoints:
- `GET /` : Lists all endpoints, input format and examples
- `GET /buildings` : API for getting buildings. Has filtering mechanism via query params
- `GET /aggregate` : API for getting aggregates like `mean`, `min`, and `max` over the data
- `POST /filtAgg` : API for getting aggregates on filtered Data. Is the combination of the previous two, but takes a `POST`.

### Endpoints details:
```
{
  "apiEndpoints":              "GET /buildings(?filterParams), GET /aggregate(?aggregationParam), POST /filtAgg",
  "filterParamsFormat":        "/buildings?<dataFeature1>=<somevalue>&<dataFeature2>=<some other value>& ...",
  "filterParamsExamples":      "/buildings?ConstYear=1922&Name=AlphaHouse",
  "aggregationParamsFormat":   "/aggregate?Field=<dataFeature>&AggBy=<Aggregation operator>",
  "aggregationParamsExamples": "/aggregate?Field=ShapeArea&AggBy=mean",
  "aggregationOperators":      "mean, min, max",
  "dataFeatures":              "Bin, ConstYear, Name, HeightRoof, FeatCode, GroundElev, ShapeArea",
  "filtAggBodyFormat": gin.H{
    "filter": gin.H{
      "<dataFeature1>": "<filter value>",
      "<dataFeature2>": "<filter value>",
    },
    "aggregate": gin.H{
      "Field": "<Data Field on which aggregation should happen>",
      "AggBy": "<mean / min / max>",
    },
  },
  "filtAggBodyExample": gin.H{
    "filter": gin.H{
      "ConstYear": 2019,
      "FeatCode":  2100,
    },
    "aggregate": gin.H{
      "Field": "GroundElev",
      "AggBy": "max>",
    },
  },
}
```

### Pre-Installation:
- Make sure MongoDB is installed
- Run the ETL script from https://github.com/anshumanvenkatesh/toposETL

### Running Steps:
- Make sure the data is available at the DB: `topos` in the Collection: `testCollection`
- Clone Repo
- Install `go` dependencies by running `go get` (Alternatively use `go build` to get a Binary instead
- run by: `go run main.go`
- The server is serving at `localhost:8080` and using a MongoDB running on default address `mongodb://localhost:27017`

### Future Plans:
- Write tests
- Insert, Delete, Update Operations

### UPDATES

- Added Postman collection that has a set of requests that can be used to test the application
