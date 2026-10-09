package server

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/builtins"
	"github.com/ts91/urscript-ls/internal/codeaction"
	"github.com/ts91/urscript-ls/internal/completion"
	"github.com/ts91/urscript-ls/internal/definition"
	"github.com/ts91/urscript-ls/internal/diagnostics"
	"github.com/ts91/urscript-ls/internal/docsymbols"
	"github.com/ts91/urscript-ls/internal/document"
	"github.com/ts91/urscript-ls/internal/folding"
	"github.com/ts91/urscript-ls/internal/format"
	"github.com/ts91/urscript-ls/internal/highlight"
	"github.com/ts91/urscript-ls/internal/hover"
	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/references"
	"github.com/ts91/urscript-ls/internal/selrange"
	"github.com/ts91/urscript-ls/internal/semtokens"
	"github.com/ts91/urscript-ls/internal/signature"
	"github.com/ts91/urscript-ls/internal/symbols"
	"github.com/ts91/urscript-ls/internal/workspacesymbol"
)

// Server implements protocol.Server for URScript.
type Server struct {
	protocol.UnimplementedServer

	mu       sync.Mutex
	client   protocol.Client
	store    *document.Store
	registry *builtins.Registry
	roots    []string // workspace root filesystem paths
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
func (s *Server) Initialize(_ context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	// Capture workspace roots for background scanning.
	if folders, ok := params.WorkspaceFolders.Get(); ok {
		for _, f := range folders {
			s.roots = append(s.roots, f.URI.Path())
		}
	}
	if len(s.roots) == 0 && params.RootURI != nil {
		s.roots = append(s.roots, params.RootURI.Path())
	}

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
			ReferencesProvider:             protocol.Boolean(true),
			RenameProvider:                 protocol.Boolean(true),
			DocumentFormattingProvider:     protocol.Boolean(true),
			DocumentHighlightProvider:      protocol.Boolean(true),
			FoldingRangeProvider:           protocol.Boolean(true),
			SelectionRangeProvider:         protocol.Boolean(true),
			SemanticTokensProvider: &protocol.SemanticTokensOptions{
				Legend: semtokens.Legend(),
				Full:   protocol.Boolean(true),
			},
			WorkspaceSymbolProvider:        protocol.Boolean(true),
			CodeActionProvider:             protocol.Boolean(true),
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
	go s.scanWorkspace()
	return nil
}

// scanWorkspace reads all .script/.urscript files in the workspace roots
// and indexes them so cross-file features work immediately.
func (s *Server) scanWorkspace() {
	for _, root := range s.roots {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			if ext != ".script" && ext != ".urscript" && ext != ".variables" {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			docURI := uri.File(path)
			content := string(data)

			// Only load if not already open in the editor.
			if _, ok := s.store.Get(docURI); !ok {
				s.store.Open(docURI, content)
				prog, _ := parser.Parse(content)
				idx := symbols.Build(prog, content)
				s.store.SetIndex(docURI, idx)
			}
			return nil
		})
	}
	log.Printf("workspace scan complete: %d files indexed", len(s.store.AllURIs()))
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
	idx := s.mergedIndex(params.TextDocument.URI)
	return hover.Provide(content, params.Position, s.registry, idx), nil
}

// Completion returns completion items for the current position.
func (s *Server) Completion(_ context.Context, params *protocol.CompletionParams) (protocol.CompletionResult, error) {
	idx := s.mergedIndex(params.TextDocument.URI)
	items := completion.Provide(idx, s.registry)
	return items, nil
}

// SignatureHelp returns signature information for the function call at the cursor.
func (s *Server) SignatureHelp(_ context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	idx := s.mergedIndex(params.TextDocument.URI)
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
	// Try current file first.
	idx := s.store.GetIndex(params.TextDocument.URI)
	result := definition.Provide(content, params.Position, params.TextDocument.URI, idx)
	if len(result) > 0 {
		return result, nil
	}
	// Search other workspace files.
	word := wordAt(content, params.Position)
	if word == "" {
		return nil, nil
	}
	for docURI, fileIdx := range s.store.AllIndices() {
		if docURI == params.TextDocument.URI {
			continue
		}
		for _, sym := range fileIdx.Lookup(word) {
			if sym.Kind == symbols.KindFunction || sym.Kind == symbols.KindThread ||
				sym.Kind == symbols.KindGlobal {
				return protocol.LocationSlice{{
					URI: docURI,
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(sym.Pos.Line), Character: uint32(sym.Pos.Column)},
						End:   protocol.Position{Line: uint32(sym.Pos.Line), Character: uint32(sym.Pos.Column + len(sym.Name))},
					},
				}}, nil
			}
		}
	}
	return nil, nil
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
	// Search all workspace files.
	var allLocs []protocol.Location
	for _, docURI := range s.store.AllURIs() {
		src, ok := s.store.Get(docURI)
		if !ok {
			continue
		}
		positions := references.FindAllInSource(src, word)
		allLocs = append(allLocs, references.ToLocations(positions, docURI)...)
	}
	_ = content // already used for wordAt
	return allLocs, nil
}

// Rename renames the symbol at the cursor across all workspace files.
func (s *Server) Rename(_ context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	word := wordAt(content, params.Position)
	if word == "" {
		return nil, nil
	}
	changes := make(map[uri.URI][]protocol.TextEdit)
	for _, docURI := range s.store.AllURIs() {
		src, ok := s.store.Get(docURI)
		if !ok {
			continue
		}
		positions := references.FindAllInSource(src, word)
		if len(positions) > 0 {
			changes[docURI] = references.ToTextEdits(positions, params.NewName)
		}
	}
	if len(changes) == 0 {
		return nil, nil
	}
	return &protocol.WorkspaceEdit{Changes: changes}, nil
}

// DocumentHighlight highlights all occurrences of the word under cursor.
func (s *Server) DocumentHighlight(_ context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	word := wordAt(content, params.Position)
	return highlight.Provide(content, params.Position, word), nil
}

// Symbols returns workspace symbols matching a query.
func (s *Server) Symbols(_ context.Context, params *protocol.WorkspaceSymbolParams) (protocol.WorkspaceSymbolResult, error) {
	indices := s.store.AllIndices()
	return workspacesymbol.Provide(params.Query, indices), nil
}

// CodeAction returns quick fixes for diagnostics.
func (s *Server) CodeAction(_ context.Context, params *protocol.CodeActionParams) ([]protocol.CommandOrCodeAction, error) {
	idx := s.mergedIndex(params.TextDocument.URI)
	return codeaction.Provide(params, s.registry, idx), nil
}

// SemanticTokensFull returns semantic tokens for the entire document.
func (s *Server) SemanticTokensFull(_ context.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	idx := s.mergedIndex(params.TextDocument.URI)
	return semtokens.Provide(content, s.registry, idx), nil
}

// SelectionRange returns nested selection ranges for smart expand/shrink.
func (s *Server) SelectionRange(_ context.Context, params *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	prog, _ := parser.Parse(content)
	return selrange.Provide(prog, content, params.Positions), nil
}

// FoldingRanges returns folding ranges for blocks and comment runs.
func (s *Server) FoldingRanges(_ context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	prog, _ := parser.Parse(content)
	return folding.Provide(prog, content), nil
}

// Formatting reformats the entire document.
func (s *Server) Formatting(_ context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	content, ok := s.store.Get(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	return format.Format(content, params.Options), nil
}

// mergedIndex returns a symbol index combining the current file's index
// with indices from all other workspace files.
func (s *Server) mergedIndex(docURI uri.URI) *symbols.Index {
	primary := s.store.GetIndex(docURI)
	allIdx := s.store.AllIndices()
	var others []*symbols.Index
	for u, idx := range allIdx {
		if u != docURI {
			others = append(others, idx)
		}
	}
	if len(others) == 0 {
		return primary
	}
	return symbols.Merge(primary, others...)
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

	// Semantic diagnostics use merged index so cross-file symbols are known.
	merged := s.mergedIndex(docURI)
	diags = append(diags, diagnostics.Analyze(prog, merged, s.registry)...)

	c := s.getClient()
	if c == nil {
		return nil
	}
	return c.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         docURI,
		Diagnostics: diags,
	})
}
