package server

import (
	"context"
	"sync"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/document"
	"github.com/ts91/urscript-ls/internal/hover"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

// Server implements protocol.Server for URScript.
type Server struct {
	protocol.UnimplementedServer

	mu       sync.Mutex
	client   protocol.Client
	store    *document.Store
	registry *builtins.Registry
}

// NewServer creates a new URScript language server.
func NewServer() *Server {
	return &Server{
		store:    document.NewStore(),
		registry: builtins.NewRegistry(),
	}
}

// SetClient stores the client proxy used to push notifications.
func (s *Server) SetClient(c protocol.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.client = c
}

func (s *Server) getClient() protocol.Client {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.client
}

// Initialize handles the initialize request.
func (s *Server) Initialize(_ context.Context, _ *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	openClose := true
	syncKind := protocol.TextDocumentSyncKindFull
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: &openClose,
				Change:    &syncKind,
			},
			HoverProvider: protocol.Boolean(true),
		},
		ServerInfo: protocol.ServerInfo{
			Name:    "urscript-ls",
			Version: protocol.NewOptional[string]("0.1.0"),
		},
	}, nil
}

// Initialized is called after the client confirms initialization.
func (s *Server) Initialized(_ context.Context, _ *protocol.InitializedParams) error {
	return nil
}

// Shutdown is called before Exit.
func (s *Server) Shutdown(_ context.Context) error {
	return nil
}

// Exit terminates the server.
func (s *Server) Exit(_ context.Context) error {
	return nil
}

// DidOpen handles opening a text document.
func (s *Server) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.store.Open(params.TextDocument.URI, params.TextDocument.Text)
	return s.publishDiagnostics(ctx, params.TextDocument.URI)
}

// DidChange handles text document changes (full sync).
func (s *Server) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if whole, ok := change.(*protocol.TextDocumentContentChangeWholeDocument); ok {
			s.store.Update(params.TextDocument.URI, whole.Text)
		}
	}
	return s.publishDiagnostics(ctx, params.TextDocument.URI)
}

// DidClose handles closing a text document.
func (s *Server) DidClose(_ context.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.store.Close(params.TextDocument.URI)
	return nil
}

// Hover returns hover information for the symbol at the given position.
func (s *Server) Hover(_ context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	return hover.Provide(content, params.Position, s.registry), nil
}

// publishDiagnostics parses the document and sends diagnostics to the client.
func (s *Server) publishDiagnostics(ctx context.Context, docURI uri.URI) error {
	content, ok := s.store.Get(docURI)
	if !ok {
		return nil
	}

	prog, errs := parser.Parse(content)
	idx := symbols.Build(prog, content)
	s.store.SetIndex(docURI, idx)

	diags := make([]protocol.Diagnostic, 0, len(errs))
	for _, e := range errs {
		diags = append(diags, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(e.Line), Character: uint32(e.Column)},
				End:   protocol.Position{Line: uint32(e.Line), Character: uint32(e.Column)},
			},
			Severity: protocol.DiagnosticSeverityError,
			Source:   protocol.NewOptional[string]("urscript-ls"),
			Message:  protocol.String(e.Message),
		})
	}

	c := s.getClient()
	if c == nil {
		return nil
	}
	return c.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         docURI,
		Diagnostics: diags,
	})
}
