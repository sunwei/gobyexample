package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/afero"
	"golang.org/x/tools/txtar"
	"os"
)

func main() {
	ofs := New([]afero.Fs{
		projectModuleFs(), mythemeModuleFs()})
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

func projectModuleFs() afero.Fs {
	ps := "-- a.md --\n" +
		"project: a\n" +
		"-- c.md --\n" +
		"project: c"
	return fsFromTxtTar(ps)
}

func mythemeModuleFs() afero.Fs {
	ms := "-- a.md --\n" +
		"mytheme: a\n" +
		"-- b.md --\n" +
		"mytheme: b"
	return fsFromTxtTar(ms)
}

func fsFromTxtTar(s string) afero.Fs {
	data := txtar.Parse([]byte(s))
	fs := afero.NewMemMapFs()
	for _, f := range data.Files {
		if err := afero.WriteFile(
			fs,
			f.Name,
			bytes.TrimSuffix(f.Data, []byte("\n")),
			0o666); err != nil {
			panic(err)
		}
	}
	return fs
}

// OverlayFs is a filesystem that overlays multiple filesystems.
// It's by default a read-only filesystem.
// For all operations, the filesystems are checked in order until found.
type OverlayFs struct {
	fss []afero.Fs
}

// New creates a new OverlayFs with the given options.
func New(fss []afero.Fs) *OverlayFs {
	return &OverlayFs{
		fss: fss,
	}
}

func (ofs *OverlayFs) stat(name string) (
	afero.Fs, os.FileInfo, bool, error) {
	for _, fs := range ofs.fss {
		if fi, err := fs.Stat(name); err == nil ||
			!os.IsNotExist(err) {
			return fs, fi, false, err
		}
	}
	return nil, nil, false, os.ErrNotExist
}

// Open opens a file, returning it or an error, if any happens.
// If name is a directory, a *Dir is returned representing all directories matching name.
// Note that a *Dir must not be used after it's closed.
func (ofs *OverlayFs) Open(name string) (
	afero.File, error) {
	fs, fi, _, err := ofs.stat(name)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		panic("dir not supported yet.")
	}

	return fs.Open(name)
}

func readFile(f afero.File) string {
	b, _ := afero.ReadAll(f)
	return string(b)
}
