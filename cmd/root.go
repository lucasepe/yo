package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	appName = "yo"
	banner  = `┏┓ ┏┓ ┏━━┓
┃┗━┛┃ ┃┏┓┃
┗━┓┏┛ ┃┗┛┃
┗━━┛  ┗━━┛`
	appSummary = "The ultimate commanline YAML (or JSON) generator!"
)

func Run() *cobra.Command {
	cmd := &cobra.Command{
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Use:   fmt.Sprintf("%s <COMMAND>", appName),
		Short: appSummary,
		Long:  fmt.Sprintf("%s\n%s\n", banner, appSummary),
	}

	cmd.AddCommand(NewCmdVersion())
	cmd.AddCommand(NewCmdEval())
	cmd.AddCommand(NewCmdFunctions())

	return cmd
}
