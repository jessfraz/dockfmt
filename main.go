package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/jessfraz/dockfmt/version"
	"github.com/urfave/cli"
)

// preload initializes any global options and configuration
// before the main or sub commands are run.
func preload(c *cli.Context) (err error) {
	if c.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if len(c.Args()) < 1 {
		return errors.New("please supply filename(s)")
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "dockfmt"
	app.Version = version.VERSION
	app.Author = "@jessfraz"
	app.Email = "no-reply@butts.com"
	app.Usage = "Dockerfile format."
	app.Before = preload
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "run in debug mode",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "base",
			Usage:  "list the base image used in Dockerfile(s)",
			Action: getBase,
		},
		{
			Name:  "dump",
			Usage: "dump parsed Dockerfile(s)",
			Action: func(c *cli.Context) error {
				err := forFile(c, func(f string, nodes []*parser.Node) error {
					fmt.Println(f)
					if len(nodes) > 0 {
						fmt.Println(nodes[0].Dump())
					}
					return nil
				})
				return err
			},
		},
		{
			Name:    "format",
			Aliases: []string{"fmt"},
			Usage:   "format the Dockerfile(s)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "diff, d",
					Usage: "display diffs instead of rewriting files",
				},
				cli.BoolFlag{
					Name:  "list, l",
					Usage: "list files whose formatting differs from dockfmt's",
				},
				cli.BoolFlag{
					Name:  "write, w",
					Usage: "write result to (source) file instead of stdout",
				},
			},
			Action: format,
		},
		{
			Name:   "maintainer",
			Usage:  "list the maintainer for Dockerfile(s)",
			Action: getMaintainer,
		},
	}

	app.Run(os.Args)
}

type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func rank(images map[string]int) pairList {
	pl := make(pairList, len(images))
	i := 0
	for k, v := range images {
		pl[i] = pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func nodeSearch(search string, n *parser.Node, a map[string]int) map[string]int {
	if n.Value == search {
		if v, ok := a[n.Next.Value]; ok {
			a[n.Next.Value] = v + 1
		} else {
			a[n.Next.Value] = 1

		}
	}
	return a
}

func forFile(c *cli.Context, fnc func(string, []*parser.Node) error) error {
	for _, fn := range c.Args() {
		logrus.Debugf("File: %s", fn)
		f, err := os.Open(fn)
		if err != nil {
			return err
		}
		defer f.Close()

		d := parser.Directive{LookingForDirectives: true}
		parser.SetEscapeToken(parser.DefaultEscapeToken, &d)

		ast, err := parser.Parse(f, &d)
		if err != nil {
			return err
		}

		nodes := []*parser.Node{ast}
		if ast.Children != nil {
			nodes = append(nodes, ast.Children...)
		}
		if err := fnc(fn, nodes); err != nil {
			return err
		}
	}
	return nil
}
