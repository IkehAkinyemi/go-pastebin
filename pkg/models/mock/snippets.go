package mock

import (
	"ikehakinyemi/go-pastebin/pkg/models"
	"time"
)

var mockSnippet = &models.Snippet{
	ID: 1,
	Title: "Dead memories erupts me",
	Content: "Forgot, this is a sad record",
	Expires: time.Now(),
	Created: time.Now(),
}

type SnippetModel struct {}

func (m *SnippetModel) Insert(userEmail, title, content, expires string) (int, error)  {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest(userEmail string) ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}