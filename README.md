# 深入理解Hugo

代码实例来帮助 [深入理解Hugo](https://hugo.sunwei.xyz),

### Overview

深入理解Hugo 站点并不是由`hugo`生成的。

站点生成的主要思想是一样的：
`templates` + `content - examples` => `site - public`
既通过模板，将内容转化为静态站点。

虽然`hugo`提供了强大的模板功能，但主要支持的[内容格式](https://gohugo.io/content-management/formats/)是markdown。
并不是我们`examples`中的`*.go`文件。

再加上解析方式的不同。
在我们的示例中，需要逐行解析，分别用`blackfriday`解析注解，和`chroma`来解析go代码。
这不同于Hugo现在用`goldmark`来解析md文件，并将文件内容作为整体来进行处理，虽然提供了一些字符处理功能函数。
`readDir`可以帮我们获取样例下的所有文件，`readFile`可以帮助我们读取文本内容，`replaceRE`可以帮助我们用正则来处理内容。
但没法生成高亮的go代码块。

如果想创建一个Hugo Theme来提供这些功能，做起来就会比较别扭。
最好的方式是拓展Hugo功能，让它支持编程语言的内容格式。
但这会和Hugo现在的内容体系显得格格不入，有些别扭。
说不定可以重新开启一个Hugo的分支，专门用来支持开发人员源码讲解的站点需求。
让我们在这里加一个**TODO**，放到深入理解Hugo的实际样例章节 - 开发人员友好的Hugo。

### Building

[![test](https://github.com/sunwei/gobyexample/actions/workflows/test.yml/badge.svg)](https://github.com/sunwei/gobyexample/actions/workflows/test.yml)

安装Go，运行:

```console
$ tools/build
```

支持所见既所得，持续构建:

```console
$ tools/build-loop
```

启动本地服务:

```
$ tools/serve
```

可检查站点： `http://127.0.0.1:8000/` 

### Publishing

持续发布:

```console
$ git push
```
