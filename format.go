package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/urfave/cli"
)

type file struct {
	currentLine  int
	name         string
	originalFile []byte
}

func format(c *cli.Context) error {
	err := forFile(c, func(f string, nodes []*parser.Node) error {
		og, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		df := file{
			currentLine:  1,
			name:         f,
			originalFile: og,
		}

		var result string
		for _, n := range nodes {
			r, err := df.doFmt(n)
			if err != nil {
				return err
			}
			result += r
		}

		// display the diff if requested
		if c.Bool("diff") {
			d, err := diff(og, []byte(result))
			if err != nil {
				return fmt.Errorf("computing diff: %s", err)
			}
			fmt.Printf("diff %s dockfmt/%s\n", f, f)
			os.Stdout.Write(d)
		}

		if c.Bool("list") {
			if !bytes.Equal(og, []byte(result)) {
				fmt.Fprintln(os.Stdout, f)
			}
		}

		// write to the file
		if c.Bool("write") {
			// make a temporary backup before overwriting original
			bakname, err := backupFile(f+".", og, 0644)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(f, []byte(result), 0644); err != nil {
				os.Rename(bakname, f)
				return err
			}

			if err := os.Remove(bakname); err != nil {
				return fmt.Errorf("could not remove backup file %s: %v", bakname, err)
			}
		}

		if !c.Bool("diff") && !c.Bool("list") && !c.Bool("write") {
			os.Stdout.Write([]byte(result))
		}

		return nil
	})

	return err
}

func (df *file) doFmt(ast *parser.Node) (result string, err error) {
	// check if we are on the correct line,
	// otherwise get the comments we are missing
	if df.currentLine != ast.StartLine {
		comments, err := df.getOriginalLines(df.currentLine, ast.StartLine, df.name)
		if err != nil {
			return "", err
		}
		result += comments
	}

	// set the variables for the directive (k) and the value (v)
	k := ast.Value
	var v string
	if ast.Next != nil {
		v = ast.Next.Value
	}

	// capitalize the directive
	k = strings.ToUpper(k)

	// format per directive
	switch k {
	case "ADD":
		v = fmtCopy(ast.Next)
	case "CMD":
		v, err = fmtCmd(ast.Next)
		if err != nil {
			return "", err
		}
	case "COPY":
		v = fmtCopy(ast.Next)
	case "ENTRYPOINT":
		v, err = fmtCmd(ast.Next)
		if err != nil {
			return "", err
		}
	case "RUN":
		v = fmtRun(v)
	default:
		v = fmtCopy(ast.Next)
	}

	// print to the result
	result = fmt.Sprintf("%s\t%s\n", k, v)

	// set our current line as the start line in the next node
	// since we want the next node
	df.currentLine++
	if ast.Next != nil {
		df.currentLine = ast.Next.StartLine
	}
	return
}

func (df *file) getOriginalLines(s int, e int, fn string) (string, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(df.originalFile))
	scanner.Split(bufio.ScanLines)
	var (
		i = 1
		l string
	)
	for scanner.Scan() {
		if i >= s && i < e {
			l += scanner.Text() + "\n"
		}
		i++
	}

	return l, nil
}

func getCmd(n *parser.Node, cmd []string) []string {
	if n == nil {
		return cmd
	}
	cmd = append(cmd, n.Value)
	if len(n.Flags) > 0 {
		cmd = append(cmd, n.Flags...)
	}
	if n.Next != nil {
		for node := n.Next; node != nil; node = node.Next {
			cmd = append(cmd, node.Value)
			if len(node.Flags) > 0 {
				cmd = append(cmd, node.Flags...)
			}
		}
	}
	return cmd
}

func fmtCmd(node *parser.Node) (string, error) {
	cmd := []string{}
	cmd = getCmd(node, cmd)
	b, err := json.Marshal(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func fmtCopy(node *parser.Node) string {
	cmd := []string{}
	cmd = getCmd(node, cmd)
	return strings.Join(cmd, "\t")
}

func fmtRun(s string) string {
	s = strings.Replace(s, "apk update && apk add", "apk add --no-cache", -1)

	var r string
	cmds := strings.Split(s, "&&")
	cmds = trimAll(cmds)
	for i, c := range cmds {
		c = strings.Replace(c, "apk --no-cache add", "apk add", -1)

		// handle `apk add` commands
		if strings.HasPrefix(c, "apk add") {
			c = strings.TrimPrefix(c, "apk add")
			// remove --no-cache
			// we will add it back later
			c = strings.Replace(c, "--no-cache", "", -1)
			c = strings.Replace(c, "apk add", "", -1)
			// recreate the command
			c = "apk add --no-cache \\" + "\n" + splitLinesWord(c)
		}

		// handle `apt-get install` commands
		if strings.HasPrefix(c, "apt-get install") {
			c = strings.TrimPrefix(c, "apt-get install")
			// remove -y
			// we will add it back later
			c = strings.Replace(c, "-y", "", -1)
			c = strings.Replace(c, "apt-get install", "", -1)
			// recreate the command
			c = "apt-get install -y \\" + "\n" + splitLinesWord(c)
		}

		// we aren't on the first line add back the `&&`
		if i != 0 {
			c = "\t&& " + c
		}

		// if we aren't on the last line add a `\\n`
		if i != len(cmds)-1 {
			c += " \\\n"
		}
		r += c
	}

	// put `apt-get update && apt-get install` on one-line it's prettier
	r = strings.Replace(r, "apt-get update \\\n\t&& apt-get install", "apt-get update && apt-get install", -1)

	return r
}

func trimAll(a []string) []string {
	for i, v := range a {
		a[i] = strings.TrimSpace(v)
	}
	return a
}

func splitLinesWord(s string) string {
	a := strings.Fields(s)
	a = trimAll(a)

	var r string
	for i, v := range a {
		r += "\t" + v
		// if we aren't on the last line add a `\\n`
		if i != len(a)-1 {
			r += " \\\n"
		}
	}
	return r
}

func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "dockfmt")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "dockfmt")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}

	return
}

const chmodSupported = runtime.GOOS != "windows"

// backupFile writes data to a new file named filename<number> with permissions perm,
// with <number randomly chosen such that the file name is unique. backupFile returns
// the chosen file name.
func backupFile(filename string, data []byte, perm os.FileMode) (string, error) {
	// create backup file
	f, err := ioutil.TempFile(filepath.Dir(filename), filepath.Base(filename))
	if err != nil {
		return "", err
	}

	bakname := f.Name()
	if chmodSupported {
		err = f.Chmod(perm)
		if err != nil {
			f.Close()
			os.Remove(bakname)
			return bakname, err
		}
	}

	// write data to backup file
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	if err1 := f.Close(); err == nil {
		err = err1
	}

	return bakname, err
}
