# Gost

[中文](https://github.com/byte16/gost/blob/master/README_zh.md)

`gost` is a **Simple Tool** for Go which can help you to manage GOPATHs in an incredibly easy and convenient  way. It's very useful for you if you are working on multiple projects in Go at the same time based on isolated workspaces for reasons (for example, these projects have some same dependencies of diffrent versions) and tired of switching among different GOPATHs to run Go commands. It can be used on CI/CD tool chain or just a CommandLine compensation or partner for IDEs like Gogland which supports GOPATH management on project level(Gogland is unable to switch GOPATHs for CommandLine terminals).



`gost` is a wrapper of Go commands and can help you to run them with the GOPATH set as what you specify in real time. `gost` maintains a repository of paths with a toml file named `.gost.toml` under the executing user's home directory. We can call this repository **pathspace** and it can be managed with some well designed commands like `gost add`, `gost ls` and `gost rm`. You can also specify a `.gost.toml`s at other place by using `--config` each time you run `gost`, but it's not recommended for the possibilities to cause confusion.



## Installation

Download and install the package with the following command:

```
$ go get github.com/byte16/gost
```

This will create the gost executable under your `$GOPATH/bin` directory. You can move it to a permanent directory under `$PATH` so that you can use it without caring what the `$GOPATH` would be changed to.

You can also download the gost executable directly from the release page:

https://github.com/byte16/gost/releases



## Getting Started

### Manage pathspace

Any path that you want to use as GOPATH should be put into pathspace first. You can add a path like this:

```
$ gost add foo /home/foobar/bar
```

 `gost` would ask you to give each path item a name when you add them into pathspace with `gost add`. In the example above, `foo` is just the name for the path `/home/foobar/bar `. With the name you can specify which path to use as GOPATH  when running a go command with `gost`. The name should be unique in the pathspace. If you want to name a new path item with a name that has been used to name other path item which has exists when you run `gost add`, the old path item would be overwritten.

**Note**: Before you add a path into pathspace, you should first make sure that the directory that the path represents exists. `gost` supports symbloic links well.



`gost` supports both single path like `/home/foo/bar` and multiple path(a path list) like `/home/foobar/baz:/home/foobar/quz`(an example for Linux). Path list separator differs in different operating systems. Make sure you are using the right one. When you want to add a multiple path to the pathspace, it's **required** to use `--multi`(or `-m` for short) flag to tell `gost` explicitly. Here is an example:

```
$ gost add foobar /home/foobar/foo:/home/foobar/bar -m
```



You can use `gost ls` to see the details of the paths in pathspace. If you are sure that you would never use any path in the future, you can remove it from the pathspace with `gost rm`. For example:

```
$ gost ls
 * Name:     foo
   Path:     /home/foobar/bar
   RealPath: /home/foobar/bar 	#if Path contains symbolic links, RealPath would be the solid one.
   Multi:    false
 * Name:     foobar
   Path:     /home/foobar/foo:/home/foobar/bar
   RealPath: /home/foobar/foo:/home/foobar/bar 	#it would be resolved by gost automatically.
   Multi:    true
$ gost ls foo	# gost ls support one or more names as arguments.
 * Name:     foo
   Path:     /home/foobar/bar
   RealPath: /home/foobar/bar
   Multi:    false
$ gost rm foobar # gost rm also support one or more names as arguments.
$ gost ls 
 * Name:     foo
   Path:     /home/foobar/bar
   RealPath: /home/foobar/bar
   Multi:    false
```



### Run Go commands

It's incredibly simple to use `gost` to run a go command with a given path as GOPATH. Just remember that the path should have been added into pathspace. For example, you have added `/home/foobar/bar` into pathspace and named it with `foo` just like what you did above. Now you want to run go command below with GOPATH set as `/home/foobar/bar`:

```
$ go get -u github.com/byte16/gost
```

Normally before without `gost`, you need to set `/home/foobar/bar` as GOPATH with running `export GOPATH=/home/foobar/bar`  or adding it into `$HOME/.bash_profile` and `source` it first. Now with `gost` you don't have to do that. You just need to run the command below:

```
$ gost get -p foo -- -u github.com/byte16/gost
```

The jobs of switching GOPATH and running `go get` would all be took care of by `gost`. You would never experience any differences between the two ways except **the simpleness and convenience of using `gost`**. It could be even more simpler if `/home/foobar/bar` has existed in gost's pathspace as a single path and you are currently under `/home/foobar/bar/src` or any sub directories of it, you can just run `gost get` command without specifying `foo` with `--path` or `-p` like:

```
$ gost get -- -u github.com/byte16/gost
```

`gost` would detect that itself.

But if you want to run a go command with a multiple path as GOPATH like `/home/foobar/foo:/home/foobar/bar`, the `--path`(or `-p` for short) should never be omitted.



By the way, the double dash symbol `--` is used to tell `gost` that the flags and arguments after the double dash symbol `--` should be passed to the underlying go command directly. So you should never omit the symbol.



Now `gost` supports all the functional commands of Go and each corresponding `gost`
command has the same name with its underlying Go command. They are listed below:

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
        vet         run go tool vet on packages
For the details about how to use each of them with `gost`, just run `gost help cmdName`. For example:

```
$ gost help install
usage: gost install [-p | --path {pathName}] -- [goFlags] [goArgs]

Set specified path as GOPATH and run 'go install'. If no path has been
specified, gost would try to check if current working directory was under
any single path in its pathspace. And if gost found one, it would set it
as the GOPATH and then run 'go install'. If you want to use a multiple path
like '/home/foo/bar:home/foo/baz', it's a must to use -p or --path to
specify the pathName defined in pathspace.

        [-p | --path {pathName}]
                It's used to specify the path to set as GOPATH before running the
                corresponding go command.
        [goFlags] [goArgs]
                These content following '--' would be passed directly to the underlying
                go command.

example: gost install -p foobar -- -n github.com/baz/qux


*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*
| Below is the help info of corresponding go command: |
*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*--*


usage: go install [build flags] [packages]

Install compiles and installs the packages named by the import paths,
along with their dependencies.

For more about the build flags, see 'go help build'.
For more about specifying packages, see 'go help packages'.

See also: go build, go get, go clean.

```



## Contributing

Contributions are greatly appreciated. If you find an issue or want to contribute please file an [issue](https://github.com/byte16/gost/issues) or create a [pull request](https://github.com/byte16/gost/pulls).



## License

Gost is released under the MIT license. See [LICENSE](https://github.com/byte16/gost/blob/master/LICENSE)