package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

const baseHelp = `List the base image used in the Dockerfile(s).`

func (cmd *baseCommand) Name() string      { return "base" }
func (cmd *baseCommand) Args() string      { return "[OPTIONS] [DOCKERFILE...]" }
func (cmd *baseCommand) ShortHelp() string { return baseHelp }
func (cmd *baseCommand) LongHelp() string  { return baseHelp }
func (cmd *baseCommand) Hidden() bool      { return false }

func (cmd *baseCommand) Register(fs *flag.FlagSet) {}

type baseCommand struct{}

func (cmd *baseCommand) Run(ctx context.Context, args []string) error {
	images := map[string]int{}

	err := forFile(args, func(f *os.File, nodes []*parser.Node) error {
		for _, n := range nodes {
			images = nodeSearch("from", n, images)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// setup the tab writer
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	// print header
	fmt.Fprintln(w, "BASE\tCOUNT")

	pl := rank(images)
	for _, p := range pl {
		fmt.Fprintf(w, "%s\t%d\n", p.Key, p.Value)
	}

	w.Flush()

	return nil
}
