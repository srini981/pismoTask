package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/srini981/pismoTask/database"
	"github.com/srini981/pismoTask/models"
)

// @Summary create account api
// @Description create account api for creating accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /accounts [post]
func CreateAccount(c *gin.Context) {
	account := models.Accounts{}
	reqBody, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {

		log.Println("unable to read json body")
		Response := models.Response{Err: err, Message: "unable to read json body"}
		c.JSON(http.StatusBadRequest, Response)
		return
	}

	err = json.Unmarshal(reqBody, &account)

	if err != nil {
		log.Println("unable to parse json body")
		Response := models.Response{Err: err, Message: "unable to parse json body"}
		c.JSON(http.StatusBadRequest, Response)
		return
	}

	ctx := context.Background()
	defer ctx.Done()
	responseaccount, err := database.Client.GetAccountByDocumentNumber(ctx, account.DocumentNumber)

	if err == nil || responseaccount.DocumentNumber == account.DocumentNumber {
		response := models.Response{Err: err, Message: "failed to create  account with document number already exists"}
		c.JSON(http.StatusConflict, response)
		return

	}

	err = database.Client.CreateAccount(ctx, account)

	if err != nil {

		response := models.Response{Err: err, Message: "failed to create  account with document number"}
		c.JSON(http.StatusBadRequest, response)
		return

	}

	response := models.Response{Err: nil, Message: "account created with document number"}
	c.JSON(http.StatusOK, response)

}

// @Summary get account api
// @Description get account api for creating accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Accounts
// @Router /accounts/:ID [get]
// @Param id path int true "account ID"
func GetAccount(c *gin.Context) {
	accountID, exists := c.Params.Get("ID")

	if !exists {

		log.Println("user id not found in path")

		Response := models.Response{Err: errors.New("user ID required"), Message: "user ID not found in path"}
		c.JSON(http.StatusBadRequest, Response)
		return

	}

	id, err := strconv.Atoi(accountID)
	if err != nil {

		log.Println("invalid account id")

		Response := models.Response{Err: errors.New("invalid account id"), Message: "invalid account id"}
		c.JSON(http.StatusBadRequest, Response)
		return

	}

	ctx := context.Background()
	defer ctx.Done()
	account, err := database.Client.GetAccount(ctx, int64(id))
	if err != nil {

		err := fmt.Errorf("failed to get account details %s", err.Error())
		log.Println(err)
		Response := models.Response{Err: err, Message: "failed to get account details"}
		c.JSON(http.StatusBadRequest, Response)
		return

	}

	c.JSON(http.StatusOK, account)

}
