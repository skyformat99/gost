# Gost

[English](https://github.com/byte16/gost/blob/master/README.md)

`gost` 是一个用于Go语言的简单工具。它可以以一种简洁轻便的方式来帮助你管理GOPATH并运行Go命令。如果你同时开发多个Go项目，并且因为某些原因这些项目基于相互隔离的不同工作空间（例如这些项目基于某些相同依赖的不同版本），这导致你在运行Go命令的时候需要在不同的GOPATH间来回切换，那么`gost`对于你来说将非常有用。它可以用于持续开发和持续集成的工具链中，也可以作为Gogland这样可以在项目层级来管理GOPATH的集成开发工具在命令行方面的补偿（Gogland无法为命令行终端切换GOPATH）。



`gost` 是对Go命令的一个封装。它可以帮你实时设置GOPATH并运行Go命令。 `gost` 维护了一个路径的仓库，这个仓库是基于一个位于当前用户Home目录下的名为`.gost.toml`的toml文件建立的。 我们可以把这个路径仓库叫做**pathspace** ，可以使用像`gost add`， `gost ls` 和 `gost rm`这些被精心设计的命令来对它进行管理。你也可以在每次运行`gost`命令的时候通过`--config`来指定使用一个位于其它位置的`.gost.toml`文件，但是并不推荐这么做，因为这样可能会导致混乱。 


## 安装

通过下面的命令来下载并安装`gost`：

```
$ go get github.com/byte16/gost
```

这将在你的`$GOPATH/bin`目录下生成一个gost可执行文件。你可以将它转移到一个位于`$PATH`中的永久目录，这样不论`$GOPATH`被修改成什么，你都可以正常使用`gost`。

你也可以在本项目的分发页面中直接下载gost的可执行文件：

https://github.com/byte16/gost/releases



## 使用入门

### 管理 pathspace

任何你想要用作GOPATH的路径都需要首先被放入pathspace。你可以像这样来添加一个路径：

```
$ gost add foo /home/foobar/bar
```

当你通过`gost add`向pathspace中添加路径的时候，`gost`会要求你为所添加的路径取一个名称。 在上面的示例中，`foo`就是为路径`/home/foobar/bar `所取的名称。在使用`gost`来运行一个go命令的时候，你可以通过所取的名称来指定使用哪个路径作为GOPATH。每个名称在pathspace中都应该是唯一的。如果一个名称已经在pathspace中用于命名一个路径，当在运行`gost add`命令时你将这个名称再用于命名新的路径的时候，那么pathspace中的原有路径条目信息将会被覆盖。

**注意**: 在你把一个路径添加到pathspace之前，你需要首先确认这个路径所代表的目录确实存在。`gost`对包含软链接的路径支持良好。



`gost` 支持像`/home/foo/bar`这样的单路径和像`/home/foobar/baz:/home/foobar/quz`（这是Linux下的一个示例）这样的路径列表。在不同的操作系统中，路径列表分隔符是不同的。请确保你使用了正确的路径列表分隔符。当你需要把一个包含路径列表的路径条目添加到pathspace中的时候，需要使用`--multi`（或简写作`-m`）来显式告知`gost`。下面是一个示例：

```
$ gost add foobar /home/foobar/foo:/home/foobar/bar -m
```



你可以使用 `gost ls`来查看pathspace中的路径条目的详细信息。如果你确定某些路径条目已经无用，你可以使用`gost rm`来将它们从pathspace中移除。下面是操作示例：

```
$ gost ls
 * Name:     foo
   Path:     /home/foobar/bar
   RealPath: /home/foobar/bar 	# 如果Path包含软链接，RealPath将会是对应的实际路径
   Multi:    false
 * Name:     foobar
   Path:     /home/foobar/foo:/home/foobar/bar
   RealPath: /home/foobar/foo:/home/foobar/bar 	# RealPath是gost自动解析出来的 
   Multi:    true
$ gost ls foo	                # gost ls 也支持一个或多个名称作为参数
 * Name:     foo
   Path:     /home/foobar/bar
   RealPath: /home/foobar/bar
   Multi:    false
$ gost rm foobar                # gost rm 也支持一个或多个名称作为参数
$ gost ls 
 * Name:     foo
   Path:     /home/foobar/bar
   RealPath: /home/foobar/bar
   Multi:    false
```



### 运行Go命令

使用`gost`设置GOPATH并运行一条Go命令操作非常简单。只需要记住你用来设置GOPATH的路径已经被添加到了pathspace。比如上面的示例，你已经把`/home/foobar/bar`添加到了pathspace中并命名为`foo`，现在你想要把`/home/foobar/bar`设置为GOPATH并运行下面的Go命令:

```
$ go get -u github.com/byte16/gost
```

通常在以前没有使用`gost`的时候，你首先需要通过执行`export GOPATH=/home/foobar/bar`命令或者将该命令添加到`$HOME/.bash_profile`并执行`source $HOME/.bash_profile`来设置GOPATH。现在使用`gost`你无须像上面那样操作。 你只需要运行下面的命令即可：:

```
$ gost get -p foo -- -u github.com/byte16/gost
```

切换GOPATH和运行`go get`命令的任务完全由`gost`来进行执行。除了使用`gost`所带来的简便，你体会不到这两种方式之间在功能上会有什么差异。然而，你还可以使操作更加简单。如果`/home/foobar/bar`已经作为一条单目录条目被添加到pathspace中了，并且你当前的工作目录位于`/home/foobar/bar/src`路径下的任意层级的子目录，你可以运行`gost get`而无需通过`--path`或`-p`来指定`foo`，操作如下：

```
$ gost get -- -u github.com/byte16/gost
```

`gost` 会自动探测使用pathspace中的哪个路径来作为GOPATH。

但是，如果你想要将一条包含路径列表的路径条目（如`/home/foobar/foo:/home/foobar/bar`）来设置为GOPATH并运行Go命令的话，`--path`(或简写作`-p`)标志不能被省略。



顺便提示，双横线标记`--`是用来告诉`gost` `--`标记后的标志和参数需要被直接传递给底层的Go命令。所以请不要省略该标记。



现在`gost` 支持Go的所有功能性命令，对应的`gost`命令与底层的Go命令有相同的名称。这些命令如下所示：

        build       构建包和依赖
        clean       移除对象文件
        doc         显示包或者标识符的文档
        env         打印Go环境信息
        bug         开始报告bug
        fix         对包运行go工具fix
        fmt         对包的源文件运行gofmt
        generate    通过处理源代码来生成Go文件
        get         下载并安装包和依赖
        install     编译并安装包和依赖
        list        列出包列表
        run         编译并运行go程序
        test        对包进行测试
        tool        运行指定的go工具
        vet         对包运行go工具vet
关于如何通过`gost`来运行这些命令的详细细节， 可以运行`gost help 命令名称`来查看。 例如：

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



## 贡献代码

如果您能向该项目贡献代码，维护者将不胜感激。当您发现问题或者想要贡献代码，可以生成一个 [issue](https://github.com/byte16/gost/issues) 或一个[pull request](https://github.com/byte16/gost/pulls).



## 协议

Gost基于MIT 协议进行分发。 详见[协议](https://github.com/byte16/gost/blob/master/LICENSE)