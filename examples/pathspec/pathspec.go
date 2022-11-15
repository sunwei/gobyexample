package main

import (
	"bytes"
	"fmt"
	"github.com/sunwei/gobyexample/modules/overlayfs"
	"github.com/sunwei/gobyexample/modules/radixtree"
	"golang.org/x/tools/txtar"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func main() {
	dir, _ := os.MkdirTemp("", "hugo")
	defer os.RemoveAll(dir)

	myContentPath := filepath.Join(dir, "mycontent")
	_ = os.Mkdir(myContentPath, os.ModePerm)
	myContent2Path := filepath.Join(dir, "mycontent2")
	_ = os.Mkdir(myContent2Path, os.ModePerm)
	themePath := filepath.Join(dir, "mytheme")
	_ = os.Mkdir(themePath, os.ModePerm)

	mcs := "-- a.md --\n" +
		"mycontent: a\n" +
		"-- c.md --\n" +
		"mycontent: c"
	writeFiles(mcs, myContentPath)

	mc2s := "-- a.md --\n" +
		"mycontent2: a\n" +
		"-- d.md --\n" +
		"mycontent2: d"
	writeFiles(mc2s, myContent2Path)

	ts := "-- a.md --\n" +
		"mytheme: a\n" +
		"-- b.md --\n" +
		"mytheme: b"

	writeFiles(ts, themePath)

	p := &Paths{
		WorkingDir: dir,
		AllModules: []Module{
			{
				ProjectMod: true,
				Dir:        dir,
				Mounts: []Mount{
					{Source: "mycontent",
						Target: "content"},
					{Source: "mycontent2",
						Target: "content"},
				},
			},
			{
				ProjectMod: false,
				Dir:        dir,
				Mounts: []Mount{
					{Source: "mytheme",
						Target: "content"},
				},
			},
		},
	}

	collector := &filesystemsCollector{
		overlayMountsContent: overlayfs.New(
			[]overlayfs.AbsStatFss{}),
	}
	createOverlayFs(collector, p)

	var f fs.File
	fis, _ := collector.
		overlayMountsContent.ReadDir(filepathSeparator)
	for _, fi := range fis {
		fmt.Println(fi.Name())
		f, _ = collector.
			overlayMountsContent.Open(fi.Name())
		b, _ := io.ReadAll(f)
		fmt.Println(string(b))
	}
	defer f.Close()
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

var filepathSeparator = string(filepath.Separator)

// RootMapping describes a virtual file or directory mount.
type RootMapping struct {
	// The virtual mount.
	From string
	// The source directory or file.
	To string
	// The base of To. May be empty if an
	// absolute path was provided.
	ToBasedir string
	// Whether this is a mount in the main project.
	IsProject bool
	// The virtual mount point, e.g. "blog".
	path string
}

type Mount struct {
	Source string
	Target string
}

type Module struct {
	ProjectMod bool
	Mounts     []Mount
	Dir        string
}

type Modules []Module

type Paths struct {
	AllModules Modules
	WorkingDir string
}

// A RootMappingFs maps several roots into one.
// Note that the root of this filesystem
// is directories only, and they will be returned
// in Readdir and Readdirnames
// in the order given.
type RootMappingFs struct {
	fs            overlayfs.AbsStatFss
	rootMapToReal *radixtree.Tree
}

type filesystemsCollector struct {
	overlayMountsContent *overlayfs.OverlayFs
}

func createOverlayFs(
	collector *filesystemsCollector,
	path *Paths) {

	for _, md := range path.AllModules {
		var fromToContent []RootMapping
		for _, mount := range md.Mounts {
			rm := RootMapping{
				From:      mount.Target, // content
				To:        mount.Source, // mycontent
				ToBasedir: md.Dir,
				IsProject: md.ProjectMod,
			}
			fromToContent = append(fromToContent, rm)
		}
		rmfsContent := newRootMappingFs(fromToContent...)
		collector.overlayMountsContent = collector.
			overlayMountsContent.Append(rmfsContent)
	}
	return
}

// NewRootMappingFs creates a new RootMappingFs
// on top of the provided with root mappings with
// some optional metadata about the root.
// Note that From represents a virtual root
// that maps to the actual filename in To.
func newRootMappingFs(
	rms ...RootMapping) *RootMappingFs {
	t := radixtree.New()
	var virtualRoots []RootMapping

	for _, rm := range rms {
		key := filepathSeparator + rm.From
		mappings := getRms(t, key)
		mappings = append(mappings, rm)
		t.Insert(key, mappings)

		virtualRoots = append(virtualRoots, rm)
	}

	t.Insert(filepathSeparator, virtualRoots)

	return &RootMappingFs{
		rootMapToReal: t,
	}
}

func (m *RootMappingFs) Abs(name string) []string {
	mappings := getRms(m.rootMapToReal, name)

	var paths []string
	for _, m := range mappings {
		paths = append(paths, path.Join(
			m.ToBasedir, m.To))
	}
	return paths
}

func (m *RootMappingFs) Fss() []fs.StatFS {
	mappings := getRms(
		m.rootMapToReal, filepathSeparator)

	var fss []fs.StatFS
	for _, m := range mappings {
		fss = append(fss, os.DirFS(
			path.Join(m.ToBasedir, m.To)).(fs.StatFS))
	}
	return fss
}

func getRms(t *radixtree.Tree,
	key string) []RootMapping {
	var mappings []RootMapping
	v, found := t.Get(key)
	if found {
		mappings = v.([]RootMapping)
	}
	return mappings
}
