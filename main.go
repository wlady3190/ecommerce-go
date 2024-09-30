package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wlady3190/ecommerce/controllers"
	"github.com/wlady3190/ecommerce/database"
	"github.com/wlady3190/ecommerce/middleware"
	"github.com/wlady3190/ecommerce/routes"
)


func main()  {
	port := os.Getenv("PORT")
	if port ==""{
		port = "8080"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAdrress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAdress())
	router.GET("/deleteaddress", controllers.DeleteAdrress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":8080"))


 
}