package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/srini981/pismoTask/database"
	"github.com/srini981/pismoTask/models"
)

// credit voucher operation type variable
var creditVoucher int64 = 4

// @Summary create transaction api
// @Description create transaction api for creating transaction
// @Tags transaction
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /transaction [post]
func CreateTransaction(c *gin.Context) {

	transaction := models.Transactions{}
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("unable to read json body")
		Response := models.Response{Err: err, Message: "unable to read json body"}
		c.JSON(http.StatusBadRequest, Response)
		return
	}

	err = json.Unmarshal(reqBody, &transaction)

	if err != nil {
		log.Println("unable to parse json body")
		Response := models.Response{Err: err, Message: "unable to parse json body"}
		c.JSON(http.StatusBadRequest, Response)
		return
	}

	ctx := context.Background()
	defer ctx.Done()

	if transaction.OperationTypeID != creditVoucher && transaction.Amount > 0 || (transaction.OperationTypeID == creditVoucher && transaction.Amount < 0) {
		transaction.Amount = -1 * transaction.Amount
	}

	transaction.Balance = transaction.Amount
	transaction.EventDate = time.Now().String()

	if transaction.OperationTypeID == creditVoucher {
		err = database.Client.Discharge(ctx, transaction)

	} else {
		err = database.Client.CreateTransaction(ctx, transaction)
	}

	if err != nil {
		response := models.Response{Err: err, Message: "failed to create  transaction"}
		c.JSON(http.StatusBadRequest, response)
		return

	}

	response := models.Response{Err: nil, Message: "transaction created "}
	c.JSON(http.StatusOK, response)
}
