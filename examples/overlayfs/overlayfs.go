package main

import (
	"bytes"
	"fmt"
	"golang.org/x/tools/txtar"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	dir, _ := os.MkdirTemp("", "hugo")
	defer os.RemoveAll(dir)

	ofs := New([]fs.StatFS{
		projectModuleFs(dir), mythemeModuleFs(dir)})
	a, _ := ofs.Open("a.md")
	defer a.Close()
	b, _ := ofs.Open("b.md")
	defer b.Close()
	c, _ := ofs.Open("c.md")
	defer c.Close()

	fmt.Println(readFile(a))
	fmt.Println(readFile(b))
	fmt.Println(readFile(c))
}

func projectModuleFs(dir string) fs.StatFS {
	modulePath := filepath.Join(dir, "project")
	_ = os.Mkdir(modulePath, os.ModePerm)

	ps := "-- a.md --\n" +
		"project: a\n" +
		"-- c.md --\n" +
		"project: c"

	writeFiles(ps, modulePath)

	return os.DirFS(modulePath).(fs.StatFS)
}

func mythemeModuleFs(dir string) fs.StatFS {
	modulePath := filepath.Join(dir, "mytheme")
	_ = os.Mkdir(modulePath, os.ModePerm)

	ms := "-- a.md --\n" +
		"mytheme: a\n" +
		"-- b.md --\n" +
		"mytheme: b"

	writeFiles(ms, modulePath)

	return os.DirFS(modulePath).(fs.StatFS)
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

// OverlayFs is a filesystem that overlays multiple filesystems.
// It's by default a read-only filesystem.
// For all operations, the filesystems are checked in order until found.
type OverlayFs struct {
	fss []fs.StatFS
}

// New creates a new OverlayFs with the given options.
func New(fss []fs.StatFS) *OverlayFs {
	return &OverlayFs{
		fss: fss,
	}
}

// Open opens a file, returning it or an error, if any happens.
// If name is a directory, a *Dir is returned representing all directories matching name.
// Note that a *Dir must not be used after it's closed.
func (ofs *OverlayFs) Open(name string) (fs.File, error) {
	for _, fs := range ofs.fss {
		if _, err := fs.Stat(name); err == nil ||
			!os.IsNotExist(err) {
			return fs.Open(name)
		}
	}
	return nil, nil
}

func readFile(f fs.File) string {
	b, _ := io.ReadAll(f)
	return string(b)
}
