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

package cfg

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	// The original config info
	GostToml Cfg

	// The list of usable GOPATHs in original order
	PathList []GoPath

	// A GOPATH container for quick access
	PathMap = map[string]*GoPath{}

	// Current working directory.
	WD string

	// Current user's home dir
	UserHome string

	// Path of config file .gost.toml
	ConfigPath string
)

// Name of config file
const CONFIG = ".gost.toml"

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("[FATAL] unable to get info of current working directory: %s\n", err.Error())
		os.Exit(1)
	}

	// If wd contains symbolic links, get the real one
	// it points to. This is just a preparation for path
	// match when executing go commands without explicitly
	// specifying which path item to use.
	realWd, err := filepath.EvalSymlinks(wd)
	if err != nil {
		fmt.Printf("[FATAL] failed to eval info of current working directory: %s\n", err.Error())
		os.Exit(1)
	}
	WD = realWd

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	UserHome = home

	// Config file may have been set a custom one
	if ConfigPath == "" {
		ConfigPath = filepath.Join(home, CONFIG)
	}
}

type Cfg struct {
	GoPaths []GoPath `toml:"GoPaths"`
}

type GoPath struct {
	// It's a nick-name for the given GOPATH
	Name string `toml:"name,required"`

	// The GOPATH
	Path string `toml:"path,required"`

	// Normally RealPath is the same with Path, unless
	// Path contains symbolic links. When Path contains
	// symbolic links, RealPath would be the real path
	// that Path finally points to.
	RealPath string `toml:"realPath"`

	// This denotes that if the GOPATH is multiple.
	// (true for yes and false for no)
	Multi bool `toml:"multi"`
}

func LoadTomlCfg() {
	err := viper.Unmarshal(&GostToml)
	if err != nil {
		fmt.Printf("[FATAL] failed to load .god.yml: %s\n", err.Error())
		os.Exit(1)
	}

	for i, pathItem := range GostToml.GoPaths {
		if !validateCfg(&pathItem) {
			fmt.Printf(`[FATAL] invalid GoPath item found in %s, pleas check:
  - Name: %s
  - Path: %s
  - RealPath: %s
`,
				ConfigPath, pathItem.Name, pathItem.Path, pathItem.RealPath)
			os.Exit(1)
		}
		PathMap[pathItem.Name] = &GostToml.GoPaths[i]
	}

	PathList = GostToml.GoPaths
}

// Try to check if a given GoPath is valid
// (true for yes and false for no)
func validateCfg(path *GoPath) bool {
	if path == nil ||
		path.Name == "" ||
		path.Path == "" ||
		path.RealPath == "" {
		return false
	}
	return true
}

// Write config info into file
func WriteTomlCfg() {
	f, e := os.OpenFile(ConfigPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if e != nil {
		fmt.Printf(`[FATAL] failed to open or create config file %s:
  %s
`,
			ConfigPath, e.Error())
		os.Exit(1)
	}
	defer f.Close()

	newPathList := make([]GoPath, len(PathMap))
	i := 0
	for _, pathItem := range PathMap {
		newPathList[i] = *pathItem
		i++
	}
	PathList = newPathList

	cfg := Cfg{GoPaths: PathList}
	cfgWriter := toml.NewEncoder(f)
	err := cfgWriter.Encode(&cfg)
	if err != nil {
		fmt.Println(err)
	}
}
