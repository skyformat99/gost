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
	"github.com/byte16/gost/processor"
	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc [-p | --path {pathName}] -- [goFlags] [goArgs]",
	Short: "Set specified path as GOPATH and run 'go doc'",
	Long: `Set specified path as GOPATH and run 'go doc'. If no path has been
specified, gost would try to check if current working directory was under
any single path in its pathspace. And if gost found one, it would set it
as the GOPATH and then run 'go doc'. If you want to use a multiple path
like '/home/foo/bar:home/foo/baz', it's a must to use -p or --path to
specify the pathName defined in pathspace.

	[-p | --path {pathName}]
		It's used to specify the path to set as GOPATH before running the
		corresponding go command.
	[goFlags] [goArgs]
		These content following '--' would be passed directly to the underlying
		go command.

example: gost doc -p foobar -- github.com/baz/qux`,
	Run: func(cmd *cobra.Command, args []string) {
		processor.RunGoCmd(cmd.Name(), args)
	},
}

func init() {
	RootCmd.AddCommand(docCmd)
	docCmd.SetHelpFunc(processor.RunGoHelp)
}
