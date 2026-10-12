package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"

	"github.com/ts91/urscript-ls/internal/server"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println("urscript-ls", version)
		os.Exit(0)
	}

	// Capture stdout for the LSP stream, then redirect the global
	// os.Stdout to stderr so any stray writes (e.g. from libraries)
	// cannot corrupt the JSON-RPC framing.
	lspOut := os.Stdout
	os.Stdout = os.Stderr
	log.SetOutput(os.Stderr)

	ctx := context.Background()

	stream := jsonrpc2.NewStream(rwc{os.Stdin, lspOut})
	srv := server.NewServer(version)

	ctx, conn, client := protocol.NewServer(ctx, srv, stream)
	srv.SetClient(client)

	<-conn.Done()
	if err := conn.Err(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// rwc combines separate reader and writer into io.ReadWriteCloser.
type rwc struct {
	r *os.File
	w *os.File
}

func (c rwc) Read(p []byte) (int, error)  { return c.r.Read(p) }
func (c rwc) Write(p []byte) (int, error) { return c.w.Write(p) }
func (c rwc) Close() error {
	c.r.Close()
	return c.w.Close()
}
