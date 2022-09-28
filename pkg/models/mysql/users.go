package mysql

import (
	"database/sql"
	"ikehakinyemi/go-pastebin/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert adds new record to user table.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies a user exists with provide parameters. 
// Return user ID if exists.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}