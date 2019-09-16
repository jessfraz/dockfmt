package main

import (
	"context"
	"flag"
	"os"
	"sort"
	"strings"

	"github.com/genuinetools/pkg/cli"
	"github.com/jessfraz/dockfmt/version"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/sirupsen/logrus"
)

var (
	debug bool
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "dockfmt"
	p.Description = "Dockerfile format."

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	p.Commands = []cli.Command{
		&baseCommand{},
		&dumpCommand{},
		&formatCommand{},
		&maintainerCommand{},
		&stagesCommand{},
		&workdirCommand{},
	}

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}


		if p.FlagSet.NArg() < 1 {
			return errors.New("please pass in Dockerfile(s)")
		}
    
		return nil
	}

	// Run our program.
	p.Run()
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

func labelSearch(search string, n *parser.Node, a map[string]int) map[string]int {
	if n.Value == "label" {
		if n.Next != nil && strings.EqualFold(strings.ToLower(n.Next.Value), strings.ToLower(search)) {
			i := strings.Trim(n.Next.Next.Value, "\"")
			if v, ok := a[i]; ok {
				a[i] = v + 1
			} else {
				a[i] = 1

			}
		}
	}
	return a
}

func nodeSearch(search string, n *parser.Node, a map[string]int) map[string]int {
	if n.Value == search {
		i := strings.Trim(n.Next.Value, "\"")
		if v, ok := a[i]; ok {
			a[i] = v + 1
		} else {
			a[i] = 1

		}
	}
	return a
}

type forFileFunc func(*os.File, []*parser.Node) error

func forFile1(f *os.File, fnc forFileFunc) error {
	result, err := parser.Parse(f)
	if err != nil {
		return err
	}
	ast := result.AST
	nodes := []*parser.Node{ast}
	if ast.Children != nil {
		nodes = append(nodes, ast.Children...)
	}
	return fnc(f, nodes)
}

func forFile(args []string, fnc forFileFunc) error {
	if len(args) == 0 {
		return forFile1(os.Stdin, fnc)
	}
	for _, fn := range args {
		logrus.Debugf("parsing file: %s", fn)

		f, err := os.Open(fn)
		if err != nil {
			return err
		}
		defer f.Close()

		if err = forFile1(f, fnc); err != nil {
			return err
		}
	}
	return nil
}
