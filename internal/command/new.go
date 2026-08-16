package command

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"learning-notes-cli/internal/model"
	"learning-notes-cli/internal/storage"
)

func (a *App) runNew(ctx context.Context, args []string, store storage.Store) int {
	fs := a.newFlagSet("new")
	title := fs.String("title", "", "note title")
	tags := fs.String("tags", "", "comma or semicolon separated tags")
	content := fs.String("content", "", "Markdown content")
	stdin := fs.Bool("stdin", false, "read Markdown content from stdin")
	showHelp := fs.Bool("h", false, "show help")
	fs.BoolVar(showHelp, "help", false, "show help")
	fs.Usage = a.printNewHelp

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showHelp {
		a.printNewHelp()
		return 0
	}

	if strings.TrimSpace(*title) == "" {
		a.printError(fmt.Errorf("--title is required"))
		return 2
	}

	noteContent := *content
	if *stdin {
		data, err := io.ReadAll(a.in)
		if err != nil {
			a.printError(fmt.Errorf("read content from stdin: %w", err))
			return 1
		}
		noteContent = string(data)
	}
	if strings.TrimSpace(noteContent) == "" {
		a.printError(fmt.Errorf("note content is required; use --content or --stdin"))
		return 2
	}

	note, err := model.NewNote(*title, model.ParseTags(*tags), noteContent, time.Now().UTC())
	if err != nil {
		a.printError(err)
		return 1
	}
	if err := store.Save(ctx, note); err != nil {
		a.printError(err)
		return 1
	}

	fmt.Fprintf(a.out, "Created note %s %q\n", note.ID, note.Title)
	return 0
}

func (a *App) printNewHelp() {
	a.printCommandHelp(
		"new",
		appName+" new --title <title> [--tags <tags>] (--content <text> | --stdin)",
		`  -h, --help       show help
      --title string  note title (required)
      --tags string   comma or semicolon separated tags
      --content string
                      Markdown note content
      --stdin         read Markdown content from stdin`,
		`  learning-notes new --title "Go goroutines" --tags go,concurrency --content "Channels are safe."
  echo "# Today" | learning-notes new --title "Learning log" --stdin`,
	)
}
