package database

var (
	accountByIDQuery        = "select * from accounts where id = %d;"
	accountByDocumentNumber = "select * from accounts where Document_number = %d;"
)
