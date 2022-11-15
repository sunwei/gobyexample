package overlayfs

import (
	"bytes"
	"fmt"
	"golang.org/x/tools/txtar"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
)

func main() {
	dir, _ := os.MkdirTemp("", "hugo")
	defer os.RemoveAll(dir)

	ofs := New([]AbsStatFss{
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

	fis, _ := ofs.ReadDir(".")
	for _, fi := range fis {
		fmt.Println(fi.Name())
	}
}

func projectModuleFs(dir string) AbsStatFss {
	modulePath := filepath.Join(dir, "project")
	_ = os.Mkdir(modulePath, os.ModePerm)

	ps := "-- a.md --\n" +
		"project: a\n" +
		"-- c.md --\n" +
		"project: c"

	writeFiles(ps, modulePath)

	return &moduleFs{
		workingDir: modulePath,
	}
}

func mythemeModuleFs(dir string) AbsStatFss {
	modulePath := filepath.Join(dir, "mytheme")
	_ = os.Mkdir(modulePath, os.ModePerm)

	ms := "-- a.md --\n" +
		"mytheme: a\n" +
		"-- b.md --\n" +
		"mytheme: b"

	writeFiles(ms, modulePath)

	return &moduleFs{
		workingDir: modulePath,
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

type moduleFs struct {
	workingDir string
}

func (m *moduleFs) Abs(name string) []string {
	return []string{path.Join(m.workingDir, name)}
}

func (m *moduleFs) Fss() []fs.StatFS {
	return []fs.StatFS{os.DirFS(m.workingDir).(fs.StatFS)}
}

type AbsStatFss interface {
	// Abs returns an absolute path of file or dir.
	Abs(name string) []string

	Fss() []fs.StatFS
}

// OverlayFs is a filesystem that overlays multiple filesystems.
// It's by default a read-only filesystem.
// For all operations, the filesystems are checked in order until found.
type OverlayFs struct {
	fss []AbsStatFss
}

// New creates a new OverlayFs with the given options.
func New(fss []AbsStatFss) *OverlayFs {
	return &OverlayFs{
		fss: fss,
	}
}

func (ofs OverlayFs) Append(fss ...AbsStatFss) *OverlayFs {
	ofs.fss = append(ofs.fss, fss...)
	return &ofs
}

// Open opens a file, returning it or an error, if any happens.
// If name is a directory, a *Dir is returned representing all directories matching name.
// Note that a *Dir must not be used after it's closed.
func (ofs *OverlayFs) Open(name string) (fs.File, error) {
	bfs, fi, err := ofs.stat(name)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, os.ErrInvalid
	}

	return bfs.Open(name)
}

func (ofs *OverlayFs) Stat(name string) (os.FileInfo, error) {
	_, fi, err := ofs.stat(name)
	return fi, err
}

func (ofs *OverlayFs) stat(name string) (fs.StatFS, os.FileInfo, error) {
	for _, bfs := range ofs.fss {
		fss := bfs.Fss()
		for _, sfs := range fss {
			if fi, err := sfs.Stat(name); err == nil ||
				!os.IsNotExist(err) {
				return sfs, fi, nil
			}
		}
	}
	return nil, nil, os.ErrNotExist
}

func readFile(f fs.File) string {
	b, _ := io.ReadAll(f)
	return string(b)
}

func (ofs *OverlayFs) ReadDir(dirname string) ([]fs.FileInfo, error) {
	var fis []fs.FileInfo
	readDir := func(bfs AbsStatFss) error {
		dirs := bfs.Abs(dirname)
		for _, dir := range dirs {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
			fis = merge(fis, files)
		}

		return nil
	}

	for _, bfs := range ofs.fss {
		if err := readDir(bfs); err != nil {
			return nil, err
		}
	}

	sort.Slice(fis, func(i, j int) bool {
		return fis[i].Name() < fis[j].Name()
	})

	return fis, nil
}

func merge(upper, lower []fs.FileInfo) []fs.FileInfo {
	for _, lfi := range lower {
		var found bool
		for _, ufi := range upper {
			if lfi.Name() == ufi.Name() {
				found = true
				break
			}
		}
		if !found {
			upper = append(upper, lfi)
		}
	}
	return upper
}
