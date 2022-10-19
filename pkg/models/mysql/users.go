package mysql

import (
	"database/sql"
	"errors"
	"ikehakinyemi/go-pastebin/pkg/models"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// Insert adds new record to user table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) 
	VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate verifies a user exists with provide parameters.
// Return user ID if exists.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return 0, models.ErrInvalidCredentials
			} else {
				return 0, err
			}
		}
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}
	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return u, nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	var hashedPassword []byte
	stmt := `SELECT hashed_password FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrInvalidCredentials
		} else {
			return err
		}

	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	stmt = `UPDATE users SET hashed_password = ? WHERE id = ?`
	_, err = m.DB.Exec(stmt, string(newPasswordHash), id)
	return err
}
