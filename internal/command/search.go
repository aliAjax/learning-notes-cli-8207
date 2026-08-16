package command

import (
	"context"
	"fmt"
	"strings"

	"learning-notes-cli/internal/model"
	"learning-notes-cli/internal/search"
	"learning-notes-cli/internal/storage"
)

type searchFunc func([]model.Note, search.Query) []model.Note

func (a *App) runSearch(ctx context.Context, args []string, store storage.Store, searcher searchFunc) int {
	fs := a.newFlagSet("search")
	title := fs.String("title", "", "match note titles containing this text")
	tag := fs.String("tag", "", "match notes with this tag")
	jsonOut := fs.Bool("json", false, "output notes as JSON")
	showHelp := fs.Bool("h", false, "show help")
	fs.BoolVar(showHelp, "help", false, "show help")
	fs.Usage = a.printSearchHelp

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showHelp {
		a.printSearchHelp()
		return 0
	}

	query := search.Query{
		Text:  strings.Join(fs.Args(), " "),
		Title: *title,
		Tag:   *tag,
	}
	if query.IsEmpty() {
		a.printError(fmt.Errorf("provide a query, --title, or --tag"))
		return 2
	}

	notes, err := store.List(ctx)
	if err != nil {
		a.printError(err)
		return 1
	}
	a.printNotes(notes, *jsonOut)
	return 0
}

func (a *App) printSearchHelp() {
	a.printCommandHelp(
		"search",
		appName+" search [--title <text>] [--tag <tag>] [--json] [query...]",
		`  -h, --help       show help
      --title string  match note titles
      --tag string    match note tags
      --json          output notes as JSON`,
		`  learning-notes search goroutines
  learning-notes search --tag go
  learning-notes search --title "learning log"`,
	)
}
