# goLangWebserver
Web Server written in Go. MongoDB for serving data. Uses gin-gonic for webserver heavylifting

### Endpoints:
- `GET /` : Lists all endpoints, input format and examples
- `GET /buildings` : API for getting buildings. Has filtering mechanism via query params
- `GET /aggregate` : API for getting aggregates like `mean`, `min`, and `max` over the data
- `POST /filtAgg` : API for getting aggregates on filtered Data. Is the combination of the previous two, but takes a `POST`.

### Running Steps:
- Clone Repo
- Install `go` dependencies by running `go get` (Alternatively use `go build` to get a Binary instead
- run by: `go run main.go`
- Make sure the data is available at the DB: `topos` in the Collection: `testCollection`
- The server is serving at `localhost:8080` and using a MongoDB running on default address `mongodb://localhost:27017`

