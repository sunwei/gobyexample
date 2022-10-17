package main

import (
	"bytes"
	"fmt"
	toml "github.com/pelletier/go-toml/v2"
	"golang.org/x/tools/txtar"
	"path/filepath"
	"strings"
)

// 文件结构
// 文件名: config.toml
// 文件内容：theme = 'mytheme'
var files = "-- config.toml --\n" +
	"theme = 'mytheme'"

// Format 文件格式类型
type Format string

// TOML 支持的格式，为简单示例，只支持TOML格式
const (
	TOML Format = "toml"
)

func main() {
	// 解析上面的文件结构
	data := txtar.Parse([]byte(files))
	fmt.Println("File start:")

	// Input: 数据，格式，输出类型
	var configData []byte
	var format Format
	m := make(map[string]any)

	// 遍历解析生成的所有文件，通过File结构体获取文件名和文件数据
	// f.Name 获取文件名
	// f.Data 获取文件数据
	// 如果是config.toml文件，则获取文件数据
	for _, f := range data.Files {
		if "config.toml" == f.Name {
			configData = bytes.TrimSuffix(
				f.Data, []byte("\n"))
			format = FormatFromString(f.Name)
		}
	}

	err := UnmarshalTo(configData, format, &m)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m)
	}

	fmt.Println("File end.")
}

// FormatFromString turns formatStr, typically a file extension without any ".",
// into a Format. It returns an empty string for unknown formats.
// Hugo 实现
func FormatFromString(formatStr string) Format {
	formatStr = strings.ToLower(formatStr)
	if strings.Contains(formatStr, ".") {
		// Assume a filename
		formatStr = strings.TrimPrefix(
			filepath.Ext(formatStr), ".")
	}
	switch formatStr {
	case "toml":
		return TOML
	}

	return ""
}

// UnmarshalTo unmarshals data in format f into v.
func UnmarshalTo(data []byte, f Format, v any) error {
	var err error

	switch f {
	case TOML:
		err = toml.Unmarshal(data, v)

	default:
		return fmt.Errorf(
			"unmarshal of format %q is not supported", f)
	}

	if err == nil {
		return nil
	}

	return err
}
