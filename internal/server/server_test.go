package server

import (
	"context"
	"testing"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func setupServer(t *testing.T) (context.Context, protocol.Server, protocol.Client) {
	t.Helper()

	left, right := jsonrpc2.NewChannelStreamPair(16)
	srv := NewServer()

	srvCtx, _, client := protocol.NewServer(context.Background(), srv, left)
	srv.SetClient(client)

	clientCtx, _, serverProxy := protocol.NewClient(srvCtx, protocol.UnimplementedClient{}, right)
	_ = clientCtx

	return srvCtx, serverProxy, client
}

func TestInitialize(t *testing.T) {
	left, right := jsonrpc2.NewChannelStreamPair(16)
	srv := NewServer()

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
	if !ok || ver != "0.1.0" {
		t.Errorf("expected version '0.1.0', got %q (ok=%v)", ver, ok)
	}

	srvConn.Close()
	clientConn.Close()
}

func TestDidOpenPublishesDiagnostics(t *testing.T) {
	srv := NewServer()
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
	srv := NewServer()
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

func TestHoverUnknownSymbol(t *testing.T) {
	srv := NewServer()
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
