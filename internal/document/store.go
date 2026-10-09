package document

import (
	"sync"

	"go.lsp.dev/uri"
)

// Store holds the contents of open text documents.
type Store struct {
	mu   sync.RWMutex
	docs map[uri.URI]string
}

// NewStore creates a new document store.
func NewStore() *Store {
	return &Store{
		docs: make(map[uri.URI]string),
	}
}

// Open stores a newly opened document.
func (s *Store) Open(u uri.URI, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs[u] = content
}

// Update replaces the content of an open document.
func (s *Store) Update(u uri.URI, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs[u] = content
}

// Close removes a document from the store.
func (s *Store) Close(u uri.URI) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.docs, u)
}

// Get returns the content of a document and whether it exists.
func (s *Store) Get(u uri.URI) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	content, ok := s.docs[u]
	return content, ok
}
