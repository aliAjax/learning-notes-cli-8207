package command

func (a *App) runHelp(args []string) {
	if len(args) == 0 {
		a.printRootHelp()
		return
	}

	switch args[0] {
	case "new":
		a.printNewHelp()
	case "list":
		a.printListHelp()
	case "search":
		a.printSearchHelp()
	case "edit":
		a.printEditHelp()
	case "delete":
		a.printDeleteHelp()
	case "help":
		a.printRootHelp()
	default:
		a.printError(errUnknownCommand(args[0]))
		a.printRootHelp()
	}
}

func errUnknownCommand(name string) error {
	return &commandError{message: "unknown command " + name}
}

type commandError struct {
	message string
}

func (e *commandError) Error() string {
	return e.message
}
