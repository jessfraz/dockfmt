package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

const stageHelp = `List the stages in the Dockerfile.`

func (cmd *stagesCommand) Name() string      { return "stages" }
func (cmd *stagesCommand) Args() string      { return "[OPTIONS] DOCKERFILE" }
func (cmd *stagesCommand) ShortHelp() string { return stageHelp }
func (cmd *stagesCommand) LongHelp() string  { return stageHelp }
func (cmd *stagesCommand) Hidden() bool      { return false }

func (cmd *stagesCommand) Register(fs *flag.FlagSet) {}

type stagesCommand struct{}

func (cmd *stagesCommand) Run(ctx context.Context, args []string) error {
	images := []*parser.Node{}

	err := forFile(args, func(f *os.File, nodes []*parser.Node) error {
		for _, n := range nodes {
			if n.Value == "from" {
				images = append(images, n)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// setup the tab writer
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	// print header
	fmt.Fprintln(w, "STAGE\tINTERPOLATED")
	for i, n := range images {
		cmd, err := instructions.ParseInstruction(n)
		if err != nil {
			w.Flush()
			return err
		}
		switch stage := cmd.(type) {
		case *instructions.Stage:
			stageName := stage.Name
			interpolated := false

			if stageName == "" {
				stageName = fmt.Sprintf("stage-%d", i)
				interpolated = true
			}

			fmt.Fprintf(w, "%s\t%v\n", stageName, interpolated)
		}
	}

	w.Flush()

	return nil
}
