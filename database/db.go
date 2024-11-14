package database

import (
	"context"
	"fmt"
	"log"

	"github.com/srini981/pismoTask/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// dbclient to handle all db request
type client struct{}

var Client client

var dbclient = DB()

// interface for declaring methods for db
type clientInterface interface {
	GetAccount(context.Context, int64) (models.Accounts, error)
	GetAccountByDocumentNumber(context.Context, int64) (models.Accounts, error)
	CreateAccount(context.Context, models.Accounts) error
	CreateTransaction(context.Context, models.Transactions) error
	Discharge(context.Context, models.Transactions) error
}

// init function for postgres and redis clients
func DB() *gorm.DB {
	dsn := "host=localhost user=pg password=pass dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {

		log.Fatal("failed to connect postgres database")
		return nil
	}

	db.AutoMigrate(&models.Accounts{})
	db.AutoMigrate(&models.OperationTypes{})
	db.AutoMigrate(&models.Transactions{})
	OperationTypes := []models.OperationTypes{{1, "Normal Purchase"}, {2, "Purchase with installments"}, {3, "Withdrawal"}, {4, "Credit Voucher"}}
	db.Create(&OperationTypes)
	postgresdb, err := db.DB()

	if err != nil {

		log.Fatal("failed to create postgres db object")
		return nil

	}

	err = postgresdb.Ping()

	if err != nil {

		log.Fatal("failed to Ping postgres database")
		return nil

	}

	return db
}

// get account function to get account details
func (d *client) GetAccount(ctx context.Context, accountID int64) (models.Accounts, error) {
	account := models.Accounts{}
	query := fmt.Sprintf(accountByIDQuery, accountID)
	log.Println(query)
	tx := dbclient.Raw(query).Scan(&account)

	if (tx.Error != nil && tx.Error.Error() != "") || tx.RowsAffected == 0 {

		err := fmt.Errorf(fmt.Sprintf("failed to fetch account details from db"))
		return models.Accounts{}, err

	}

	return account, nil
}

// get account function to get account details by document number
func (d *client) GetAccountByDocumentNumber(ctx context.Context, documentNumber int64) (models.Accounts, error) {
	account := models.Accounts{}
	query := fmt.Sprintf(accountByDocumentNumber, documentNumber)
	tx := dbclient.Raw(query).Scan(&account)

	if (tx.Error != nil && tx.Error.Error() != "") || tx.RowsAffected == 0 {

		err := fmt.Errorf(fmt.Sprintf("failed to fetch account details from db "))
		return models.Accounts{}, err

	}

	return account, nil

}

// create account function to create account record in db
func (d *client) CreateAccount(ctx context.Context, account models.Accounts) error {
	tx := dbclient.Create(&account)

	if tx.Error != nil && tx.Error.Error() != "" {

		err := fmt.Errorf(fmt.Sprintf("failed to create account in db %s", tx.Error.Error()))
		return err

	}

	return nil

}

// create transaction function to create transaction record in db
func (d *client) CreateTransaction(ctx context.Context, transaction models.Transactions) error {
	tx := dbclient.Model(models.Transactions{}).Create(&transaction)

	if tx.Error != nil && tx.Error.Error() != "" {
		err := fmt.Errorf(fmt.Sprintf("failed to create transaction in db %s", tx.Error.Error()))
		return err
	}

	return nil
}

// to handle the discharge feature if an transaction of type credit voucher comes in
func (d *client) Discharge(ctx context.Context, currenttransaction models.Transactions) error {
	transctions := []models.Transactions{}

	err := dbclient.Transaction(func(tx *gorm.DB) error {

		// query the transactions in the db
		tx = dbclient.Model(models.Transactions{}).Where("account_id = ? and balance < 0", currenttransaction.AccountID).Scan(&transctions)

		if tx.Error != nil && tx.Error.Error() != "" || tx.RowsAffected == 0 {
			err := fmt.Errorf(fmt.Sprintf("failed to fetch transactions details in db %s", tx.Error.Error()))
			return err
		}

		currentBalance := currenttransaction.Amount

		// logic for handling the discharge
		for index, transaction := range transctions {

			if currentBalance == 0 {
				break
			}

			temp := (-1 * transaction.Balance) - currentBalance

			if temp <= 0 {

				currentBalance = currentBalance - (-1 * transaction.Balance)
				transctions[index].Balance = 0

			} else if temp > 0 {

				transctions[index].Balance = (-1 * temp)
				currentBalance = 0

			}

		}

		currenttransaction.Balance = currentBalance
		transctions = append(transctions, currenttransaction)

		currenttransaction.Balance = currentBalance

		// update in db
		tx = dbclient.Model(models.Transactions{}).Save(&transctions)

		if tx.Error != nil && tx.Error.Error() != "" || tx.RowsAffected == 0 {

			err := fmt.Errorf(fmt.Sprintf("failed to update transactions in db %s", tx.Error.Error()))
			return err

		}
		return nil
	})

	return err
}
