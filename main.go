package main

import (
	"log"
	// "command"
	"context"
	"fmt"
	"geoServer/routes"
	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type Building struct {
	// _id
	// the_geom   string
	Bin       int32  `bson:"BIN"`
	ConstYear int32  `bson:"CNSTRCT_YR"`
	Name      string `bson:"NAME"`
	// LSTMODDATE string
	// LSTSTATYPE string
	// DOITT_ID   int32
	HeightRoof  float64 `bson:"HEIGHTROOF"`
	FeatCode    int32   `bson:"FEAT_CODE"`
	GroundLevel int32   `bson:"GROUNDELEV"`
	ShapeArea   float64 `bson:"SHAPE_AREA"`
	// SHAPE_LEN  float64
	// BASE_BBL   int64
	// MPLUTO_BBL int64
	// GEOMSOURCE string
}

func main() {
	router := gin.Default()

	ctx := context.TODO()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)

	// Experiment with mgo driver
	// session, err := mgo.DialWithInfo(&mgo.DialInfo{
	// 	Addrs: []string{"127.0.0.1:27017"},
	// })
	// c := session.DB("topos").C("test")
	// fmt.Println("collection: ", c)

	if err != nil {
		log.Fatal(err)
	}

	type Trainer struct {
		// ID   bson.ObjectId `bson:"_id,omitempty"`
		Name string `bson:"name"`
		Age  int32  `bson:"age"`
		City string `bson:"city"`
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	fmt.Println("Connected to MongoDB!")
	// dbs, err := client.ListDatabases(context.TODO())

	// fmt.Println("db list: ", dbs)
	collection := client.Database("topos").Collection("testCollection")
	fmt.Println("collection: ", collection)

	// ash := Trainer{"Ash", 10, "Pallet Town"}
	// misty := Trainer{"Misty", 10, "Cerulean City"}
	// brock := Trainer{"Brock", 15, "Pewter City"}

	// trainers := []interface{}{ash, misty, brock}

	// insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	router.GET("/", routes.NotFound())
	router.GET("/trial", routes.Trial(client))
	router.GET("/building", routes.BuildingByName(client))
	router.GET("/buildingYear", routes.BuildingByConstructionYear(client))
	// router.GET("/", routes.NotFound())

	router.Run()

}
