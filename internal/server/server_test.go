package server

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/ts91/urscript-ls/internal/parser"
	"github.com/ts91/urscript-ls/internal/symbols"
)

func setupServer(t *testing.T) (context.Context, protocol.Server, protocol.Client) {
	t.Helper()

	left, right := jsonrpc2.NewChannelStreamPair(16)
	srv := NewServer("test")

	srvCtx, _, client := protocol.NewServer(context.Background(), srv, left)
	srv.SetClient(client)

	clientCtx, _, serverProxy := protocol.NewClient(srvCtx, protocol.UnimplementedClient{}, right)
	_ = clientCtx

	return srvCtx, serverProxy, client
}

func TestInitialize(t *testing.T) {
	left, right := jsonrpc2.NewChannelStreamPair(16)
	srv := NewServer("test")

	ctx := context.Background()
	ctx, srvConn, client := protocol.NewServer(ctx, srv, left)
	srv.SetClient(client)
	srvConn.Go(ctx, protocol.ServerHandler(srv, nil))

	_, clientConn, serverProxy := protocol.NewClient(ctx, protocol.UnimplementedClient{}, right)
	clientConn.Go(ctx, protocol.ClientHandler(protocol.UnimplementedClient{}, nil))

	result, err := serverProxy.Initialize(ctx, &protocol.InitializeParams{})
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	info := result.ServerInfo
	if info.Name != "urscript-ls" {
		t.Errorf("expected server name 'urscript-ls', got %q", info.Name)
	}

	ver, ok := info.Version.Get()
	if !ok || ver != "test" {
		t.Errorf("expected version 'test', got %q (ok=%v)", ver, ok)
	}

	srvConn.Close()
	clientConn.Close()
}

func TestPublishDiagnosticsUsesClientURI(t *testing.T) {
	srv := NewServer("test")
	ctx := context.Background()

	diags := make(chan *protocol.PublishDiagnosticsParams, 1)
	srv.SetClient(&mockClient{diagsCh: diags})

	// Editor URI with a distinct string form from the normalized store key.
	clientURI := uri.URI("file:///tmp/client-alias.urscript")
	storeKey := uri.File("/tmp/client-alias.urscript")
	srv.store.Open(clientURI, "unknown_fn()\n")

	if err := srv.publishDiagnostics(ctx, storeKey); err != nil {
		t.Fatalf("publishDiagnostics: %v", err)
	}
	params := <-diags
	if params.URI != clientURI {
		t.Errorf("expected publish URI %q, got %q", clientURI, params.URI)
	}
}

func TestDidOpenPublishesDiagnostics(t *testing.T) {
	srv := NewServer("test")
	ctx := context.Background()

	diags := make(chan *protocol.PublishDiagnosticsParams, 1)
	mc := &mockClient{diagsCh: diags}

	left, right := jsonrpc2.NewChannelStreamPair(16)
	ctx, srvConn, _ := protocol.NewServer(ctx, srv, left)
	srv.SetClient(mc)
	srvConn.Go(ctx, protocol.ServerHandler(srv, nil))

	_, clientConn, serverProxy := protocol.NewClient(ctx, mc, right)
	clientConn.Go(ctx, protocol.ClientHandler(mc, nil))

	_, err := serverProxy.Initialize(ctx, &protocol.InitializeParams{})
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	testURI := uri.File("/tmp/test.urscript")
	err = serverProxy.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        testURI,
			LanguageID: "urscript",
			Version:    1,
			Text:       "def broken(:\nend",
		},
	})
	if err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	params := <-diags
	if params.URI != testURI {
		t.Errorf("expected URI %q, got %q", testURI, params.URI)
	}
	if len(params.Diagnostics) == 0 {
		t.Error("expected diagnostics for broken program, got none")
	}

	srvConn.Close()
	clientConn.Close()
}

func TestHover(t *testing.T) {
	srv := NewServer("test")
	ctx := context.Background()

	testURI := uri.File("/tmp/test.urscript")
	srv.store.Open(testURI, "movej(q)")

	result, err := srv.Hover(ctx, &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
			Position:     protocol.Position{Line: 0, Character: 2},
		},
	})
	if err != nil {
		t.Fatalf("Hover failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected hover result, got nil")
	}
	mc, ok := result.Contents.(*protocol.MarkupContent)
	if !ok {
		t.Fatalf("expected *MarkupContent, got %T", result.Contents)
	}
	if mc.Kind != protocol.MarkupKindMarkdown {
		t.Errorf("expected markdown, got %v", mc.Kind)
	}
}

func TestScanWorkspaceRepublishExamples(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "examples"))
	if err != nil {
		t.Fatal(err)
	}
	mainPath := filepath.Join(root, "main.script")
	mainData, err := os.ReadFile(mainPath)
	if err != nil {
		t.Skipf("examples/main.script: %v", err)
	}

	srv := NewServer("test")
	srv.roots = []string{root}
	ctx := context.Background()

	diags := make(chan *protocol.PublishDiagnosticsParams, 16)
	srv.SetClient(&mockClient{diagsCh: diags})

	mainURI := uri.File(mainPath)
	srv.store.Open(mainURI, string(mainData))

	if err := srv.publishDiagnostics(ctx, mainURI); err != nil {
		t.Fatalf("publishDiagnostics: %v", err)
	}
	params := <-diags
	if !diagIsNotDefined(params.Diagnostics, "compute_offset") {
		t.Fatal("expected stale 'not defined' for compute_offset before scan")
	}

	srv.scanWorkspace()

	for {
		params = <-diags
		if params.URI != mainURI {
			continue
		}
		if diagIsNotDefined(params.Diagnostics, "compute_offset") {
			t.Error("expected compute_offset 'not defined' cleared after workspace scan")
		}
		break
	}
}

func TestRepublishDiagnosticsAfterCrossFileIndex(t *testing.T) {
	srv := NewServer("test")
	ctx := context.Background()

	diags := make(chan *protocol.PublishDiagnosticsParams, 4)
	mc := &mockClient{diagsCh: diags}
	srv.SetClient(mc)

	mainURI := uri.File("/tmp/main.urscript")
	utilsURI := uri.File("/tmp/utils.urscript")

	srv.store.Open(mainURI, "other_func()\n")

	if err := srv.publishDiagnostics(ctx, mainURI); err != nil {
		t.Fatalf("publishDiagnostics: %v", err)
	}
	params := <-diags
	if !diagContains(params.Diagnostics, "not defined") {
		t.Fatal("expected 'not defined' before cross-file index is available")
	}

	utilsSrc := "def other_func():\nend\n"
	srv.store.OpenFromDisk(utilsURI, utilsSrc)
	utilsProg, _ := parser.Parse(utilsSrc)
	srv.store.SetIndex(utilsURI, symbols.Build(utilsProg, utilsSrc))

	if err := srv.republishOpenEditorDiagnostics(ctx); err != nil {
		t.Fatalf("republishOpenEditorDiagnostics: %v", err)
	}

	for {
		params = <-diags
		if params.URI != mainURI {
			continue
		}
		if diagContains(params.Diagnostics, "not defined") {
			t.Errorf("expected no 'not defined' after republish, got: %v", params.Diagnostics)
		}
		break
	}
}

func diagContains(diags []protocol.Diagnostic, substr string) bool {
	for _, d := range diags {
		if s, ok := d.Message.(protocol.String); ok && strings.Contains(string(s), substr) {
			return true
		}
	}
	return false
}

func diagIsNotDefined(diags []protocol.Diagnostic, name string) bool {
	want := "'" + name + "' is not defined"
	for _, d := range diags {
		if s, ok := d.Message.(protocol.String); ok && string(s) == want {
			return true
		}
	}
	return false
}

func TestHoverUnknownSymbol(t *testing.T) {
	srv := NewServer("test")
	ctx := context.Background()

	testURI := uri.File("/tmp/test.urscript")
	srv.store.Open(testURI, "myCustomFunc()")

	result, err := srv.Hover(ctx, &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
			Position:     protocol.Position{Line: 0, Character: 2},
		},
	})
	if err != nil {
		t.Fatalf("Hover failed: %v", err)
	}
	if result != nil {
		t.Error("expected nil hover for unknown symbol")
	}
}

// mockClient captures PublishDiagnostics calls.
type mockClient struct {
	protocol.UnimplementedClient
	diagsCh chan<- *protocol.PublishDiagnosticsParams
}

func (m *mockClient) PublishDiagnostics(_ context.Context, params *protocol.PublishDiagnosticsParams) error {
	select {
	case m.diagsCh <- params:
	default:
	}
	return nil
}
