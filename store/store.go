package store

import (
	"github.com/mgeri/snippetbox/pkg/models"
)

// SnippetStore is the persistent store of Snippet
type SnippetStore interface {
	Insert(title, content, expires string) (int, error)
	Get(id int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

// UserStore is the persistent store of User
type UserStore interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Get(id int) (*models.User, error)
}
