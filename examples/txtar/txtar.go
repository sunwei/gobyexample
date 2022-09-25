package main

import (
	"bytes"
	"fmt"
	"path/filepath"

	"golang.org/x/tools/txtar"
)

// 文件结构
// 文件名: config.toml
// 文件内容：theme = 'mytheme'
var files = "-- config.toml --\n" +
	"theme = 'mytheme'"

func main() {
	// 解析上面的文件结构
	data := txtar.Parse([]byte(files))
	fmt.Println("File start:")

	// 遍历解析生成的所有文件，通过File结构体获取文件名和文件数据
	//
	//type File struct {
	//	Name string // name of file ("foo/bar.txt")
	//	Data []byte // text content of file
	//}
	for _, f := range data.Files {
		filename := filepath.Join("workingDir", f.Name)
		data := bytes.TrimSuffix(f.Data, []byte("\n"))

		fmt.Println(filename)
		fmt.Println(string(data))
	}
	fmt.Println("File end.")
}
