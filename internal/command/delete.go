package command

import (
	"context"
	"fmt"
	"strings"

	"learning-notes-cli/internal/storage"
)

func (a *App) runDelete(ctx context.Context, args []string, store storage.Store) int {
	fs := a.newFlagSet("delete")
	id := fs.String("id", "", "note id")
	showHelp := fs.Bool("h", false, "show help")
	fs.BoolVar(showHelp, "help", false, "show help")
	fs.Usage = a.printDeleteHelp

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showHelp {
		a.printDeleteHelp()
		return 0
	}

	noteID, err := resolveID(*id, fs.Args())
	if err != nil {
		a.printError(err)
		return 2
	}

	if err := store.Delete(ctx, noteID); err != nil {
		if strings.Contains(err.Error(), storage.ErrNotFound.Error()) {
			a.printError(fmt.Errorf("note %q not found", noteID))
		} else {
			a.printError(err)
		}
		return 1
	}

	fmt.Fprintf(a.out, "Deleted note %s\n", noteID)
	return 0
}

func (a *App) printDeleteHelp() {
	a.printCommandHelp(
		"delete",
		appName+" delete --id <id>",
		`  -h, --help     show help
      --id string   note id (required)`,
		`  learning-notes delete --id 20260816-...`,
	)
}
