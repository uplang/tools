// Package main provides the UP CLI application.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
	up "github.com/uplang/go"
)

var (
	// version is set at build time via ldflags
	version = "dev"
	// commit is set at build time via ldflags
	commit = "none"
	// date is set at build time via ldflags
	date = "unknown"
	// builtBy is set at build time via ldflags
	builtBy = "unknown"
)

// App provides the CLI application functionality.
type App struct {
	parser   *up.Parser
	output   io.Writer
	input    io.Reader
	exitFunc func(int)
}

// NewApp creates a new CLI application with dependency injection.
func NewApp(p *up.Parser, out io.Writer, in io.Reader, exitFunc func(int)) *App {
	return &App{
		parser:   p,
		output:   out,
		input:    in,
		exitFunc: exitFunc,
	}
}

// DefaultApp creates a new CLI application with default dependencies.
func DefaultApp() *App {
	return NewApp(
		up.NewParser(),
		os.Stdout,
		os.Stdin,
		os.Exit,
	)
}

// Run executes the CLI application.
func (a *App) Run(args []string) error {
	app := &cli.App{
		Name:    "up",
		Usage:   "UP is a tool for managing UP (Unified Properties) documents",
		Version: version,
		UsageText: `up <command> [arguments]

The commands are:
    parse       parse UP documents and output as JSON
    format      format UP documents (alias: fmt)
    validate    validate UP documents against schemas (alias: vet)
    eval        evaluate dynamic namespaces
    convert     convert between UP and other formats
    lsp         start the UP language server
    repl        start interactive REPL
    tool        run specified UP tool
    version     print UP version

Use "up <command> -h" for more information about a command.`,
		Commands: []*cli.Command{
			a.parseCommand(),
			a.formatCommand(),
			a.validateCommand(),
			a.evalCommand(),
			a.convertCommand(),
			a.templateCommand(),
			a.lspCommand(),
			a.replCommand(),
			a.toolCommand(),
			a.versionCommand(),
		},
		Before: func(c *cli.Context) error {
			return nil
		},
	}

	return app.Run(args)
}

// parseCommand creates the parse command.
func (a *App) parseCommand() *cli.Command {
	return &cli.Command{
		Name:    "parse",
		Aliases: []string{"p"},
		Usage:   "Parse a UP document and output as JSON",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input file (default: stdin)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file (default: stdout)",
			},
			&cli.BoolFlag{
				Name:  "pretty",
				Usage: "Pretty print JSON output",
			},
		},
		Action: a.handleParse,
	}
}

// formatCommand creates the format command.
func (a *App) formatCommand() *cli.Command {
	return &cli.Command{
		Name:    "format",
		Aliases: []string{"fmt", "f"},
		Usage:   "Format and validate a UP document",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input file (default: stdin)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file (default: stdout)",
			},
		},
		Action: a.handleFormat,
	}
}

// validateCommand creates the validate command.
func (a *App) validateCommand() *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"vet", "v"},
		Usage:   "Validate a UP document syntax",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input file (default: stdin)",
			},
		},
		Action: a.handleValidate,
	}
}

// evalCommand creates the eval command.
func (a *App) evalCommand() *cli.Command {
	return &cli.Command{
		Name:    "eval",
		Aliases: []string{"e"},
		Usage:   "Evaluate dynamic namespaces in a UP document",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input file (default: stdin)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file (default: stdout)",
			},
			&cli.StringFlag{
				Name:  "ns-path",
				Usage: "Namespace search path",
				Value: "./up-namespaces",
			},
			&cli.BoolFlag{
				Name:  "pretty",
				Usage: "Pretty print output",
			},
		},
		Action: a.handleEval,
	}
}

// convertCommand creates the convert command.
func (a *App) convertCommand() *cli.Command {
	return &cli.Command{
		Name:    "convert",
		Aliases: []string{"c"},
		Usage:   "Convert between UP and other formats",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Input file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output file",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "from",
				Usage: "Input format (json, yaml, toml) - auto-detected if not specified",
			},
			&cli.StringFlag{
				Name:     "to",
				Usage:    "Output format (up, json, yaml, toml)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "pretty",
				Usage: "Pretty print output",
			},
		},
		Action: a.handleConvert,
	}
}

// lspCommand creates the lsp command.
func (a *App) lspCommand() *cli.Command {
	return &cli.Command{
		Name:    "lsp",
		Aliases: []string{"language-server"},
		Usage:   "Start the UP language server",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug logging",
			},
			&cli.StringFlag{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "Log file path (default: stderr)",
			},
		},
		Action: a.handleLSP,
	}
}

// replCommand creates the repl command.
func (a *App) replCommand() *cli.Command {
	return &cli.Command{
		Name:  "repl",
		Usage: "Start interactive REPL",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug output",
			},
		},
		Action: a.handleREPL,
	}
}

// toolCommand creates the tool command.
func (a *App) toolCommand() *cli.Command {
	return &cli.Command{
		Name:      "tool",
		Usage:     "Run specified UP tool",
		UsageText: "up tool <name> [arguments]",
		Action:    a.handleTool,
	}
}

// versionCommand creates the version command.
func (a *App) versionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print version information",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "short",
				Aliases: []string{"s"},
				Usage:   "Print only the version number",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Print version information as JSON",
			},
		},
		Action: a.handleVersion,
	}
}

// handleParse processes the parse command.
func (a *App) handleParse(c *cli.Context) error {
	input, err := a.getInput(c.String("input"))
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	defer a.closeIfFile(input)

	output, err := a.getOutput(c.String("output"))
	if err != nil {
		return fmt.Errorf("failed to open output: %w", err)
	}
	defer a.closeIfFile(output)

	doc, err := a.parser.ParseDocument(input)
	if err != nil {
		return fmt.Errorf("failed to parse document: %w", err)
	}

	return a.writeJSON(output, doc, c.Bool("pretty"))
}

// handleFormat processes the format command.
func (a *App) handleFormat(c *cli.Context) error {
	input, err := a.getInput(c.String("input"))
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	defer a.closeIfFile(input)

	output, err := a.getOutput(c.String("output"))
	if err != nil {
		return fmt.Errorf("failed to open output: %w", err)
	}
	defer a.closeIfFile(output)

	doc, err := a.parser.ParseDocument(input)
	if err != nil {
		return fmt.Errorf("failed to parse document: %w", err)
	}

	return a.writeUP(output, doc)
}

// handleValidate processes the validate command.
func (a *App) handleValidate(c *cli.Context) error {
	input, err := a.getInput(c.String("input"))
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	defer a.closeIfFile(input)

	_, err = a.parser.ParseDocument(input)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Fprintf(a.output, "✓ Document is valid\n")
	return nil
}

// getInput returns an io.ReadCloser for the input source.
func (a *App) getInput(filename string) (io.ReadCloser, error) {
	if filename == "" {
		return io.NopCloser(a.input), nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// getOutput returns an io.WriteCloser for the output destination.
func (a *App) getOutput(filename string) (io.WriteCloser, error) {
	if filename == "" {
		return &nopWriteCloser{a.output}, nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// closeIfFile closes the reader/writer if it's a file.
func (a *App) closeIfFile(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v", err)
		}
	}
}

// writeJSON writes the document as JSON.
func (a *App) writeJSON(w io.Writer, doc *up.Document, pretty bool) error {
	var data []byte
	var err error

	if pretty {
		data, err = json.MarshalIndent(doc, "", "  ")
	} else {
		data, err = json.Marshal(doc)
	}

	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// writeUP writes the document back as formatted UP.
func (a *App) writeUP(w io.Writer, doc *up.Document) error {
	for _, node := range doc.Nodes {
		line := a.formatNode(node)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// formatNode formats a single node back to UP format.
func (a *App) formatNode(node up.Node) string {
	key := node.Key
	if node.Type != "" {
		key += "!" + node.Type
	}

	switch v := node.Value.(type) {
	case string:
		if strings.Contains(v, "\n") {
			return fmt.Sprintf("%s ```\n%s\n```", key, v)
		}
		return fmt.Sprintf("%s %s", key, v)
	case up.Block:
		lines := []string{key + " {"}
		for k, val := range v {
			subNode := up.Node{Key: k, Value: val}
			lines = append(lines, "  "+a.formatNode(subNode))
		}
		lines = append(lines, "}")
		return strings.Join(lines, "\n")
	case up.List:
		if len(v) == 0 {
			return key + " []"
		}
		lines := []string{key + " ["}
		for _, item := range v {
			lines = append(lines, fmt.Sprintf("  %v", item))
		}
		lines = append(lines, "]")
		return strings.Join(lines, "\n")
	default:
		return fmt.Sprintf("%s %v", key, v)
	}
}

// templateCommand creates the template command.
func (a *App) templateCommand() *cli.Command {
	return &cli.Command{
		Name:    "template",
		Aliases: []string{"t"},
		Usage:   "Process UP templates with overlays and composition",
		Subcommands: []*cli.Command{
			{
				Name:  "process",
				Usage: "Process a template and output the result",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "Input template file",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Output file (default: stdout)",
					},
					&cli.BoolFlag{
						Name:  "json",
						Usage: "Output as JSON instead of UP",
					},
					&cli.BoolFlag{
						Name:  "pretty",
						Usage: "Pretty print output",
					},
				},
				Action: a.handleTemplateProcess,
			},
			{
				Name:  "validate",
				Usage: "Validate a template",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "Input template file",
						Required: true,
					},
				},
				Action: a.handleTemplateValidate,
			},
		},
	}
}

// handleTemplateProcess processes a template
func (a *App) handleTemplateProcess(c *cli.Context) error {
	engine := up.NewTemplateEngine()
	doc, err := engine.ProcessTemplate(c.String("input"))
	if err != nil {
		return fmt.Errorf("template processing failed: %w", err)
	}

	output, err := a.getOutput(c.String("output"))
	if err != nil {
		return fmt.Errorf("failed to open output: %w", err)
	}
	defer a.closeIfFile(output)

	if c.Bool("json") {
		return a.writeJSON(output, doc, c.Bool("pretty"))
	}
	return a.writeUP(output, doc)
}

// handleTemplateValidate validates a template
func (a *App) handleTemplateValidate(c *cli.Context) error {
	engine := up.NewTemplateEngine()
	_, err := engine.ProcessTemplate(c.String("input"))
	if err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	fmt.Fprintf(a.output, "✓ Template is valid\n")
	return nil
}

// handleEval processes the eval command.
func (a *App) handleEval(c *cli.Context) error {
	fmt.Fprintf(a.output, "Eval command not yet implemented\n")
	fmt.Fprintf(a.output, "This will evaluate dynamic namespaces in UP documents\n")
	return nil
}

// handleConvert processes the convert command.
func (a *App) handleConvert(c *cli.Context) error {
	fmt.Fprintf(a.output, "Convert command not yet implemented\n")
	fmt.Fprintf(a.output, "This will convert between UP and other formats (JSON, YAML, TOML)\n")
	return nil
}

// handleLSP starts the language server.
func (a *App) handleLSP(c *cli.Context) error {
	// Try to exec up-language-server
	args := []string{"up-language-server"}
	if c.Bool("debug") {
		args = append(args, "--debug")
	}
	if logFile := c.String("log"); logFile != "" {
		args = append(args, "--log", logFile)
	}

	return execTool("up-language-server", args[1:])
}

// handleREPL starts the interactive REPL.
func (a *App) handleREPL(c *cli.Context) error {
	// Try to exec up-repl
	args := []string{}
	if c.Bool("debug") {
		args = append(args, "--debug")
	}

	return execTool("up-repl", args)
}

// handleTool dispatches to external tools.
func (a *App) handleTool(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("tool name required")
	}

	toolName := c.Args().First()
	toolArgs := c.Args().Tail()

	return execTool("up-"+toolName, toolArgs)
}

// handleVersion prints version information.
func (a *App) handleVersion(c *cli.Context) error {
	if c.Bool("short") {
		fmt.Fprintf(a.output, "%s\n", version)
		return nil
	}

	if c.Bool("json") {
		versionInfo := map[string]string{
			"version": version,
			"commit":  commit,
			"date":    date,
			"builtBy": builtBy,
		}
		data, err := json.MarshalIndent(versionInfo, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintf(a.output, "%s\n", data)
		return nil
	}

	fmt.Fprintf(a.output, "up version %s\n", version)
	fmt.Fprintf(a.output, "  commit:  %s\n", commit)
	fmt.Fprintf(a.output, "  date:    %s\n", date)
	fmt.Fprintf(a.output, "  built by: %s\n", builtBy)
	return nil
}

// execTool executes an external tool binary.
func execTool(name string, args []string) error {
	// Try to find the tool in PATH
	path, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Try same directory as up binary first
	toolPath := strings.TrimSuffix(path, "up") + name
	if _, err := os.Stat(toolPath); os.IsNotExist(err) {
		// Try PATH
		var lookErr error
		toolPath, lookErr = exec.LookPath(name)
		if lookErr != nil {
			return fmt.Errorf("tool %q not found (install with: go install github.com/uplang/tools/%s@latest)", name, strings.TrimPrefix(name, "up-"))
		}
	}

	// Execute the tool
	cmd := exec.Command(toolPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("failed to run tool: %w", err)
	}

	return nil
}

// nopWriteCloser wraps an io.Writer to implement io.WriteCloser.
type nopWriteCloser struct {
	io.Writer
}

func (nwc *nopWriteCloser) Close() error {
	return nil
}

func main() {
	app := DefaultApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

