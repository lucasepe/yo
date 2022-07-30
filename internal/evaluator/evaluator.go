package evaluator

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/lucasepe/yo/internal/parser"
	"gopkg.in/yaml.v2"
)

type Evaluator struct {
	JSON bool
}

func (r *Evaluator) Eval(gens []parser.Generator) error {
	for _, g := range gens {
		if r.JSON {
			if err := toJSON(os.Stdout, g); err != nil {
				return err
			}
		} else {
			if err := toYAML(os.Stdout, g); err != nil {
				return err
			}
		}
	}

	return nil
}

func toYAML(w io.Writer, g parser.Generator) (err error) {
	dat, err := yaml.Marshal(g.Get())
	if err != nil {
		return err
	}

	_, err = w.Write(dat)
	return err
}

func toJSON(w io.Writer, g parser.Generator) (err error) {
	dat, err := json.Marshal(g.Get())
	if err != nil {
		return err
	}

	var out bytes.Buffer
	err = json.Indent(&out, dat, "", "   ")
	if err != nil {
		return err
	}
	_, err = w.Write(out.Bytes())
	return err
}
