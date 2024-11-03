package handler

// handlers/handlers_test.go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/srini981/pismoTask/models"
)

func (m *MockService) mockMethodCreateTransaction(ctx context.Context, account models.Transactions) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router for testing

	// Test case 1: No name query parameter
	account := models.Transactions{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          -50,
	}

	testcases := []TestCases{{
		account: []byte(`{"AccountID": 1,
		"OperationTypeID": 1,
		"Amount":          -50}`), mockMethodCreateTransactiontinput: models.Transactions{AccountID: 1,
			OperationTypeID: 1,
			Amount:          -50},
		mockMethodCreateTransactionoutput: fmt.Errorf(fmt.Sprintf("failed to fetch account details from db ")),
		statusCode:                        http.StatusBadRequest},
		{account: []byte(`{"AccountID": 1,
		"OperationTypeID": 1,
		"Amount":          -50`), mockMethodCreateTransactiontinput: models.Transactions{AccountID: 1,
			OperationTypeID: 1,
			Amount:          -50},
			mockMethodCreateTransactionoutput: fmt.Errorf(fmt.Sprintf("failed to fetch account details from db ")),
			statusCode:                        http.StatusBadRequest},
		{account: []byte(`{"AccountID": 1,
			"OperationTypeID": 1,
			"Amount":          -50}`),
			mockMethodCreateTransactiontinput: models.Transactions{AccountID: 1,
				OperationTypeID: 1,
				Amount:          -50},
			mockMethodCreateTransactionoutput: nil,
			statusCode:                        http.StatusOK},
	}

	for _, testcase := range testcases {
		r := gin.Default()
		mockService := new(MockService)

		// Replace the actual service function with a mock
		r.POST("/transactions", func(c *gin.Context) {
			account := models.Transactions{}
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

			err = mockService.mockMethodCreateTransaction(c.Request.Context(), account)

			if err != nil {

				response := models.Response{Err: err, Message: "failed to create  account with document number"}
				c.JSON(http.StatusBadRequest, response)
				return

			}

			response := models.Response{Err: nil, Message: "account created with document number"}
			c.JSON(http.StatusOK, response)

		})

		mockService.On("mockMethodCreateTransaction", context.Background(), testcase.mockMethodCreateTransactiontinput).Return(testcase.mockMethodCreateTransactionoutput)

		w := httptest.NewRecorder()

		bode, err := json.Marshal(account)
		if err != nil {
			log.Println("failed to parse json")
			return
		}
		b := bytes.NewReader(bode)
		req, _ := http.NewRequest(http.MethodPost, "/transactions", b)
		r.ServeHTTP(w, req)
		assert.Equal(t, testcase.statusCode, w.Code)
		mockService.AssertExpectations(t)
	}
}
