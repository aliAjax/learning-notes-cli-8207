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

func (a *App) runEdit(ctx context.Context, args []string, store storage.Store) int {
	fs := a.newFlagSet("edit")
	id := fs.String("id", "", "note id")
	title := fs.String("title", "", "new note title")
	tags := fs.String("tags", "", "new comma or semicolon separated tags")
	content := fs.String("content", "", "new Markdown content")
	stdin := fs.Bool("stdin", false, "read new Markdown content from stdin")
	showHelp := fs.Bool("h", false, "show help")
	fs.BoolVar(showHelp, "help", false, "show help")
	fs.Usage = a.printEditHelp

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showHelp {
		a.printEditHelp()
		return 0
	}

	noteID, err := resolveID(*id, fs.Args())
	if err != nil {
		a.printError(err)
		return 2
	}

	titleSet := containsFlag(fs, "title")
	tagsSet := containsFlag(fs, "tags")
	contentSet := containsFlag(fs, "content")
	if !titleSet && !tagsSet && !contentSet && !*stdin {
		a.printError(fmt.Errorf("provide at least one of --title, --tags, --content, or --stdin"))
		return 2
	}

	note, err := store.Get(ctx, noteID)
	if err != nil {
		if strings.Contains(err.Error(), storage.ErrNotFound.Error()) {
			a.printError(fmt.Errorf("note %q not found", noteID))
		} else {
			a.printError(err)
		}
		return 1
	}

	if titleSet {
		title := strings.TrimSpace(*title)
		if title == "" {
			a.printError(fmt.Errorf("--title cannot be empty"))
			return 2
		}
		note.Title = title
	}
	if tagsSet {
		note.Tags = model.ParseTags(*tags)
	}
	if *stdin {
		data, err := io.ReadAll(a.in)
		if err != nil {
			a.printError(fmt.Errorf("read content from stdin: %w", err))
			return 1
		}
		note.Content = string(data)
	} else if contentSet {
		note.Content = *content
	}
	note.UpdatedAt = time.Now().UTC()

	if err := store.Save(ctx, note); err != nil {
		a.printError(err)
		return 1
	}
	fmt.Fprintf(a.out, "Updated note %s %q\n", note.ID, note.Title)
	return 0
}

func (a *App) printEditHelp() {
	a.printCommandHelp(
		"edit",
		appName+" edit --id <id> [--title <title>] [--tags <tags>] [--content <text> | --stdin]",
		`  -h, --help       show help
      --id string     note id (required)
      --title string  new note title
      --tags string   new comma or semicolon separated tags
      --content string
                      new Markdown note content
      --stdin         read new Markdown content from stdin`,
		`  learning-notes edit --id 20260816-... --tags go,notes
  learning-notes edit --id 20260816-... --title "Updated title"`,
	)
}

func resolveID(flagValue string, positional []string) (string, error) {
	if flagValue != "" && len(positional) > 0 {
		return "", fmt.Errorf("provide the note id either as --id or positional, not both")
	}
	if flagValue == "" {
		if len(positional) != 1 {
			return "", fmt.Errorf("--id is required")
		}
		flagValue = positional[0]
	}
	flagValue = strings.TrimSpace(flagValue)
	if flagValue == "" {
		return "", fmt.Errorf("note id cannot be empty")
	}
	return flagValue, nil
}
