package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

const dumpHelp = `Dump parsed Dockerfile(s).`

func (cmd *dumpCommand) Name() string      { return "dump" }
func (cmd *dumpCommand) Args() string      { return "[OPTIONS] DOCKERFILE [DOCKERFILE...]" }
func (cmd *dumpCommand) ShortHelp() string { return dumpHelp }
func (cmd *dumpCommand) LongHelp() string  { return dumpHelp }
func (cmd *dumpCommand) Hidden() bool      { return false }

func (cmd *dumpCommand) Register(fs *flag.FlagSet) {}

type dumpCommand struct{}

func (cmd *dumpCommand) Run(ctx context.Context, args []string) error {
	return forFile(args, func(f *os.File, nodes []*parser.Node) error {
		fmt.Println(f.Name())
		if len(nodes) > 0 {
			fmt.Println(nodes[0].Dump())
		}
		return nil
	})
}
