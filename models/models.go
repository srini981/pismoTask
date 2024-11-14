package models

type Accounts struct {
	ID             int64 `gorm:"primaryKey"`
	DocumentNumber int64 `gorm:"unique" json:"document_number"`
	AccountName    string
}

type Response struct {
	Err      error
	Response interface{}
	Message  string
}

type OperationTypes struct {
	OperationTypeID int64  `gorm:"primaryKey"`
	Description     string `gorm:"unique"`
}

type Transactions struct {
	TransactionId   int64   `gorm:"primaryKey"`
	AccountID       int64   `json:"account_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
	Balance         float64 `json:"balance"`
	EventDate       string  `json:"event_date"`
}
