package routes

import (
	"context"
	"fmt"
	// "github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	// "reflect"
)

// type Building struct {
// 	// _id
// 	// the_geom   string
// 	Bin       int32  `bson:"BIN" `
// 	ConstYear int32  `bson:"CNSTRCT_YR" `
// 	Name      string `bson:"NAME" binding:"exists"`
// 	// LSTMODDATE string
// 	// LSTSTATYPE string
// 	// DOITT_ID   int32
// 	HeightRoof  float64 `bson:"HEIGHTROOF"`
// 	FeatCode    int32   `bson:"FEAT_CODE"`
// 	GroundLevel int32   `bson:"GROUNDELEV"`
// 	ShapeArea   float64 `bson:"SHAPE_AREA"`
// 	// SHAPE_LEN  float64
// 	// BASE_BBL   int64
// 	// MPLUTO_BBL int64
// 	// GEOMSOURCE string
// }

// var bsonMap = map[string]string{ // Built for easy lookup for bson tags
// 	"Bin":         "BIN",
// 	"ConstYear":   "CNSTRCT_YR",
// 	"Name":        "NAME",
// 	"HeightRoof":  "HEIGHTROOF",
// 	"FeatCode":    "FEAT_CODE",
// 	"GroundLevel": "GROUNDELEV",
// 	"ShapeArea":   "SHAPE_AREA",
// }

type AggResult struct {
	Id    int64   `bson:"_id"`
	Value float64 `bson:"value"`
}

var aggOpMap = map[string]string{ // Mapping aggregation api values to mongo query
	"mean": "$avg",
	"min":  "$min",
	"max":  "$max",
}

var aggColMap = map[string]string{
	"Bin":         "$BIN",
	"ConstYear":   "$CNSTRCT_YR",
	"Name":        "$NAME",
	"HeightRoof":  "$HEIGHTROOF",
	"FeatCode":    "$FEAT_CODE",
	"GroundLevel": "$GROUNDELEV",
	"ShapeArea":   "$SHAPE_AREA",
}

type AggQuery struct {
	Collection string
	AggBy      string
}

func GetAggregatedValue(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
		collection := client.Database("topos").Collection("testCollection")
		fmt.Println("collection: ", collection)

		var q AggQuery

		if c.BindQuery(&q) == nil {
			fmt.Println("url query: ", q)
		}

		query := []bson.M{
			bson.M{
				"$group": bson.M{
					"_id":   0,
					"value": bson.M{aggOpMap[q.AggBy]: aggColMap[q.Collection]},
				},
			}}

		fmt.Println("query: ", query)
		var aggResult AggResult

		cur, err := collection.Aggregate(context.TODO(), query)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err,
			})
			log.Panic("something wrong in query: ", err)
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			fmt.Println("cur: ", cur)
			err := cur.Decode(&aggResult)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"error": err,
				})
				log.Panic(err)
			}
			fmt.Printf("Found document: %+v\n", aggResult)
		}
		if err := cur.Err(); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err,
			})
			log.Panic(err)
		}

		fmt.Printf("Found a single document: %+v\n", aggResult)
		c.JSON(http.StatusOK, gin.H{
			"msg": aggResult.Value,
		})
	}
}
