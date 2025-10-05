package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/uplang/tools/language-server/server"
	"github.com/urfave/cli/v2"
)

var version = "1.0.0"

func main() {
	app := &cli.App{
		Name:    "up-language-server",
		Usage:   "Language Server Protocol implementation for UP",
		Version: version,
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
		Action: func(c *cli.Context) error {
			// Setup logger
			logger := setupLogger(c.Bool("debug"), c.String("log"))

			logger.Info("UP Language Server starting",
				slog.String("version", version),
				slog.Bool("debug", c.Bool("debug")),
			)

			// Create and start server
			srv := server.NewServer(logger)

			ctx := context.Background()
			if err := srv.Run(ctx); err != nil {
				logger.Error("Server error", slog.Any("error", err))
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setupLogger(debug bool, logFile string) *slog.Logger {
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
			os.Exit(1)
		}
		handler = slog.NewJSONHandler(f, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

