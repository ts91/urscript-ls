package server

import (
	"context"
	"strings"
	"sync"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/completion"
	"github.com/ts91/urscript-ls/internal/definition"
	"github.com/ts91/urscript-ls/internal/diagnostics"
	"github.com/ts91/urscript-ls/internal/docsymbols"
	"github.com/ts91/urscript-ls/internal/document"
	"github.com/ts91/urscript-ls/internal/hover"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/references"
	"github.com/ts91/urscript-ls/internal/signature"
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
			HoverProvider:          protocol.Boolean(true),
			DefinitionProvider:     protocol.Boolean(true),
			DocumentSymbolProvider: protocol.Boolean(true),
			ReferencesProvider:     protocol.Boolean(true),
			RenameProvider:         protocol.Boolean(true),
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{"."},
			},
			SignatureHelpProvider: &protocol.SignatureHelpOptions{
				TriggerCharacters:   []string{"(", ","},
				RetriggerCharacters: []string{","},
			},
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
	idx := s.store.GetIndex(params.TextDocument.URI)
	return hover.Provide(content, params.Position, s.registry, idx), nil
}

// Completion returns completion items for the current position.
func (s *Server) Completion(_ context.Context, params *protocol.CompletionParams) (protocol.CompletionResult, error) {
	idx := s.store.GetIndex(params.TextDocument.URI)
	items := completion.Provide(idx, s.registry)
	return items, nil
}

// SignatureHelp returns signature information for the function call at the cursor.
func (s *Server) SignatureHelp(_ context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	idx := s.store.GetIndex(params.TextDocument.URI)
	return signature.Provide(content, params.Position, s.registry, idx), nil
}

// DocumentSymbol returns the symbol outline for a document.
func (s *Server) DocumentSymbol(_ context.Context, params *protocol.DocumentSymbolParams) (protocol.DocumentSymbolResult, error) {
	idx := s.store.GetIndex(params.TextDocument.URI)
	return docsymbols.Provide(idx), nil
}

// Definition returns the definition location for the symbol at the cursor.
func (s *Server) Definition(_ context.Context, params *protocol.DefinitionParams) (protocol.DefinitionResult, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	idx := s.store.GetIndex(params.TextDocument.URI)
	return definition.Provide(content, params.Position, params.TextDocument.URI, idx), nil
}

// References returns all locations where the symbol at the cursor is used.
func (s *Server) References(_ context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	word := wordAt(content, params.Position)
	if word == "" {
		return nil, nil
	}
	positions := references.FindAllInSource(content, word)
	return references.ToLocations(positions, params.TextDocument.URI), nil
}

// Rename renames the symbol at the cursor across the document.
func (s *Server) Rename(_ context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	word := wordAt(content, params.Position)
	if word == "" {
		return nil, nil
	}
	positions := references.FindAllInSource(content, word)
	if len(positions) == 0 {
		return nil, nil
	}
	edits := references.ToTextEdits(positions, params.NewName)
	return &protocol.WorkspaceEdit{
		Changes: map[uri.URI][]protocol.TextEdit{
			params.TextDocument.URI: edits,
		},
	}, nil
}

// wordAt extracts the identifier at the given position.
func wordAt(content string, pos protocol.Position) string {
	lines := strings.Split(content, "\n")
	line := int(pos.Line)
	col := int(pos.Character)
	if line >= len(lines) {
		return ""
	}
	text := lines[line]
	if col >= len(text) {
		return ""
	}
	if !isIdentRune(rune(text[col])) {
		return ""
	}
	start := col
	for start > 0 && isIdentRune(rune(text[start-1])) {
		start--
	}
	end := col
	for end < len(text) && isIdentRune(rune(text[end])) {
		end++
	}
	return text[start:end]
}

func isIdentRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
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

	// Semantic diagnostics (warnings).
	diags = append(diags, diagnostics.Analyze(prog, idx, s.registry)...)

	c := s.getClient()
	if c == nil {
		return nil
	}
	return c.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         docURI,
		Diagnostics: diags,
	})
}
