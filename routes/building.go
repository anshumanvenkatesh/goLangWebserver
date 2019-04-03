package routes

import (
	"context"
	"fmt"
	"geoServer/errors"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"reflect"
)

type Building struct {
	// _id
	// the_geom   string
	Bin       int32  `bson:"BIN" `
	ConstYear int32  `bson:"CNSTRCT_YR" `
	Name      string `bson:"NAME" binding:"exists"`
	// LSTMODDATE string
	// LSTSTATYPE string
	// DOITT_ID   int32
	HeightRoof float64 `bson:"HEIGHTROOF"`
	FeatCode   int32   `bson:"FEAT_CODE"`
	GroundElev int32   `bson:"GROUNDELEV"`
	ShapeArea  float64 `bson:"SHAPE_AREA"`
	// SHAPE_LEN  float64
	// BASE_BBL   int64
	// MPLUTO_BBL int64
	// GEOMSOURCE string
}

var bsonMap = map[string]string{ // Built for easy lookup for bson tags
	"Bin":        "BIN",
	"ConstYear":  "CNSTRCT_YR",
	"Name":       "NAME",
	"HeightRoof": "HEIGHTROOF",
	"FeatCode":   "FEAT_CODE",
	"GroundElev": "GROUNDELEV",
	"ShapeArea":  "SHAPE_AREA",
}

func Intro() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ctx := context.TODO(),
		c.JSON(http.StatusOK, gin.H{
			"apiEndpoints":              "/buildings(?filterParams), /aggregate(?aggregationParam), /",
			"filterParamsFormat":        "/buildings?<dataFeature1>=<somevalue>&<dataFeature2>=<some other value>& ...",
			"aggregationParamsFormat":   "/aggregate?Field=<dataFeature>&AggBy=<Aggregation operator>",
			"filterParamsExamples":      "/buildings?ConstYear=1922&Name=AlphaHouse",
			"aggregationParamsExamples": "/aggregate?Field=ShapeArea&AggBy=mean",
			"aggregationOperators":      "mean, min, max",
			"dataFeatures":              "Bin, ConstYear, Name, HeightRoof, FeatCode, GroundElev, ShapeArea",
		})
	}
}

func GetBuildingsData(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
		collection := client.Database("topos").Collection("testCollection")
		var building Building

		if c.BindQuery(&building) == nil {
			fmt.Println(building)
		}

		m := structs.Map(building)
		query := bson.M{}

		for k, v := range m {
			fmt.Println("key: ", bsonMap[k])
			fmt.Println("value: ", v)
			// @TODO: Need better BindQuery mechanism to filter out optional params from query
			if v != reflect.Zero(reflect.TypeOf(v)).Interface() { // Has a non 0 or "" value. Seriously Go??
				query[bsonMap[k]] = v
			}
		}

		fmt.Println("query: ", query)
		var result Building
		var results = []Building{}
		cur, err := collection.Find(ctx, query)
		if err != nil {
			errors.DatabaseError(c, err)
			log.Fatal(err)
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, result)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": gin.H{
				"vals": results,
			},
		})
	}
}
