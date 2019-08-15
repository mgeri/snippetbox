package store

import (
	"errors"

	"github.com/mgeri/snippetbox/pkg/models"
)

var ErrNoRecord = errors.New("models: no matching record found")

// SnippetStore is the persistent store of Snippet
type SnippetStore interface {
	Insert(title, content, expires string) (int, error)
	Get(id int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}
