package routes

import (
	"context"
	"fmt"
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
	HeightRoof  float64 `bson:"HEIGHTROOF"`
	FeatCode    int32   `bson:"FEAT_CODE"`
	GroundLevel int32   `bson:"GROUNDELEV"`
	ShapeArea   float64 `bson:"SHAPE_AREA"`
	// SHAPE_LEN  float64
	// BASE_BBL   int64
	// MPLUTO_BBL int64
	// GEOMSOURCE string
}

var bsonMap = map[string]string{
	"Bin":         "BIN",
	"ConstYear":   "CNSTRCT_YR",
	"Name":        "NAME",
	"HeightRoof":  "HEIGHTROOF",
	"FeatCode":    "FEAT_CODE",
	"GroundLevel": "GROUNDELEV",
	"ShapeArea":   "SHAPE_AREA",
}

func BuildingByName(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		collection := client.Database("topos").Collection("testCollection")
		fmt.Println("collection: ", collection)
		filter := bson.D{{"NAME", "Pheasant Aviary"}}
		var result Building
		// var result Trainer

		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found a single document: %+v\n", result.FeatCode)
		c.JSON(http.StatusOK, gin.H{
			"msg": result.ShapeArea,
		})
	}
}

func BuildingByConstructionYear(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
		collection := client.Database("topos").Collection("testCollection")
		filter := bson.D{{"CNSTRCT_YR", 2016}}
		var result Building
		// var result Trainer

		cur, err := collection.Find(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Found document: %+v\n", result)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found a single document: %+v\n", result.FeatCode)
		c.JSON(http.StatusOK, gin.H{
			"msg": result,
		})
	}
}

func getTag(name string) string {
	t := reflect.TypeOf(Building{})
	x, _ := t.FieldByName(name)
	return x.Tag.Get("bson")
}

func Trial(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
		collection := client.Database("topos").Collection("testCollection")
		// q := c.Request.URL.Query()
		// params := c.Param("vals")
		// fmt.Printf("params: ", params)
		var building Building

		if c.BindQuery(&building) == nil {
			fmt.Println(building)
		}

		m := structs.Map(building)
		// _v := structs.Values(building)
		// queries := []bson.M{}
		queries := bson.M{}

		// fields := reflect.TypeOf(building)
		// // values := reflect.ValueOf(building)
		// num := fields.NumField()

		// for i := 0; i < num; i++ {
		// 	field := fields.Field(i)
		// 	// value := values.Field(i)

		// 	t := reflect.TypeOf(Building{})
		// 	x, _ := t.FieldByName(field.Name)

		// 	bsonType := x.Tag.Get("bson")

		// 	// bson_type := getTag(field)
		// 	r := reflect.ValueOf(building)
		// 	f := reflect.Indirect(r).FieldByName(field.Name)
		// 	query[bsonType] = f
		// 	fmt.Println("\nfstring: ", f)
		// 	// fmt.Println("\nfstring type: ", reflect.TypeOf(f))
		// 	fmt.Println("\n bsonType: ", bsonType)
		// 	fmt.Println("\n bsonType: ", bsonType)
		// }

		for k, v := range m {
			// query[k] = v
			fmt.Println("key: ", bsonMap[k])
			fmt.Println("value: ", v)
			if v != reflect.Zero(reflect.TypeOf(v)).Interface() { // Has a non NIL value
				fmt.Println("---------------------Its non zero!!")
				queries[bsonMap[k]] = v
			}
		}

		// for k, v := range _v {
		// 	// query[k] = v
		// 	fmt.Println("key: ", k)
		// 	fmt.Println("value: ", v)
		// }

		fmt.Println("queries: ", queries)
		var result Building
		var results = []Building{}
		cur, err := collection.Find(ctx, queries)
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, result)
			fmt.Printf("Found document: %+v\n", result)
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
