package handler

// handlers/handlers_test.go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/srini981/pismoTask/models"
	"github.com/stretchr/testify/mock"
)

// Mock the services package
type MockService struct {
	mock.Mock
}

type TestCases struct {
	account                           []byte
	mockMethodGetAccountinput         int64
	mockMethodGetAccountoutput        mockMethodGetAccountoutput
	mockMethodCreateAccountinput      models.Accounts
	mockMethodCreateTransactiontinput models.Transactions
	mockMethodCreateTransactionoutput error
	mockMethodCreateAccountoutput     error
	statusCode                        int
	err                               error
}

type mockMethodGetAccountoutput struct {
	Account models.Accounts
	err     error
}

func (m *MockService) mockMethodGetAccount(ctx context.Context, account int64) (models.Accounts, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(models.Accounts), args.Error(1)
}

func (m *MockService) mockMethodCreateAccount(ctx context.Context, account models.Accounts) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func TestCreateAccount(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router for testing

	// Test case 1: No name query parameter
	account := models.Accounts{
		DocumentNumber: 1234,
	}

	testcases := []TestCases{{account: []byte(`{"documentNmber":12434}`),
		mockMethodGetAccountinput: int64(1234), mockMethodGetAccountoutput: mockMethodGetAccountoutput{models.Accounts{}, fmt.Errorf(fmt.Sprintf("failed to fetch account details from db "))}, mockMethodCreateAccountinput: models.Accounts{
			DocumentNumber: 1234,
		}, mockMethodCreateAccountoutput: nil, statusCode: http.StatusBadRequest},
		{account: []byte(`{"documentNmber":12434}`), mockMethodGetAccountinput: int64(1234), mockMethodGetAccountoutput: mockMethodGetAccountoutput{models.Accounts{}, nil}, mockMethodCreateAccountinput: models.Accounts{
			DocumentNumber: 1234,
		}, mockMethodCreateAccountoutput: nil, statusCode: http.StatusOK},
		{account: []byte(`{"documentNmber":12434},}`), mockMethodGetAccountinput: int64(1234), mockMethodGetAccountoutput: mockMethodGetAccountoutput{models.Accounts{}, fmt.Errorf(fmt.Sprintf("failed to fetch account details from db "))}, mockMethodCreateAccountinput: models.Accounts{
			DocumentNumber: 1234,
		}, mockMethodCreateAccountoutput: nil, statusCode: http.StatusBadRequest}}

	for _, testcase := range testcases {
		r := gin.Default()
		mockService := new(MockService)

		// Replace the actual service function with a mock
		r.POST("/accounts", func(c *gin.Context) {
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

			responseaccount, err := mockService.mockMethodGetAccount(c.Request.Context(), account.DocumentNumber)

			if err != nil || responseaccount.DocumentNumber == account.DocumentNumber {

				response := models.Response{Err: err, Message: "failed to create  account with document number already exists"}
				c.JSON(http.StatusBadRequest, response)
				return

			}

			err = mockService.mockMethodCreateAccount(c.Request.Context(), account)

			if err != nil {

				response := models.Response{Err: err, Message: "failed to create  account with document number"}
				c.JSON(http.StatusBadRequest, response)
				return

			}

			response := models.Response{Err: nil, Message: "account created with document number"}
			c.JSON(http.StatusOK, response)

		})

		mockService.On("mockMethodGetAccount", context.Background(), testcase.mockMethodGetAccountinput).Return(testcase.mockMethodGetAccountoutput.Account, testcase.mockMethodGetAccountoutput.err)

		mockService.On("mockMethodCreateAccount", context.Background(), testcase.mockMethodCreateAccountinput).Return(testcase.mockMethodCreateAccountoutput).Maybe()
		w := httptest.NewRecorder()

		bode, err := json.Marshal(account)
		if err != nil {

		}
		b := bytes.NewReader(bode)
		req, _ := http.NewRequest(http.MethodPost, "/accounts", b)
		r.ServeHTTP(w, req)
		assert.Equal(t, testcase.statusCode, w.Code)
		mockService.AssertExpectations(t)

	}
}

func TestGetAccount(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router for testing

	// Test case 1: No name query parameter
	account := models.Accounts{
		DocumentNumber: 1234,
	}

	testcases := []TestCases{{account: []byte(`{"documentNmber":12434}`),
		mockMethodGetAccountinput: int64(1234), mockMethodGetAccountoutput: mockMethodGetAccountoutput{models.Accounts{}, fmt.Errorf(fmt.Sprintf("failed to fetch account details from db "))}, statusCode: http.StatusBadRequest},
		{account: []byte(`{"documentNmber":12434}`), mockMethodGetAccountinput: int64(1234), mockMethodGetAccountoutput: mockMethodGetAccountoutput{models.Accounts{}, nil}, statusCode: http.StatusOK},
		{account: []byte(`{"documentNmber":12434},}`), mockMethodGetAccountinput: int64(1234), mockMethodGetAccountoutput: mockMethodGetAccountoutput{models.Accounts{}, fmt.Errorf(fmt.Sprintf("failed to fetch account details from db "))}, statusCode: http.StatusBadRequest}}

	for _, testcase := range testcases {
		r := gin.Default()
		mockService := new(MockService)

		// Replace the actual service function with a mock
		r.GET("/accounts/:ID", func(c *gin.Context) {
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
			account, err := mockService.mockMethodGetAccount(ctx, int64(id))
			if err != nil {

				err := fmt.Errorf("failed to get account details %s", err.Error())
				log.Println(err)
				Response := models.Response{Err: err, Message: "failed to get account details"}
				c.JSON(http.StatusBadRequest, Response)
				return

			}

			response := models.Response{Err: nil, Message: "account details fetched successfully with id", Response: account}
			c.JSON(http.StatusOK, response)

		})

		mockService.On("mockMethodGetAccount", context.Background(), testcase.mockMethodGetAccountinput).Return(testcase.mockMethodGetAccountoutput.Account, testcase.mockMethodGetAccountoutput.err)

		w := httptest.NewRecorder()

		bode, err := json.Marshal(account)
		if err != nil {

		}
		b := bytes.NewReader(bode)
		req, _ := http.NewRequest(http.MethodGet, "/accounts/1234", b)
		r.ServeHTTP(w, req)
		assert.Equal(t, testcase.statusCode, w.Code)
		mockService.AssertExpectations(t)
	}
}
