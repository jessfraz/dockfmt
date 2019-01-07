package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

const workdirHelp = `List the workdirs for the Dockerfile(s).`

func (cmd *workdirCommand) Name() string      { return "workdir" }
func (cmd *workdirCommand) Args() string      { return "[OPTIONS] DOCKERFILE [DOCKERFILE...]" }
func (cmd *workdirCommand) ShortHelp() string { return workdirHelp }
func (cmd *workdirCommand) LongHelp() string  { return workdirHelp }
func (cmd *workdirCommand) Hidden() bool      { return false }

func (cmd *workdirCommand) Register(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.noRank, "no-rank", false, "turn off ranking of WORKDIRs")
	fs.BoolVar(&cmd.noRank, "N", false, "turn off ranking of WORKDIRs")

	fs.BoolVar(&cmd.skipHeader, "skip-header", false, "skip printing headers")
	fs.BoolVar(&cmd.skipHeader, "S", false, "skip printing headers")
}

type workdirCommand struct {
	noRank     bool
	skipHeader bool
}

func (cmd *workdirCommand) Run(ctx context.Context, args []string) error {
	workdirs := map[string]int{}

	err := forFile(args, func(f string, nodes []*parser.Node) error {
		for _, n := range nodes {
			workdirs = nodeSearch("workdir", n, workdirs)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// setup the tab writer
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	if cmd.noRank {
		if !cmd.skipHeader {
			// print header
			fmt.Fprintln(w, "WORKDIR")
		}

		for workdir := range workdirs {
			fmt.Fprintf(w, "%s\n", workdir)
		}
	} else {
		if !cmd.skipHeader {
			// print header
			fmt.Fprintln(w, "WORKDIR\tCOUNT")
		}

		pl := rank(workdirs)
		for _, p := range pl {
			fmt.Fprintf(w, "%s\t%d\n", p.Key, p.Value)
		}
	}

	w.Flush()
	return nil
}
