package routes

import (
	"context"
	"fmt"
	// "github.com/fatih/structs"
	"geoServer/errors"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"reflect"
)

type AggResult struct {
	ID    int64   `bson:"_id"`
	Value float64 `bson:"value"`
}

var aggOpMap = map[string]string{ // Mapping aggregation api values to mongo query
	"mean": "$avg",
	"min":  "$min",
	"max":  "$max",
}

var aggColMap = map[string]string{
	"Bin":        "$BIN",
	"ConstYear":  "$CNSTRCT_YR",
	"Name":       "$NAME",
	"HeightRoof": "$HEIGHTROOF",
	"FeatCode":   "$FEAT_CODE",
	"GroundElev": "$GROUNDELEV",
	"ShapeArea":  "$SHAPE_AREA",
}

type AggQuery struct {
	Field string
	AggBy string
}

var validAgg = gin.H{
	"validFields":                `Bin, ConstYear, Name, HeightRoof, FeatCode, GroundElev, ShapeArea`,
	"validAggregatingOperations": `mean, min, max`,
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

		if aggOpMap[q.AggBy] == "" || aggColMap[q.Field] == "" {
			errors.BadAggregate(c)
			return
		}

		if aggOpMap[q.AggBy] == "Name" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Good one! You cannot aggregate names. No, seriously!",
			})
			return
		}

		query := []bson.M{
			bson.M{
				"$group": bson.M{
					"_id":   0,
					"value": bson.M{aggOpMap[q.AggBy]: aggColMap[q.Field]},
				},
			}}

		fmt.Println("query: ", query)
		var aggResult AggResult

		cur, err := collection.Aggregate(context.TODO(), query)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err,
				"meta":  validAgg,
			})
			// log.Panic("something wrong in query: ", err)
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

type filAggQuery struct {
	Filter    Building `json:"filter"`
	Aggregate AggQuery `json:"aggregate"`
}

func GetFilteredAggregatedValue(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
		collection := client.Database("topos").Collection("testCollection")

		// var building Building
		var q filAggQuery
		if c.BindJSON(&q) == nil {
			fmt.Println("url query: ", q)
		}
		qAgg := q.Aggregate
		qFilter := q.Filter

		if aggOpMap[qAgg.AggBy] == "" || aggColMap[qAgg.Field] == "" {
			errors.BadAggregate(c)
			return
		}

		if aggOpMap[qAgg.AggBy] == "Name" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Good one! You cannot aggregate names. No, seriously!",
			})
			return
		}

		m := structs.Map(qFilter)
		filterQuery := bson.M{}

		for k, v := range m {
			fmt.Println("key: ", bsonMap[k])
			fmt.Println("value: ", v)
			// @TODO: Need better BindQuery mechanism to filter out optional params from query
			if v != reflect.Zero(reflect.TypeOf(v)).Interface() { // Has a non 0 or "" value. Seriously Go??
				filterQuery[bsonMap[k]] = v
			}
		}

		aggregateQuery := createAggQuery(qAgg)

		tempFilterQuery := []bson.M{}
		tempFilterQuery = append(tempFilterQuery, bson.M{"$match": filterQuery})
		aggregateQuery = append(tempFilterQuery, aggregateQuery...)
		var aggResult AggResult
		cur, err := collection.Aggregate(context.TODO(), aggregateQuery)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err,
				"meta":  validAgg,
			})
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			fmt.Println("--------------cur: ", cur)
			err := cur.Decode(&aggResult)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"error": err,
				})
				log.Panic(err)
			}
		}
		if err := cur.Err(); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err,
			})
			log.Panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"meeesg": aggResult,
		})
	}
}

func createAggQuery(q AggQuery) []bson.M {
	query := []bson.M{
		bson.M{
			"$group": bson.M{
				"_id":   0,
				"value": bson.M{aggOpMap[q.AggBy]: aggColMap[q.Field]},
			},
		}}
	return query
}
