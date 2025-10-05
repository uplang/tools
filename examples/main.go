package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorYellow = "\033[33m"
)

func main() {
	app := &cli.App{
		Name:  "examples",
		Usage: "Test runner for UP namespace example files",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose output",
			},
		},
		Action: runExamples,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runExamples(c *cli.Context) error {
	fmt.Println("=========================================")
	fmt.Println("UP Namespace Examples Runner")
	fmt.Println("=========================================")
	fmt.Println()

	// Get the directory to search (default: current directory)
	searchDir := "."
	if c.NArg() > 0 {
		searchDir = c.Args().Get(0)
	}

	// Find all .up files in examples directories
	var upFiles []string
	err := filepath.WalkDir(searchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".up") && strings.Contains(path, "/examples/") {
			upFiles = append(upFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	if len(upFiles) == 0 {
		fmt.Println("No .up example files found in examples/ directories")
		os.Exit(1)
	}

	fmt.Printf("Found %d example files\n\n", len(upFiles))

	// Check if UP CLI is available
	upCmd, err := exec.LookPath("up")
	hasUpCLI := err == nil

	// Process each file
	for i, file := range upFiles {
		// Extract namespace from path
		parts := strings.Split(file, string(filepath.Separator))
		namespace := ""
		for j, part := range parts {
			if j > 0 && parts[j-1] == "." || (j > 1 && parts[j-2] == "..") {
				namespace = part
				break
			}
			if strings.Contains(file, "/"+part+"/examples/") {
				namespace = part
				break
			}
		}
		if namespace == "" {
			namespace = "unknown"
		}

		basename := filepath.Base(file)

		fmt.Printf("%s[%d/%d]%s Testing: %s%s%s - %s\n",
			colorBlue, i+1, len(upFiles), colorReset,
			colorGreen, namespace, colorReset,
			basename)
		fmt.Printf("File: %s\n\n", file)

		if hasUpCLI {
			// Try to parse the file with UP CLI
			cmd := exec.Command(upCmd, "parse", "-i", file, "--pretty")
			output, err := cmd.CombinedOutput()

			fmt.Println("Output:")
			if err != nil {
				fmt.Printf("  %s(Parse error - may require namespace evaluation)%s\n", colorYellow, colorReset)
				if len(output) > 0 {
					fmt.Printf("  %s\n", string(output))
				}
			} else {
				fmt.Println(string(output))
			}
		} else {
			fmt.Printf("%sNote: UP CLI not installed%s\n", colorYellow, colorReset)
			fmt.Println("  Install with: go install github.com/uplang/go/cmd/up@latest")
			fmt.Printf("  To test manually: up parse -i %s --pretty\n", file)
		}

		fmt.Println()
		fmt.Println("---")
		fmt.Println()
	}

	fmt.Println("=========================================")
	fmt.Println("All examples processed!")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("Note: These examples show UP syntax parsing.")
	fmt.Println("To evaluate with actual namespace functions:")
	fmt.Println("1. Build namespaces: cd <namespace> && go build")
	fmt.Println("2. Use UP CLI with full namespace support")
	fmt.Println()

	return nil
}

