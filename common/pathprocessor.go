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
	"github.com/byte16/gost/cfg"
	"path/filepath"
	"strings"
)

// Try to get a valid GOPATH from pathspace
func GetGoPath(pathName string) (string, error) {
	var path string
	// when empty pathName given, current WD should be take
	// into consider.
	if pathName == "" {
		path = getGoPathByWD()
	} else {
		path = getGoPathByName(pathName)
	}

	if path == "" {
		return "", fmt.Errorf("failed to get GOPATH")
	}
	return path, nil
}

// Try to get GOPATH to use by current working directory
func getGoPathByWD() string {
	for _, p := range cfg.PathList {
		if strings.HasPrefix(cfg.WD, filepath.Join(p.RealPath, "src")) {
			return p.Path
		}
	}
	return ""
}

// Try to get GOPATH to use by name which has been listed in %HOME/.gost.yml
func getGoPathByName(pathName string) string {
	if p, ok := cfg.PathMap[pathName]; ok {
		return p.Path
	}
	return ""
}

// Add given path info into config
func AddPath(pathName, path string, multi bool) error {
	// Check if path and -m flag are consistent
	isMulti := IsMultiPath(path)
	if isMulti != multi {
		if isMulti {
			fmt.Print(`[FAIL] if you want to add a multiple path(a path list) which would be
used as GOPATH once selected in the future, please use -m(or --multi) flag.
`)
		} else {
			fmt.Printf(`[FAIL] if want to add single path(not a path list) which would be
used as GOPATH once selected in the future, please don't' use -m(or --multi)
flag.
`)
		}
		return fmt.Errorf("inconsistent path and -m(--multi) flag")
	}

	realPath := path
	if !multi {
		// When path is not a path list, try to get the real path which
		// the path points to, in case it contains symbolic links. This
		// is just a preparation for path match when executing go commands
		// without explicitly specifying which path item to use.
		tmpPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			fmt.Printf(`[FATAL] failed to resolve given path: %s
%s
`,
				path, err.Error())
			return err
		}
		realPath = tmpPath
	}

	pathItem, exist := cfg.PathMap[pathName]
	if exist && pathItem != nil {
		fmt.Printf(`[WARNING] path name '%s' has exist in pathspace:
  name:     %s
  path:     %s
  realPath: %s
  multi:    %v
It will be overwritten to:
  name:     %s
  realPath: %s
  path:     %s
  multi:    %v
`,
			pathName,
			pathItem.Name, pathItem.Path, pathItem.RealPath, pathItem.Multi,
			pathName, path, realPath, multi)

		pathItem.Path = path
		pathItem.RealPath = realPath
		pathItem.Multi = multi
	} else {
		pathItem = &cfg.GoPath{
			Name:     pathName,
			Path:     path,
			RealPath: realPath,
			Multi:    multi,
		}
		cfg.PathMap[pathName] = pathItem
	}

	cfg.WriteTomlCfg()
	return nil
}

// Check if a given path is a single path but not a
// path list.
func IsMultiPath(path string) bool {
	if strings.Contains(path, string(filepath.ListSeparator)) {
		return true
	}
	return false
}

// Remove path item from config file according to path name.
func RmPath(pathNames []string) {
	orgnQty := len(cfg.PathMap)
	for _, pathName := range pathNames {
		_, exist := cfg.PathMap[pathName]
		if !exist {
			fmt.Printf(`[WARNING] path item '%s' doesn't exist in pathspace, please check.
`,
				pathName)
			continue
		}
		delete(cfg.PathMap, pathName)
	}

	if orgnQty == len(cfg.PathMap) {
		fmt.Println("[INFO] no valid path item has been removed")
		return
	}

	cfg.WriteTomlCfg()
}

// List the path info whose name has been list as arguments.
// If there is no argument, list all the info of the path items.
func ListPath(pathNames []string) {
	if len(cfg.PathMap) == 0 {
		fmt.Printf(`[INFO] there is no path item in the pathspace`)
		return
	}

	if len(pathNames) == 0 {
		for _, pathItem := range cfg.PathList {
			printPath(&pathItem)
		}
		return
	}

	invalids := []string{}
	for _, pathName := range pathNames {
		pathItem, exist := cfg.PathMap[pathName]
		if !exist {
			invalids = append(invalids, pathName)
			continue
		}
		printPath(pathItem)
	}

	// List out the paths not exist
	if len(invalids) > 0 {
		fmt.Println("[WARNING] the items below are not in pathspace:")
		for _, item := range invalids {
			fmt.Printf("   %s\n", item)
		}
	}
}

func printPath(pathItem *cfg.GoPath) {
	fmt.Printf(
		` * Name:     %s
   Path:     %s
   RealPath: %s
   Multi:    %v
`,
		pathItem.Name, pathItem.Path, pathItem.RealPath, pathItem.Multi)
}
