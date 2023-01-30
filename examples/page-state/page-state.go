package main

import (
	"fmt"
	"html/template"
)

func main() {
	outputFormats := createOutputFormats()
	renderFormats := initRenderFormats(outputFormats)

	s := &site{
		outputFormats: outputFormats,
		renderFormats: renderFormats,
	}

	ps := &pageState{
		pageOutputs: nil,
		pageOutput:  nil,
		pageCommon:  &pageCommon{m: &pageMeta{kind: KindPage}},
	}
	ps.init(s)

	// prepare
	ps.pageOutput = ps.pageOutputs[0]

	// render
	fmt.Println(ps.targetPaths().TargetFilename)
	fmt.Println(ps.Content())
	fmt.Println(ps.m.kind)
}

type site struct {
	outputFormats map[string]Formats
	renderFormats Formats
}

type pageState struct {
	// This slice will be of same length as the number of global slice of output
	// formats (for all sites).
	pageOutputs []*pageOutput

	// This will be shifted out when we start to render a new output format.
	*pageOutput

	// Common for all output formats.
	*pageCommon
}

func (p *pageState) init(s *site) {
	pp := newPagePaths(s)
	p.pageOutputs = make([]*pageOutput, len(s.renderFormats))
	for i, f := range s.renderFormats {
		ft, found := pp.targetPaths[f.Name]
		if !found {
			panic("target path not found")
		}
		providers := struct{ targetPather }{ft}
		po := &pageOutput{
			f:                      f,
			pagePerOutputProviders: providers,
			ContentProvider:        nil,
		}
		contentProvider := newPageContentOutput(po)
		po.ContentProvider = contentProvider
		p.pageOutputs[i] = po
	}
}

func newPageContentOutput(po *pageOutput) *pageContentOutput {
	cp := &pageContentOutput{
		f: po.f,
	}
	initContent := func() {
		cp.content = template.HTML("<p>hello content</p>")
	}

	cp.initMain = func() {
		initContent()
	}
	return cp
}

func newPagePaths(s *site) pagePaths {
	outputFormats := s.renderFormats
	targets := make(map[string]targetPathsHolder)

	for _, f := range outputFormats {
		target := "/" + "blog" + "/" + f.BaseName +
			"." + f.MediaType.SubType
		paths := TargetPaths{
			TargetFilename: target,
		}
		targets[f.Name] = targetPathsHolder{
			paths: paths,
		}
	}
	return pagePaths{
		targetPaths: targets,
	}
}

type pagePaths struct {
	targetPaths map[string]targetPathsHolder
}

type targetPathsHolder struct {
	paths TargetPaths
}

func (t targetPathsHolder) targetPaths() TargetPaths {
	return t.paths
}

type pageOutput struct {
	f Format

	// These interface provides the functionality that is specific for this
	// output format.
	pagePerOutputProviders
	ContentProvider

	// May be nil.
	cp *pageContentOutput
}

// pageContentOutput represents the Page content for a given output format.
type pageContentOutput struct {
	f        Format
	initMain func()
	content  template.HTML
}

func (p *pageContentOutput) Content() any {
	p.initMain()
	return p.content
}

// these will be shifted out when rendering a given output format.
type pagePerOutputProviders interface {
	targetPather
}

type targetPather interface {
	targetPaths() TargetPaths
}

type TargetPaths struct {
	// Where to store the file on disk relative to the publish dir. OS slashes.
	TargetFilename string
}

type ContentProvider interface {
	Content() any
}

type pageCommon struct {
	m *pageMeta
}

type pageMeta struct {
	// kind is the discriminator that identifies the different page types
	// in the different page collections. This can, as an example, be used
	// to to filter regular pages, find sections etc.
	// Kind will, for the pages available to the templates, be one of:
	// page, home, section, taxonomy and term.
	// It is of string type to make it easy to reason about in
	// the templates.
	kind string
}

func initRenderFormats(
	outputFormats map[string]Formats) Formats {
	return outputFormats[KindPage]
}

func createOutputFormats() map[string]Formats {
	m := map[string]Formats{
		KindPage: {HTMLFormat},
	}

	return m
}

const (
	KindPage = "page"
)

var HTMLType = newMediaType("text", "html")

// HTMLFormat An ordered list of built-in output formats.
var HTMLFormat = Format{
	Name:      "HTML",
	MediaType: HTMLType,
	BaseName:  "index",
}

func newMediaType(main, sub string) Type {
	t := Type{
		MainType:  main,
		SubType:   sub,
		Delimiter: "."}
	return t
}

type Type struct {
	MainType  string `json:"mainType"`  // i.e. text
	SubType   string `json:"subType"`   // i.e. html
	Delimiter string `json:"delimiter"` // e.g. "."
}

type Format struct {
	// The Name is used as an identifier. Internal output formats (i.e. HTML and RSS)
	// can be overridden by providing a new definition for those types.
	Name string `json:"name"`

	MediaType Type `json:"-"`

	// The base output file name used when not using "ugly URLs", defaults to "index".
	BaseName string `json:"baseName"`
}

type Formats []Format
