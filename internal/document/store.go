package document

import (
	"path/filepath"
	"sync"

	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/symbols"
)

// NormalizeURI returns a canonical file URI for map keys and cross-file lookup.
func NormalizeURI(u uri.URI) uri.URI {
	if u == "" {
		return u
	}
	path := filepath.Clean(u.Path())
	if path == "" {
		return u
	}
	return uri.File(path)
}

// entry holds a document's content and its symbol index.
type entry struct {
	content     string
	index       *symbols.Index
	editorOpen  bool           // true if currently open in the editor
	clientURI   uri.URI        // URI the editor used (for PublishDiagnostics)
	diskContent string         // content from workspace scan (preserved across editor close)
	diskIndex   *symbols.Index // index from workspace scan
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

// Open stores a newly opened document (from the editor).
func (s *Store) Open(u uri.URI, content string) {
	clientURI := u
	u = NormalizeURI(u)
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.docs[u]; ok {
		// File was already indexed from disk; overlay with editor content.
		e.content = content
		e.index = nil
		e.editorOpen = true
		e.clientURI = clientURI
	} else {
		s.docs[u] = &entry{content: content, editorOpen: true, clientURI: clientURI}
	}
}

// OpenFromDisk stores a document discovered by the workspace scanner.
func (s *Store) OpenFromDisk(u uri.URI, content string) {
	u = NormalizeURI(u)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.docs[u]; ok {
		return // already open in editor, don't overwrite
	}
	s.docs[u] = &entry{
		content:     content,
		diskContent: content,
	}
}

// Update replaces the content of an open document.
func (s *Store) Update(u uri.URI, content string) {
	u = NormalizeURI(u)
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.docs[u]; ok {
		e.content = content
		e.index = nil // invalidate until next parse
	} else {
		s.docs[u] = &entry{content: content}
	}
}

// Close handles the editor closing a document. If the file was loaded
// from disk (workspace scan), it reverts to the disk content and index
// rather than being deleted.
func (s *Store) Close(u uri.URI) {
	u = NormalizeURI(u)
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.docs[u]
	if !ok {
		return
	}
	if e.diskContent != "" {
		// Revert to disk-scanned state.
		e.content = e.diskContent
		e.index = e.diskIndex
		e.editorOpen = false
	} else {
		delete(s.docs, u)
	}
}

// Get returns the content of a document and whether it exists.
func (s *Store) Get(u uri.URI) (string, bool) {
	u = NormalizeURI(u)
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
	u = NormalizeURI(u)
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.docs[u]; ok {
		e.index = idx
		// If not editor-open, this is a disk scan index — save it.
		if !e.editorOpen {
			e.diskIndex = idx
		}
	}
}

// GetIndex returns the symbol index for a document, if available.
func (s *Store) GetIndex(u uri.URI) *symbols.Index {
	u = NormalizeURI(u)
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

// AllContents returns the content of every document in the store.
func (s *Store) AllContents() map[uri.URI]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[uri.URI]string, len(s.docs))
	for u, e := range s.docs {
		result[u] = e.content
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

// EditorOpenURIs returns URIs of documents currently open in the editor.
func (s *Store) EditorOpenURIs() []uri.URI {
	s.mu.RLock()
	defer s.mu.RUnlock()
	uris := make([]uri.URI, 0, len(s.docs))
	for u, e := range s.docs {
		if e.editorOpen {
			uris = append(uris, u)
		}
	}
	return uris
}

// ClientURI returns the URI to use when publishing LSP notifications for a
// document (the editor's URI, not necessarily the normalized store key).
func (s *Store) ClientURI(u uri.URI) uri.URI {
	u = NormalizeURI(u)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.docs[u]; ok && e.clientURI != "" {
		return e.clientURI
	}
	return u
}
