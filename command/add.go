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
	"os"
)

var multi bool

// addCmd represents the addpath command
var addCmd = &cobra.Command{
	Use:   "add <pathName> <path> [-m | --multi]",
	Short: "Add a path to pathspace.",
	Long: `Add a path to the pathspace maintained by gost.

	<pathName>
		It's a given name to be used to represent <path>. It's unique within the pathspace.
		When one is trying to add a new path with a name existed in the pathspace, the original
		path info would be overwritten.
	<path>
		It's the path to be add to the pathspace and it would be set as GOPATH once selected
		when running a go command. It can be a single absolute path or multiple separated by
		the os's legal PATH separator(':' in Darwin for example). When it's multiple, a '-m'
		or '--multi' flag should be marked.
	[-m | --multi]
		It's used to mark that the <path> is multiple like '/home/foo/bar:/home/foo/baz'

example: gost add foobar '/home/foo/bar:/home/foo/baz' -m`,
	Args:             cobra.ExactArgs(2),
	PersistentPreRun: common.Nop,
	Run: func(cmd *cobra.Command, args []string) {
		err := common.AddPath(args[0], args[1], multi)
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(&multi, "multi", "m", false, "path field would contains multiple paths")
	addCmd.SetHelpFunc(common.RunCommonHelp)
}
