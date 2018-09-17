package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

const maintainerHelp = `List the maintainer for the Dockerfile(s).`

func (cmd *maintainerCommand) Name() string      { return "maintainer" }
func (cmd *maintainerCommand) Args() string      { return "[OPTIONS] DOCKERFILE [DOCKERFILE...]" }
func (cmd *maintainerCommand) ShortHelp() string { return maintainerHelp }
func (cmd *maintainerCommand) LongHelp() string  { return maintainerHelp }
func (cmd *maintainerCommand) Hidden() bool      { return false }

func (cmd *maintainerCommand) Register(fs *flag.FlagSet) {}

type maintainerCommand struct{}

func (cmd *maintainerCommand) Run(ctx context.Context, args []string) error {
	maintainers := map[string]int{}

	err := forFile(args, func(f string, nodes []*parser.Node) error {
		for _, n := range nodes {
			maintainers = nodeSearch("maintainer", n, maintainers)
			maintainers = labelSearch("maintainer", n, maintainers)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// setup the tab writer
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	// print header
	fmt.Fprintln(w, "MAINTAINER\tCOUNT")

	pl := rank(maintainers)
	for _, p := range pl {
		fmt.Fprintf(w, "%s\t%d\n", p.Key, p.Value)
	}

	w.Flush()
	return nil
}
