package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokha04/go-e-commerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			ctx.Abort()
			return
		}

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var addresses models.Address
		addresses.Address_ID = primitive.NewObjectID()

		if err = ctx.BindJSON(&addresses); err != nil {
			ctx.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userId}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := UserCollection.Aggregate(c, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var addressInfo []bson.M
		err = pointCursor.All(c, &addressInfo)
		if err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: userId}}
			update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err = UserCollection.UpdateOne(c, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			ctx.IndentedJSON(http.StatusBadRequest, "not allowed")
		}

		defer cancel()
		c.Done()
	}
}

func EditHomeAddresss() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
			ctx.Abort()
			return
		}

		var editAddress models.Address
		if err := ctx.BindJSON(&editAddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key: "address.0.street_name", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "something went wrong")
			return
		}

		defer cancel()
		c.Done()
		ctx.IndentedJSON(http.StatusOK, "successfully updated home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
			ctx.Abort()
			return
		}

		var editAddress models.Address
		if err := ctx.BindJSON(&editAddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "something went wrong")
			return
		}

		defer cancel()
		c.Done()
		ctx.IndentedJSON(http.StatusOK, "successfully updated work address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			ctx.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, "wrong command")
			return
		}

		defer cancel()
		c.Done()
		ctx.IndentedJSON(http.StatusOK, "successfully deleted")
	}
}
