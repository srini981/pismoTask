package database

import (
	"context"
	"testing"
	"time"

	"github.com/srini981/pismoTask/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	dsn := "host=localhost user=pg password=pass dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Accounts{})
	db.AutoMigrate(&models.OperationTypes{})
	db.AutoMigrate(&models.Transactions{})
	return db, err
}

func TestCreateAccount(t *testing.T) {
	client, err := setupTestDB()
	assert.NoError(t, err)

	defer dbclient.Exec("DROP TABLE accounts;") // Clean up
	dbclient = client

	user := models.Accounts{DocumentNumber: 123124214}
	ctx := context.Background()
	err = Client.CreateAccount(ctx, user)
	assert.Equal(t, nil, err)
	// Test duplicate email
	err = Client.CreateAccount(ctx, user)
	assert.Error(t, err)
	// Should return an error due to unique constraint
}

func TestGetAccount(t *testing.T) {
	client, err := setupTestDB()
	assert.NoError(t, err)
	dbclient = client
	defer dbclient.Exec("delete from  accounts where id =1")

	user := models.Accounts{DocumentNumber: 123124214}
	ctx := context.Background()
	err = Client.CreateAccount(ctx, user)
	assert.Equal(t, nil, err)
	// Test duplicate email
	_, err = Client.GetAccount(ctx, 1)
	assert.Equal(t, nil, err)
	_, err = Client.GetAccount(ctx, 3)
	assert.Error(t, err)
	// Should return an error due to unique constraint
}

func TestGetAccountByDocumentID(t *testing.T) {
	client, err := setupTestDB()
	assert.NoError(t, err)
	dbclient = client
	defer dbclient.Exec("delete from  accounts where id =1")

	user := models.Accounts{DocumentNumber: 123124214}
	ctx := context.Background()
	err = Client.CreateAccount(ctx, user)
	assert.Equal(t, nil, err)
	// Test duplicate email
	_, err = Client.GetAccountByDocumentNumber(ctx, user.DocumentNumber)
	assert.Equal(t, nil, err)
	_, err = Client.GetAccountByDocumentNumber(ctx, 3)
	assert.Error(t, err)
	// Should return an error due to unique constraint
}

func TestCreateTransaction(t *testing.T) {
	client, err := setupTestDB()
	assert.NoError(t, err)

	defer dbclient.Exec("DROP TABLE transactions;") // Clean up
	dbclient = client

	user := models.Transactions{AccountID: 1, OperationTypeID: 3, Amount: 50, EventDate: time.Time{}.String()}
	ctx := context.Background()
	err = Client.CreateTransaction(ctx, user)
	assert.Equal(t, nil, err)
	// Test duplicate email
}
