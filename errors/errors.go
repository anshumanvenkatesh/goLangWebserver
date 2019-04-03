package errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var validAgg = gin.H{
	"validFields":                `Bin, ConstYear, Name, HeightRoof, FeatCode, GroundElev, ShapeArea`,
	"validAggregatingOperations": `mean, min, max`,
}

func BadAggregate(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Oops! Either aggregating field or" +
			" aggregating operation or both is/are invalid" +
			" (maybe even the grammatical structure of this error message?)",
		"meta": validAgg,
	})
}

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "Umm... Maybe you were aiming for /aggregate or /buildings? That's all for the show now.",
		})
	}
}

func DatabaseError(c *gin.Context, err error) {
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"msg":     "Our DB exploded while trying to service you.",
		"details": err,
	})
}
