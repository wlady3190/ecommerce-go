package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wlady3190/ecommerce/database"
	"github.com/wlady3190/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

// Creando constructor o inicializando la función
func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}

}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("productId is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productId is empty"))
			return
		}
		userQueryID := c.Query("userId")
		if userQueryID == "" {
			log.Println("userId is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "successfully added to cart")

	}

}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("productId is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productId is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("userId is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = database.RemoveCartItem(ctx, *app.prodCollection, *app.userCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "successfully remove from cart")

	}

}

func GetItemFromCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User

		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Not found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointcursor, err :=UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		response := gin.H{
			"total": 0,
			"cart": filledcart.UserCart,
		}

		// for _, json := range listing{
		// 	c.IndentedJSON(200, json["total"])
		// 	c.IndentedJSON(200, filledcart.UserCart)
		// }
		// ctx.Done()

		if len(listing) >0{
			if total, ok := listing[0]["total"]; ok {
				response["total"] = total
			}
		}

		c.IndentedJSON(http.StatusOK, response)



	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Panic("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userID is empty"))
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully placed the order")

	}

}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {

		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("productId is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productId is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("userId is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully place the order")

	}
}
