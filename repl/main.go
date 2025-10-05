package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/uplang/go"
	"github.com/urfave/cli/v2"
)

var version = "1.0.0"

func main() {
	app := &cli.App{
		Name:    "up-repl",
		Usage:   "Interactive REPL for UP (Unified Properties)",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug output",
			},
		},
		Action: runREPL,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runREPL(c *cli.Context) error {
	logger := setupLogger(c.Bool("debug"))

	rl, err := readline.New("up> ")
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	parser := up.NewParser()

	fmt.Println("UP REPL v" + version)
	fmt.Println("Type UP expressions and press Enter to evaluate.")
	fmt.Println("Commands: .help, .exit, .clear")
	fmt.Println()

	var buffer strings.Builder
	multiline := false

	for {
		prompt := "up> "
		if multiline {
			prompt = "... "
		}
		rl.SetPrompt(prompt)

		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		// Handle commands
		if strings.HasPrefix(line, ".") {
			if handleCommand(line, &buffer, &multiline) {
				break
			}
			continue
		}

		// Build multiline input
		if multiline {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)

		// Check for block/list delimiters
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "[") {
			multiline = true
			continue
		}

		if line == "}" || line == "]" {
			multiline = false
		}

		// If not in multiline mode, try to parse
		if !multiline {
			input := buffer.String()
			if input != "" {
				evaluateInput(parser, input, logger)
			}
			buffer.Reset()
		}
	}

	return nil
}

func handleCommand(cmd string, buffer *strings.Builder, multiline *bool) bool {
	switch cmd {
	case ".help", ".h":
		fmt.Println("Commands:")
		fmt.Println("  .help, .h     - Show this help")
		fmt.Println("  .exit, .quit  - Exit REPL")
		fmt.Println("  .clear, .c    - Clear input buffer")
		fmt.Println("\nExamples:")
		fmt.Println("  name John Doe")
		fmt.Println("  age!int 30")
		fmt.Println("  config {")
		fmt.Println("    debug!bool true")
		fmt.Println("  }")

	case ".exit", ".quit", ".q":
		fmt.Println("Goodbye!")
		return true

	case ".clear", ".c":
		buffer.Reset()
		*multiline = false
		fmt.Println("Buffer cleared.")

	default:
		fmt.Printf("Unknown command: %s (type .help for commands)\n", cmd)
	}

	return false
}

func evaluateInput(parser *up.Parser, input string, logger *slog.Logger) {
	doc, err := parser.ParseDocument(strings.NewReader(input))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Pretty print the parsed document
	fmt.Println("Parsed:")
	for _, node := range doc.Nodes {
		fmt.Printf("  %s", node.Key)
		if node.Type != "" {
			fmt.Printf("!%s", node.Type)
		}
		fmt.Printf(": %v\n", formatValue(node.Value))
	}
}

func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case map[string]any:
		return fmt.Sprintf("{ %d entries }", len(val))
	case []any:
		return fmt.Sprintf("[ %d items ]", len(val))
	default:
		return fmt.Sprintf("%v", val)
	}
}

func setupLogger(debug bool) *slog.Logger {
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	return slog.New(handler)
}

