package document

import (
	"sync"

	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/symbols"
)

// entry holds a document's content and its symbol index.
type entry struct {
	content string
	index   *symbols.Index
}

// Store holds the contents of open text documents.
type Store struct {
	mu   sync.RWMutex
	docs map[uri.URI]*entry
}

// NewStore creates a new document store.
func NewStore() *Store {
	return &Store{
		docs: make(map[uri.URI]*entry),
	}
}

// Open stores a newly opened document.
func (s *Store) Open(u uri.URI, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs[u] = &entry{content: content}
}

// Update replaces the content of an open document.
func (s *Store) Update(u uri.URI, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.docs[u]; ok {
		e.content = content
		e.index = nil // invalidate until next parse
	} else {
		s.docs[u] = &entry{content: content}
	}
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
	e, ok := s.docs[u]
	if !ok {
		return "", false
	}
	return e.content, true
}

// SetIndex stores the symbol index for a document.
func (s *Store) SetIndex(u uri.URI, idx *symbols.Index) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.docs[u]; ok {
		e.index = idx
	}
}

// GetIndex returns the symbol index for a document, if available.
func (s *Store) GetIndex(u uri.URI) *symbols.Index {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.docs[u]; ok {
		return e.index
	}
	return nil
}

// AllIndices returns the symbol index from every document in the store.
func (s *Store) AllIndices() map[uri.URI]*symbols.Index {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[uri.URI]*symbols.Index, len(s.docs))
	for u, e := range s.docs {
		if e.index != nil {
			result[u] = e.index
		}
	}
	return result
}

// AllURIs returns the URIs of all documents in the store.
func (s *Store) AllURIs() []uri.URI {
	s.mu.RLock()
	defer s.mu.RUnlock()
	uris := make([]uri.URI, 0, len(s.docs))
	for u := range s.docs {
		uris = append(uris, u)
	}
	return uris
}
