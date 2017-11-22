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

package command

import (
	"github.com/byte16/gost/common"
	"github.com/spf13/cobra"
)

// bugCmd represents the bug command
var bugCmd = &cobra.Command{
	Use:   "bug [-p | --path {pathName}]",
	Short: "Set specified path as GOPATH and run 'go bug'.",
	Long: `Set specified path as GOPATH and run 'go bug'. If no path has been
specified, gost would try to check if current working directory was under
any single path in its pathspace. And if gost found one, it would set it
as the GOPATH and then run 'go bug'. If you want to use a multiple path
like '/home/foo/bar:home/foo/baz', it's a must to use -p or --path to
specify the pathName defined in pathspace.

	[-p | --path {pathName}]
		It's used to specify the path to set as GOPATH before running the
		go command.

example: gost bug -p foobar
`,
	Run: func(cmd *cobra.Command, args []string) {
		common.RunGoCmd(cmd.Name(), args)
	},
}

func init() {
	RootCmd.AddCommand(bugCmd)
	bugCmd.SetHelpFunc(common.RunGoHelp)
}
