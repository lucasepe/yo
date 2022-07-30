package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lucasepe/yo/internal/evaluator"
	"github.com/lucasepe/yo/internal/parser"
	"github.com/lucasepe/yo/internal/stdin"
	"github.com/lucasepe/yo/internal/strvals"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewCmdEval() *cobra.Command {
	opt := &evalCmd{
		optJSON: false,
	}

	cmd := &cobra.Command{
		Use: "eval [--json] <EXPRESSION SYNTAX>...",
		//DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Short:                 fmt.Sprintf("Evaluate a %s object notation syntax", strings.ToUpper(appName)),
		Example:               opt.examples(),
		RunE:                  opt.run,
	}

	cmd.Flags().BoolVarP(&opt.optJSON, "json", "j", opt.optJSON, "output format JSON (default: YAML)")
	cmd.Flags().StringSliceVar(&opt.setValues, "set", []string{}, "key=value pairs (take precedence over -values)")
	cmd.Flags().StringSliceVarP(&opt.values, "values", "f", []string{}, "specify values in a YAML or JSON files")

	return cmd
}

type evalCmd struct {
	optJSON   bool
	setValues []string
	values    []string
}

func (r *evalCmd) run(cmd *cobra.Command, args []string) error {
	// load and merge datasources
	ds, err := vals(r.values, r.setValues)
	if err != nil {
		return err
	}

	res, err := r.parseArgsOrStdIn(args, ds)
	if err != nil {
		return err
	}

	e := evaluator.Evaluator{JSON: r.optJSON}
	return e.Eval(res)
}

func (r *evalCmd) parseArgsOrStdIn(args []string, data map[string]interface{}) ([]parser.Generator, error) {
	if len(args) == 0 {
		return parser.ParseString(stdin.Input(), data)
	}
	return parser.ParseString(strings.Join(args, " "), data)
}

func (r *evalCmd) examples() string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	fmt.Fprintf(w, "  %s eval 'apiVersion=v1 kind=Namespace metadata={name=myxxx labels.name=myxxx}\n", appName)

	fmt.Fprintf(w, "  %s eval -j 'home = (env \"HOME\")'\n", appName)

	fmt.Fprintf(w, "  %s eval 'apiVersion=v1 kind=Secret metadata.name=mysecret type=Opaque ", appName)
	fmt.Fprintf(w, "data={ password=(b64enc \"PASS\") username=(b64enc \"USER\") }'")
	return buf.String()
}

// HELM CODE
// I really like how you can set values with helm... so using their code:
// https://github.com/kubernetes/helm/blob/master/cmd/helm/install.go

// vals merges values from files specified via -f/--values and
// directly via --set, marshaling them to YAML
func vals(valueFiles []string, values []string) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, filePath := range valueFiles {
		currentMap := map[string]interface{}{}

		var bytes []byte
		var err error
		if strings.TrimSpace(filePath) == "-" {
			bytes, err = ioutil.ReadAll(os.Stdin)
		} else {
			bytes, err = ioutil.ReadFile(filePath)
		}

		if err != nil {
			return map[string]interface{}{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return map[string]interface{}{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
		}
		// Merge with the previous map
		base = mergeValues(base, currentMap)
	}

	// User specified a value via --set
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return map[string]interface{}{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return base, nil
}

// Merges source and destination map, preferring values from the source map
func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
