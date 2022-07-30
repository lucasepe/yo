package cmd

import (
	"fmt"
	"strings"

	"github.com/lucasepe/yo/internal/table"
	"github.com/lucasepe/yo/internal/template"
	"github.com/spf13/cobra"
)

// NewCmdFunctions creates a command object for the "funcs" command
func NewCmdFunctions() *cobra.Command {
	return &cobra.Command{
		Use:   "funcs",
		Short: "Print the builtin functions documentation",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			tbl := &table.TextTable{}
			tbl.SetHeader(strings.ToUpper("Function"), strings.ToUpper("Summary"))

			names := template.Names()
			for i, n := range names {
				summary := template.Summary(n)
				sample := template.Usage(n)

				tbl.AddRow(n, summary)
				tbl.AddRow("", "")
				tbl.AddRow("", fmt.Sprintf(">> %s eval 'key = ( %s )'", appName, sample))

				if i < len(names)-1 {
					tbl.AddRowLine()
				}
			}

			fmt.Println(tbl.Draw())

		},
	}
}
