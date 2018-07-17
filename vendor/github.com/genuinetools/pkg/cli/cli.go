package cli

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
)

const (
	// GitCommitKey is the key for the program's GitCommit data.
	GitCommitKey ContextKey = "program.GitCommit"
	// VersionKey is the key for the program's Version data.
	VersionKey ContextKey = "program.Version"
)

// ContextKey defines the type for holding keys in the context.
type ContextKey string

// Program defines the struct for holding information about the program.
type Program struct {
	// Name of the program. Defaults to path.Base(os.Args[0]).
	Name string
	// Description of the program.
	Description string
	// Version of the program.
	Version string
	// GitCommit information for the program.
	GitCommit string

	// Commands in the program.
	Commands []Command
	// FlagSet holds the common/global flags for the program.
	FlagSet *flag.FlagSet

	// Before defines a function to execute before any subcommands are run,
	// but after the context is ready.
	// If a non-nil error is returned, no subcommands are run.
	Before func(context.Context) error
	// After defines a function to execute after any subcommands are run,
	// but after the subcommand has finished.
	// It is run even if the subcommand returns an error.
	After func(context.Context) error

	// Action is the function to execute when no subcommands are specified.
	// It gives the user back the arguments after the flags have been parsed.
	Action func(context.Context, []string) error
}

// Command defines the interface for each command in a program.
type Command interface {
	Name() string      // "foobar"
	Args() string      // "<baz> [quux...]"
	ShortHelp() string // "Foo the first bar"
	LongHelp() string  // "Foo the first bar meeting the following conditions..."

	// Hidden indicates whether the command should be hidden from the help output.
	Hidden() bool

	// Register command specific flags.
	Register(*flag.FlagSet)
	// Run executes the function for the command with a context and the command arguments.
	Run(context.Context, []string) error
}

// NewProgram creates a new Program with some reasonable defaults for Name,
// Description, and Version.
func NewProgram() *Program {
	return &Program{
		Name:        filepath.Base(os.Args[0]),
		Description: "A new command line program.",
		Version:     "0.0.0",
	}
}

// Run is the entry point for the program. It parses the arguments and executes
// the commands.
func (p *Program) Run() {
	ctx := p.defaultContext()

	// Pass the os.Args through so we can more easily unit test.
	printUsage, err := p.run(ctx, os.Args)
	if err == nil && !printUsage {
		return
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if printUsage {
		if err != nil {
			// Print an extra new line to seperate from the usage output.
			fmt.Fprintln(os.Stderr)
		}
		p.usage(ctx)
	}
	os.Exit(1)
}

func (p *Program) run(ctx context.Context, args []string) (bool, error) {
	// Append the version command to the list of commands by default.
	p.Commands = append(p.Commands, &versionCommand{})

	// IF
	// args is <nil>
	// OR
	// args is less than zero
	// OR
	// we have more than one arg and it equals help OR is a help flag
	// THEN
	// printUsage
	if args == nil ||
		len(args) < 1 ||
		(len(args) > 1 && contains([]string{"-h", "--help", "help"}, args[1])) {
		return true, nil
	}

	// If we do not have an action set and we have no commands, print the usage
	// and exit.
	if p.Action == nil && len(p.Commands) < 2 {
		return true, nil
	}

	// Check if the command exists.
	var commandExists bool
	if len(args) > 1 && in(args[1], p.Commands) {
		commandExists = true
	}

	// If we are not running a commands we know, then automatically
	// run the main action of the program instead.
	// Also enter this loop if we weren't passed any arguments.
	if p.Action != nil &&
		(len(args) < 2 || !commandExists) {
		return p.runAction(ctx, args)
	}

	// Return early if we didn't enter the single action logic and
	// the command does not exist or we were passed no commands.
	if len(args) < 2 {
		return true, nil
	}
	if !commandExists {
		return true, fmt.Errorf("%s: no such command", args[1])
	}

	// Iterate over the commands in the program.
	for _, command := range p.Commands {
		if args[1] == command.Name() {
			// Set the default flagset if our flagset is undefined.
			if p.FlagSet == nil {
				p.FlagSet = defaultFlagSet(p.Name)
			}

			// Register the subcommand flags in with the common/global flags.
			command.Register(p.FlagSet)

			// Override the usage text to something nicer.
			p.resetCommandUsage(command)

			// Parse the flags the user gave us.
			if err := p.FlagSet.Parse(args[2:]); err != nil {
				return false, err
			}

			if p.Before != nil {
				if err := p.Before(ctx); err != nil {
					return false, err
				}
			}

			// Run the command with the context and post-flag-processing args.
			if err := command.Run(ctx, p.FlagSet.Args()); err != nil {
				if p.After != nil {
					p.After(ctx)
				}

				return false, err
			}

			// Run the after function.
			if p.After != nil {
				if err := p.After(ctx); err != nil {
					return false, err
				}
			}
		}
	}

	// Done.
	return false, nil
}

func (p *Program) runAction(ctx context.Context, args []string) (bool, error) {
	// Set the default flagset if our flagset is undefined.
	if p.FlagSet == nil {
		p.FlagSet = defaultFlagSet(p.Name)
	}

	// Override the usage text to something nicer.
	p.FlagSet.Usage = func() {
		p.usage(ctx)
	}

	// Parse the flags the user gave us.
	if err := p.FlagSet.Parse(args[1:]); err != nil {
		return true, nil
	}

	// Run the main action _if_ we are not in the loop for the version command
	// that is added by default.
	if p.Before != nil {
		if err := p.Before(ctx); err != nil {
			return false, err
		}
	}

	// Run the action with the context and post-flag-processing args.
	if err := p.Action(ctx, p.FlagSet.Args()); err != nil {
		// Run the after function.
		if p.After != nil {
			p.After(ctx)
		}

		return false, err
	}

	// Run the after function.
	if p.After != nil {
		if err := p.After(ctx); err != nil {
			return false, err
		}
	}

	// Done.
	return false, nil
}

func (p *Program) usage(ctx context.Context) error {
	fmt.Fprintf(os.Stderr, "%s -  %s.\n\n", p.Name, strings.TrimSuffix(strings.TrimSpace(p.Description), "."))
	fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", p.Name)
	fmt.Fprintln(os.Stderr)

	// Print information about the common/global flags.
	if p.FlagSet != nil {
		resetFlagUsage(p.FlagSet)
	}

	// Print information about the commands.
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr)

	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	for _, command := range p.Commands {
		if !command.Hidden() {
			fmt.Fprintf(w, "\t%s\t%s\n", command.Name(), command.ShortHelp())
		}
	}
	w.Flush()

	fmt.Fprintln(os.Stderr)
	return nil
}

func (p *Program) resetCommandUsage(command Command) {
	p.FlagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s %s %s\n", p.Name, command.Name(), command.Args())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, strings.TrimSpace(command.LongHelp()))
		fmt.Fprintln(os.Stderr)
		resetFlagUsage(p.FlagSet)
	}
}

type mflag struct {
	name     string
	defValue string
	usage    string
}

// byName implements sort.Interface for []mflag based on the name field.
type byName []mflag

func (n byName) Len() int      { return len(n) }
func (n byName) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool {
	return strings.TrimPrefix(n[i].name, "-") < strings.TrimPrefix(n[j].name, "-")
}

func resetFlagUsage(fs *flag.FlagSet) {
	var (
		hasFlags   bool
		flagBlock  bytes.Buffer
		flagMap    = []mflag{}
		flagWriter = tabwriter.NewWriter(&flagBlock, 0, 4, 2, ' ', 0)
	)

	fs.VisitAll(func(f *flag.Flag) {
		hasFlags = true

		// Default-empty string vars should read "(default: <none>)"
		// rather than the comparatively ugly "(default: )".
		defValue := f.DefValue
		if defValue == "" {
			defValue = "<none>"
		}

		// Add a double dash if the name is only one character long.
		name := f.Name
		if len(name) > 1 {
			name = "-" + name
		}

		// Try and find duplicates (or the shortcode flags and combine them.
		// Like: -, --password
		for k, v := range flagMap {
			if v.usage == f.Usage {
				if len(v.name) <= 2 {
					// We already had the shortcode, let's append.
					v.name = fmt.Sprintf("%s, -%s", v.name, name)
				} else {
					v.name = fmt.Sprintf("%s, -%s", name, v.name)
				}
				flagMap[k].name = v.name

				// Return here.
				return
			}
		}

		flagMap = append(flagMap, mflag{
			name:     name,
			defValue: defValue,
			usage:    f.Usage,
		})
	})

	// Sort by name and preserve order on output.
	sort.Sort(byName(flagMap))
	for i := 0; i < len(flagMap); i++ {
		fmt.Fprintf(flagWriter, "\t-%s\t%s (default: %s)\n", flagMap[i].name, flagMap[i].usage, flagMap[i].defValue)
	}

	flagWriter.Flush()

	if !hasFlags {
		return // Return early.
	}

	fmt.Fprintln(os.Stderr, "Flags:")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, flagBlock.String())
}

func defaultFlagSet(n string) *flag.FlagSet {
	// Create the default flagset with a debug flag.
	return flag.NewFlagSet(n, flag.ExitOnError)
}

func in(a string, c []Command) bool {
	for _, b := range c {
		if b.Name() == a {
			return true
		}
	}
	return false
}

func contains(match []string, s string) bool {
	// Iterate over the items to match.
	for _, m := range match {
		if s == m {
			return true
		}
	}
	return false
}

func (p *Program) defaultContext() context.Context {
	// Create the context with the values we need to pass to the version command.
	ctx := context.WithValue(context.Background(), GitCommitKey, p.GitCommit)
	return context.WithValue(ctx, VersionKey, p.Version)
}
