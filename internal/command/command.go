package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"learning-notes-cli/internal/config"
	"learning-notes-cli/internal/search"
	"learning-notes-cli/internal/storage"
)

const appName = "learning-notes"

// App owns the CLI's input and output streams.
type App struct {
	in     io.Reader
	out    io.Writer
	errOut io.Writer
}

// NewApp creates a command runner.
func NewApp(in io.Reader, out, errOut io.Writer) *App {
	return &App{in: in, out: out, errOut: errOut}
}

// Run dispatches the first command in args.
func (a *App) Run(ctx context.Context, args []string) int {
	dataDir, rest, err := parseGlobalFlags(args)
	if err != nil {
		a.printError(err)
		a.printRootHelp()
		return 2
	}

	if len(rest) == 0 {
		a.printRootHelp()
		return 0
	}
	if isRootHelp(rest[0]) {
		a.printRootHelp()
		return 0
	}

	cfg, err := a.resolveConfig(dataDir)
	if err != nil {
		a.printError(err)
		return 2
	}

	command := rest[0]
	commandArgs := rest[1:]
	store := storage.NewMarkdownStore(cfg.DataDir)
	searcher := search.Filter

	switch command {
	case "new":
		return a.runNew(ctx, commandArgs, store)
	case "list":
		return a.runList(ctx, commandArgs, store)
	case "search":
		return a.runSearch(ctx, commandArgs, store, searcher)
	case "edit":
		return a.runEdit(ctx, commandArgs, store)
	case "delete":
		return a.runDelete(ctx, commandArgs, store)
	case "help":
		a.runHelp(commandArgs)
		return 0
	default:
		a.printError(fmt.Errorf("unknown command %q", command))
		a.printRootHelp()
		return 2
	}
}

func (a *App) resolveConfig(flagValue string) (config.Config, error) {
	dataDir := strings.TrimSpace(flagValue)
	if dataDir == "" {
		dataDir = strings.TrimSpace(os.Getenv("LEARNING_NOTES_DATA_DIR"))
	}
	if dataDir == "" {
		dataDir = strings.TrimSpace(os.Getenv("NOTES_DATA_DIR"))
	}
	if dataDir == "" {
		defaultDir, err := config.DefaultDataDir()
		if err != nil {
			return config.Config{}, err
		}
		dataDir = defaultDir
	}

	cfg, err := config.New(dataDir)
	if err != nil {
		return config.Config{}, err
	}
	if err := cfg.EnsureDataDir(); err != nil {
		return config.Config{}, err
	}
	return cfg, nil
}

func parseGlobalFlags(args []string) (string, []string, error) {
	dataDir := ""
	position := 0

	for position < len(args) {
		arg := args[position]
		switch {
		case arg == "--data-dir":
			if position+1 >= len(args) {
				return "", nil, fmt.Errorf("--data-dir requires a value")
			}
			dataDir = args[position+1]
			position += 2
		case strings.HasPrefix(arg, "--data-dir="):
			dataDir = strings.TrimPrefix(arg, "--data-dir=")
			position++
		case arg == "-d":
			if position+1 >= len(args) {
				return "", nil, fmt.Errorf("-d requires a value")
			}
			dataDir = args[position+1]
			position += 2
		case strings.HasPrefix(arg, "-d="):
			dataDir = strings.TrimPrefix(arg, "-d=")
			position++
		default:
			return dataDir, args[position:], nil
		}
	}

	return dataDir, []string{}, nil
}

func isRootHelp(arg string) bool {
	return arg == "-h" || arg == "--help"
}

func (a *App) newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(a.errOut)
	return fs
}

func (a *App) printError(err error) {
	fmt.Fprintf(a.errOut, "error: %s\n", err)
}

func (a *App) printRootHelp() {
	fmt.Fprintf(a.out, `%s is a local Markdown learning notes tool.

Usage:
  %s [--data-dir <dir>] <command> [arguments]

Commands:
  new       create a Markdown note
  list      list all notes
  search    search notes by title or tag
  edit      edit an existing note
  delete    delete a note
  help      show help for the CLI or a command

Global options:
  -d, --data-dir <dir>   notes storage directory
                         env: LEARNING_NOTES_DATA_DIR or NOTES_DATA_DIR
  -h, --help             show this help

Run "%s help <command>" for command-specific help.
`, appName, appName, appName)
}

func (a *App) printCommandHelp(name string, usage string, options string, examples string) {
	fmt.Fprintf(a.out, "Usage:\n  %s\n\nOptions:\n%s\n", usage, options)
	if examples != "" {
		fmt.Fprintf(a.out, "\nExamples:\n%s\n", examples)
	}
}

func containsFlag(flagSet *flag.FlagSet, names ...string) bool {
	found := false
	flagSet.Visit(func(f *flag.Flag) {
		for _, name := range names {
			if f.Name == name {
				found = true
				return
			}
		}
	})
	return found
}

func isUsageError(err error) bool {
	return errors.Is(err, flag.ErrHelp)
}
