package command

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"learning-notes-cli/internal/model"
	"learning-notes-cli/internal/storage"
)

func (a *App) runList(ctx context.Context, args []string, store storage.Store) int {
	fs := a.newFlagSet("list")
	jsonOut := fs.Bool("json", false, "output notes as JSON")
	showHelp := fs.Bool("h", false, "show help")
	fs.BoolVar(showHelp, "help", false, "show help")
	fs.Usage = a.printListHelp

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showHelp {
		a.printListHelp()
		return 0
	}
	if fs.NArg() != 0 {
		a.printError(fmt.Errorf("list does not accept positional arguments"))
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

func (a *App) printNotes(notes []model.Note, jsonOut bool) {
	if jsonOut {
		encoder := json.NewEncoder(a.out)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(notes); err != nil {
			a.printError(fmt.Errorf("encode notes JSON: %w", err))
		}
		return
	}

	if len(notes) == 0 {
		fmt.Fprintln(a.out, "No notes found.")
		return
	}

	writer := tabwriter.NewWriter(a.out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tUPDATED\tTITLE\tTAGS")
	for _, note := range notes {
		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\n",
			note.ID,
			note.UpdatedAt.Local().Format("2006-01-02 15:04"),
			singleLine(note.Title),
			singleLine(note.TagString()),
		)
	}
	_ = writer.Flush()
}

func (a *App) printListHelp() {
	a.printCommandHelp(
		"list",
		appName+" list [--json]",
		`  -h, --help   show help
      --json     output notes as JSON`,
		`  learning-notes list
  learning-notes list --json`,
	)
}

func singleLine(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
