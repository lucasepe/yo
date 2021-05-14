package stdin

import (
	"io"
	"strings"

	"github.com/chzyer/readline"
)

func Input() string {
	lines := []string{}

	l, err := readline.NewEx(&readline.Config{
		Prompt:              ">> ",
		HistoryFile:         "/tmp/readline.tmp",
		AutoComplete:        autoCompleter(),
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		ln, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(ln) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		lines = append(lines, strings.TrimSpace(ln))
	}

	return strings.Join(lines, " ")
}

func autoCompleter() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("apiVersion"),
		readline.PcItem("kind"),
		readline.PcItem("metadata"),
		readline.PcItem("spec"),
		readline.PcItem("labels"),
	)
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
