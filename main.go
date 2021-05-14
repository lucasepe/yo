package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/ko/internal/parser"
	"github.com/lucasepe/ko/internal/stdin"
	"gopkg.in/yaml.v2"
)

const (
	banner = `┬ ┬┌─┐
└┬┘│ │
 ┴ └─┘ YAML object generator`
)

var (
	optJSON    bool
	optVersion bool

	commit string
)

func main() {
	configureFlags()

	if optVersion {
		fmt.Printf("%s version: %s\n", appName(), commit)
		os.Exit(0)
	}

	res, err := parseArgsOrStdIn()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	for _, g := range res {
		if optJSON {
			if err := toJSON(os.Stdout, g); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := toYAML(os.Stdout, g); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		fmt.Printf("%s\n\n", banner)
		//fmt.Print("Commandline Object Notation.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s <EXPRESSION SYNTAX>...\n\n", name)

		fmt.Print("EXAMPLES:\n\n")
		fmt.Printf("  %s -json user = { name=foo age=30 type=C }\n", name)
		fmt.Printf("  %s spec = { credentials.source=Secret credentials.secretRef = {name=aws-creds key=creds} }\n\n", name)

		fmt.Print("FLAGS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Print("  -help\n\tprints this message\n")
		fmt.Println()

		fmt.Println("crafted with passion @ 2021 by Luca Sepe <luca.sepe@gmail.com>")
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.CommandLine.BoolVar(&optJSON, "json", false, "output format JSON (default: YAML)")
	flag.CommandLine.BoolVar(&optVersion, "v", false, "print current version and exit")

	flag.CommandLine.Parse(os.Args[1:])
}

func appName() string {
	return filepath.Base(os.Args[0])
}

func parseArgsOrStdIn() ([]parser.Generator, error) {
	if len(flag.Args()) == 0 {
		return parser.ParseString(stdin.Input())
	}
	return parser.ParseString(strings.Join(flag.Args(), " "))
}

func toYAML(w io.Writer, g parser.Generator) error {
	dat, err := yaml.Marshal(g.Get())
	if err != nil {
		return err
	}

	w.Write(dat)

	return nil
}

func toJSON(w io.Writer, g parser.Generator) error {
	dat, err := json.Marshal(g.Get())
	if err != nil {
		return err
	}

	var out bytes.Buffer
	if err := json.Indent(&out, dat, "", "   "); err != nil {
		return err
	}
	w.Write(out.Bytes())

	return nil
}
