package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/urfave/cli"
)

func getBase(c *cli.Context) error {
	images := map[string]int{}

	err := forFile(c, func(f string, nodes []*parser.Node) error {
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
