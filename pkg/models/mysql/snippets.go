package mysql

import (
	"database/sql"
	"errors"
	"ikehakinyemi/go-pastebin/pkg/models"
)

// SnippetModel will wrap a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert will insert a new snippet into the database.
func (m *SnippetModel) Insert(userEmail, title, content, expires string) (int, error) {
	// SQL statement
	stmt := `INSERT INTO snippets (title, content, created, expires, owner_email)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`
	// Use Exec() to execute statement
	result, err := m.DB.Exec(stmt, title, content, expires, userEmail)
	if err != nil {
		return 0, err
	}
	// Use LastInsertId() method to get ID of above inserted record.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Initialize an empty models.Snippets objects struct to hold return record.
	s := &models.Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Latest will return the 10 most recently created snippets.
func (m *SnippetModel) Latest(userEmail string) ([]*models.Snippet, error) {
	// SQL statement.
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND owner_email = ? ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize empty slice to hold models.Snippets objects.
	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
