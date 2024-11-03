package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/srini981/pismoTask/docs"
	"github.com/srini981/pismoTask/handler"
	"github.com/srini981/pismoTask/models"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// all the handlers to handle all the incoming http requests for all micrservices
func AddRoutes(route *gin.Engine) {
	route.Handle(http.MethodGet, "health", func(c *gin.Context) {
		c.JSON(200, models.Response{Err: nil, Message: " application is up and running successfully"})
		return
	})

	route.Handle(http.MethodPost, "/accounts", handler.CreateAccount)
	route.Handle(http.MethodGet, "/accounts/:ID", handler.GetAccount)
	route.Handle(http.MethodPost, "/transactions", handler.CreateTransaction)
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
