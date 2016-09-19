package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/urfave/cli"
)

func getMaintainer(c *cli.Context) error {
	maintainers := map[string]int{}

	err := forFile(c, func(f string, nodes []*parser.Node) error {
		for _, n := range nodes {
			maintainers = nodeSearch("maintainer", n, maintainers)
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
