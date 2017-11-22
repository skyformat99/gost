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
	"fmt"
	"os"

	"github.com/byte16/gost/cfg"
	"github.com/byte16/gost/common"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var (
	pathName string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "<subCommand> [flags... {args...}] -- [goFlags...] [goArgs...]",
	Short: "A Simple Tool but useful to help you to manage GOPATHs and run go commands",
	Long: `Gost is a Simple Tool which can help you to manage GOPATHs.
It's very useful for you if you are working on multiple projects at the
same time based on isolated workspaces for reasons and tired of switching
among different GOPATHs. It can be used on CI/CD tool chain or just a
CommandLine compensation or partner for IDEs like Gogland which supports
GOPATH management on project level.

Gost maintains a repository of paths. You can use 'gost add', 'gost rm'
and 'gost ls' to manage it. By default, gost would create a toml file
named '.gost.toml' under the executing user's home directory to record
info of path items. You can also specify one at other place by using
--config each time you run gost, but it's not recommended.

Gost would ask you to give each path item a name when you add them into
gost's pathspace with 'gost add'. With the name you can specify which
path to use as GOPATH to run go a command. The name should be unique in
gost's pathspace. If you want to name a new path item with a name that
has been used to name other path item which has exists when you run
'gost add', the old path item would be overwritten.

Gost supports both single path like '/home/foo/bar' and multiple path(a
path list) like '/home/foobar/baz:/home/foobar/quz'(example for Linux).
Path list separator differs in different operating systems. Make sure
you are using the right one. When you want to add a multiple path to the
pathspace, it's required to use --multi(or -m for short) flag to tell gost
explicitly.

It's incredibly simple to use gost to run a go command with a given path
as GOPATH. Just remember that the path should have been added into pathspace.
For example, you have add '/home/foo/bar' into pathspace with the name 'foobar'.
Now you want to run go command below with GOPATH set as '/home/foo/bar':
	go get -u github.com/byte16/gost
You can just run line below with gost:
	gost get -p foobar -- -u github.com/byte16/gost
The jobs of switching GOPATH and running 'go get' would all be took care of by
gost. You would never experience any differences between the two ways except
the simpleness and convenience of using gost. It could be even more simpler if
'/home/foo/bar' has existed in gost's pathspace as a single path and you are
currently under '/home/foo/bar' or any sub directories of it, you can just run
'gost get' command without specifying foobar with --path or -p like:
	gost get -- -u github.com/byte16/gost
Gost would detect that itself.

But if you want to run a go command with a multiple path as GOPATH like
'/home/foo/bar:/home/foo/baz', the --path(or -p for short) should never be omitted.
By the way, the double dash symbol '--' is used to tell gost that the flags and
arguments after it should be passed to the underlying go command directly. Unless
you have no flag or argument for the gost command, you should never omit it.

Now gost supports all the functional commands of go and each corresponding gost
command has the same name with its underlying go command. They are listed below:
        build       compile packages and dependencies
        clean       remove object files
        doc         show documentation for package or symbol
        env         print Go environment information
        bug         start a bug report
        fix         run go tool fix on packages
        fmt         run gofmt on package sources
        generate    generate Go files by processing source
        get         download and install packages and dependencies
        install     compile and install packages and dependencies
        list        list packages
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         run go tool vet on packages
For the details about how to use each of them with gost, just run 'gost help cmdName'.

If you want to get the source code or report a bug, please visit
https://github.com/byte16/gost
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "help" {
			return
		}
		common.SetGoPath(pathName)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.gost.yml)")

	RootCmd.PersistentFlags().StringVarP(&pathName, "path", "p", "",
		"the name defined in $HOME/.gost.yml for the GOPATH which you want to set before running the command")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.SetUsageFunc(common.CommonUsage)
	RootCmd.SetHelpFunc(common.RunCommonHelp)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		cfg.ConfigPath = cfgFile
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gost" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gost")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		cfg.LoadTomlCfg()
	} else if cfgFile!="" {
		fmt.Printf(`[FATAL] There seems to be something wrong about config file %s:
%s
Please fix it first.
`,
			cfg.ConfigPath, err.Error())
		os.Exit(1)
	}
}
