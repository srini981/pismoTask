package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/srini981/pismoTask/router"
)

func main() {
	log.Println("starting  service at port 8004")

	ginRouter := gin.New()

	log.Println("initizlization database")

	log.Println(" database initilized successfully")

	log.Println("adding routes to service")

	router.AddRoutes(ginRouter)

	log.Println("added routes to  service")
	ginRouter.Use(gin.Logger())
	ginRouter.Run(":8004")

	log.Println("service started")
}
