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

// rmCmd represents the rmpath command
var rmCmd = &cobra.Command{
	Use:   "rm <pathName...>",
	Short: "Remove paths specified by pathName args from gost pathspace.",
	Long: `Remove paths specified by pathName args from gost pathspace. The pathNames
should have have been defined in gost's pathspace with 'gost add'. It accepts
at least one pathName as argument.

	[pathName...]
		These pathNames specified which paths would be removed by gost.

example: gost rm foo bar`,
	Args:             cobra.MinimumNArgs(1),
	PersistentPreRun: common.Nop,
	Run: func(cmd *cobra.Command, args []string) {
		common.RmPath(args)
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
	rmCmd.SetHelpFunc(common.RunCommonHelp)
}
