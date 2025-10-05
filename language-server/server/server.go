package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/uplang/go"
	"go.lsp.dev/protocol"
)

// Server implements the UP Language Server
type Server struct {
	logger *slog.Logger
	parser *up.Parser

	// Document cache
	documents map[protocol.DocumentURI]*Document

	// Server capabilities
	capabilities protocol.ServerCapabilities
}

// Document represents a cached document
type Document struct {
	URI     protocol.DocumentURI
	Version int32
	Content string
	Parsed  *up.Document
	Errors  []protocol.Diagnostic
}

// NewServer creates a new UP language server
func NewServer(logger *slog.Logger) *Server {
	return &Server{
		logger:    logger,
		parser:    up.NewParser(),
		documents: make(map[protocol.DocumentURI]*Document),
		capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
			},
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{"!", "@", "."},
			},
			HoverProvider: true,
			DocumentFormattingProvider: true,
			DocumentSymbolProvider: true,
			DefinitionProvider: true,
		},
	}
}

// Run starts the language server
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("UP Language Server started")

	// Read from stdin, write to stdout (LSP protocol)
	return s.serve(ctx, os.Stdin, os.Stdout)
}

func (s *Server) serve(ctx context.Context, in io.Reader, out io.Writer) error {
	// This would use the LSP JSON-RPC handler
	// For now, simplified implementation
	s.logger.Info("Server listening on stdio")

	// In a real implementation, this would:
	// 1. Read JSON-RPC messages from stdin
	// 2. Route to appropriate handlers
	// 3. Write JSON-RPC responses to stdout

	// Placeholder for actual LSP message loop
	<-ctx.Done()
	return nil
}

// HandleInitialize handles the initialize request
func (s *Server) HandleInitialize(params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	s.logger.Info("Client initializing",
		slog.String("client", params.ClientInfo.Name),
		slog.String("version", params.ClientInfo.Version),
	)

	return &protocol.InitializeResult{
		Capabilities: s.capabilities,
		ServerInfo: &protocol.ServerInfo{
			Name:    "up-language-server",
			Version: "1.0.0",
		},
	}, nil
}

// HandleDidOpen handles document open notifications
func (s *Server) HandleDidOpen(params *protocol.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI
	s.logger.Debug("Document opened", slog.String("uri", string(uri)))

	doc := &Document{
		URI:     uri,
		Version: params.TextDocument.Version,
		Content: params.TextDocument.Text,
	}

	// Parse the document
	s.parseDocument(doc)

	s.documents[uri] = doc

	// Send diagnostics
	return s.publishDiagnostics(doc)
}

// HandleDidChange handles document change notifications
func (s *Server) HandleDidChange(params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI
	doc, exists := s.documents[uri]
	if !exists {
		return fmt.Errorf("document not found: %s", uri)
	}

	s.logger.Debug("Document changed", slog.String("uri", string(uri)))

	// Update content (full sync)
	if len(params.ContentChanges) > 0 {
		doc.Content = params.ContentChanges[0].Text
		doc.Version = params.TextDocument.Version
	}

	// Reparse
	s.parseDocument(doc)

	// Send diagnostics
	return s.publishDiagnostics(doc)
}

// HandleDidClose handles document close notifications
func (s *Server) HandleDidClose(params *protocol.DidCloseTextDocumentParams) error {
	uri := params.TextDocument.URI
	s.logger.Debug("Document closed", slog.String("uri", string(uri)))

	delete(s.documents, uri)
	return nil
}

// HandleCompletion provides completions
func (s *Server) HandleCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	uri := params.TextDocument.URI
	doc, exists := s.documents[uri]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", uri)
	}

	s.logger.Debug("Completion requested",
		slog.String("uri", string(uri)),
		slog.Int("line", int(params.Position.Line)),
		slog.Int("char", int(params.Position.Character)),
	)

	items := s.getCompletions(doc, params.Position)

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

// HandleHover provides hover information
func (s *Server) HandleHover(params *protocol.HoverParams) (*protocol.Hover, error) {
	uri := params.TextDocument.URI
	doc, exists := s.documents[uri]
	if !exists {
		return nil, nil
	}

	s.logger.Debug("Hover requested",
		slog.String("uri", string(uri)),
		slog.Int("line", int(params.Position.Line)),
	)

	return s.getHover(doc, params.Position), nil
}

// parseDocument parses a document and updates diagnostics
func (s *Server) parseDocument(doc *Document) {
	doc.Errors = nil

	parsed, err := s.parser.ParseDocument(strings.NewReader(doc.Content))
	if err != nil {
		// Parse error - create diagnostic
		doc.Errors = append(doc.Errors, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 100},
			},
			Severity: protocol.DiagnosticSeverityError,
			Source:   "up-parser",
			Message:  err.Error(),
		})
		return
	}

	doc.Parsed = parsed
}

// publishDiagnostics sends diagnostics to the client
func (s *Server) publishDiagnostics(doc *Document) error {
	// In real implementation, this would send via LSP protocol
	s.logger.Debug("Publishing diagnostics",
		slog.String("uri", string(doc.URI)),
		slog.Int("count", len(doc.Errors)),
	)
	return nil
}

// getCompletions generates completion items
func (s *Server) getCompletions(doc *Document, pos protocol.Position) []protocol.CompletionItem {
	items := []protocol.CompletionItem{
		{
			Label:  "!int",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "Integer type annotation",
		},
		{
			Label:  "!bool",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "Boolean type annotation",
		},
		{
			Label:  "!string",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "String type annotation",
		},
		{
			Label:  "!dur",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "Duration type annotation",
		},
		{
			Label:  "!ts",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "Timestamp type annotation",
		},
		{
			Label:  "!uuid",
			Kind:   protocol.CompletionItemKindKeyword,
			Detail: "UUID type annotation",
		},
	}

	// Add namespace completions
	namespaceItems := []protocol.CompletionItem{
		{
			Label:  "@string.",
			Kind:   protocol.CompletionItemKindModule,
			Detail: "String manipulation functions",
		},
		{
			Label:  "@math.",
			Kind:   protocol.CompletionItemKindModule,
			Detail: "Math operations",
		},
		{
			Label:  "@time.",
			Kind:   protocol.CompletionItemKindModule,
			Detail: "Time and date functions",
		},
		{
			Label:  "@random.",
			Kind:   protocol.CompletionItemKindModule,
			Detail: "Random generation",
		},
	}

	items = append(items, namespaceItems...)
	return items
}

// getHover provides hover information
func (s *Server) getHover(doc *Document, pos protocol.Position) *protocol.Hover {
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: "**UP Document**\n\nHover information for UP files.",
		},
	}
}

