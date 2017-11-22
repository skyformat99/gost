// Copyright © 2017 Franco Li <uint64@yeah.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// Find the path represented by pathName in the pathspace and
// then set it as the GOPATH.
func SetGoPath(pathName string) {
	goPath, err := GetGoPath(pathName)
	if err != nil {
		fmt.Println(`[FATAL] Unable to get valid GOPATH.
The GOPATH you want to use must have been defined in gost's pathspace.
Specify the name of the GOPATH which you want to use with --path (or -p for
short). When the name is not specified explicitly, gost would try to match
the current wording directory with paths in gost's pathspace (check with
'gost ls'). If there is no path item in gost's pathspace, you can add one
with 'gost add' (use 'gost help add' for detail info).
`)
		os.Exit(1)
	}

	err = os.Setenv("GOPATH", goPath)
	if err != nil {
		fmt.Printf("[FATAL] failed to set GOPATH=%s: %s\n", goPath, err.Error())
		os.Exit(1)
	}

	fmt.Printf("[INFO] gost will use %s as GOPATH\n", goPath)
}

// Run go command named cmdName and pass the args to it.
func RunGoCmd(cmdName string, args []string) {
	args = append([]string{cmdName}, args...)
	c := exec.Command("go", args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fmt.Printf("[WARNING] %s\n", err.Error())
		os.Exit(1)
	}
}

// First print the help info of 'gost cmd' and then
// run go help to get the help info of 'go cmd'.
func RunGoHelp(cmd *cobra.Command, args []string) {
	RunCommonHelp(cmd, args)

	fmt.Println(`

*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*
| Below is the help info of corresponding go command: |
*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*

`)

	name := cmd.Name()
	c := exec.Command("go", "help", name)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fmt.Printf("[WARNING] %s\n", err.Error())
		os.Exit(1)
	}
}

func RunCommonHelp(cmd *cobra.Command, args []string) {
	cmd.Usage()
	fmt.Println(cmd.Long)
}

func Nop(cmd *cobra.Command, args []string) {
}

// It's used to define the format of usage info of a command.
func CommonUsage(cmd *cobra.Command) error {
	fmt.Printf(`usage: gost %s

`,
		cmd.Use)
	return nil
}
