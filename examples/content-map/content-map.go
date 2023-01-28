package main

import (
	"bytes"
	"fmt"
	"github.com/sunwei/gobyexample/modules/radixtree"
	"golang.org/x/tools/txtar"
	"os"
	"path/filepath"
)

func main() {
	dir, _ := os.MkdirTemp("", "hugo")
	defer os.RemoveAll(dir)

	myContentPath := filepath.Join(dir, "blog")
	_ = os.Mkdir(myContentPath, os.ModePerm)

	mcs := "-- post.md --\n" +
		"---\n" +
		"title: \"Post Title\"\n" +
		"---\n" +
		"### first blog\n" +
		"Hugo & Caddy > WordPress & Apache"
	writeFiles(mcs, myContentPath)

	bundlePath := filepath.Join(myContentPath, "post")
	f, err := os.Stat(bundlePath + ".md")
	if err != nil {
		fmt.Println(err)
		return
	}
	n := &contentNode{
		fi:   f,
		path: "blog/post.md",
	}

	baseKey := myContentPath + "/"
	m := newContentMap()
	m.sections.Insert(baseKey, &contentNode{
		path: "blog",
	})

	key := baseKey + cmBranchSeparator + "post" + cmLeafSeparator
	m.pages.Insert(key, n)

	v, _ := m.sections.Get(baseKey)
	fmt.Println("=--=", baseKey, v)

	v, _ = m.pages.Get(key)
	fmt.Println("=--=", key, v)
	cn := v.(*contentNode)
	fmt.Println(cn.fi.Name(), cn.fi.Size())
}

const (
	cmBranchSeparator = "__hb_"
	cmLeafSeparator   = "__hl_"
)

type contentNode struct {
	// Set if source is a file.
	// We will soon get other sources.
	fi os.FileInfo

	// The source path. Unix slashes. No leading slash.
	path string
}

type contentTree struct {
	Name string
	*radixtree.Tree
}

type contentMap struct {
	pages    *contentTree
	sections *contentTree
}

func newContentMap() *contentMap {
	return &contentMap{
		pages:    &contentTree{Name: "pages", Tree: radixtree.New()},
		sections: &contentTree{Name: "sections", Tree: radixtree.New()},
	}
}

func writeFiles(s string, dir string) {
	data := txtar.Parse([]byte(s))

	for _, f := range data.Files {
		if err := os.WriteFile(
			filepath.Join(dir, f.Name),
			bytes.TrimSuffix(f.Data, []byte("\n")),
			os.ModePerm); err != nil {
			panic(err)
		}
	}
}
