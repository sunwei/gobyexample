package main

import (
	"fmt"
	"strings"
)

type configProvider struct {
	DefaultLanguage string
}

// Language manages specific-language configuration.
type Language struct {
	Lang string

	// If set per language, this tells Hugo that all content files without any
	// language indicator (e.g. my-page.en.md) is in this language.
	// This is usually a path relative to the working dir, but it can be an
	// absolute directory reference. It is what we get.
	// For internal use.
	ContentDir string

	// Global config.
	// For internal use.
	Cfg configProvider
}

// Site contains all the information relevant for constructing a static
// site.  The basic flow of information is as follows:
//
// 1. A list of Files is parsed and then converted into Pages.
//
//  2. Pages contain sections (based on the file they were generated from),
//     aliases and slugs (included in a pages frontmatter) which are the
//     various targets that will get generated.  There will be canonical
//     listing.  The canonical path can be overruled based on a pattern.
//
//  3. Taxonomies are created via configuration and will present some aspect of
//     the final page and typically a perm url.
//
//  4. All Pages are passed through a template based on their desired
//     layout based on numerous different elements.
//
// 5. The entire collection of files is written to disk.
type Site struct {
	language *Language

	// Output formats defined in site config per Page Kind, or some defaults
	// if not set.
	// Output formats defined in Page front matter will override these.
	outputFormats map[string]Formats

	// All the output formats and media types available for this site.
	// These values will be merged from the Hugo defaults, the site config and,
	// finally, the language settings.
	outputFormatsConfig Formats
	mediaTypesConfig    Types
}

func main() {
	cusCfg := configProvider{
		DefaultLanguage: "en",
	}

	lang := &Language{
		Lang:       cusCfg.DefaultLanguage,
		ContentDir: "mycontent",
		Cfg:        cusCfg,
	}

	mediaTypes := DecodeTypes()
	formats := DecodeFormats(mediaTypes)
	outputFormats := createSiteOutputFormats(formats)

	s := &Site{
		language: lang,

		outputFormats:       outputFormats,
		outputFormatsConfig: formats,
		mediaTypesConfig:    mediaTypes,
	}

	fmt.Println("Site:")
	fmt.Printf("%#v\n", s)
}

// Type (also known as MIME type and content type) is a two-part identifier for
// file formats and format contents transmitted on the Internet.
// For Hugo's use case, we use the top-level type name / subtype name + suffix.
// One example would be application/svg+xml
// If suffix is not provided, the sub type will be used.
// See // https://en.wikipedia.org/wiki/Media_type
type Type struct {
	MainType  string `json:"mainType"`  // i.e. text
	SubType   string `json:"subType"`   // i.e. html
	Delimiter string `json:"delimiter"` // e.g. "."
}

// Type returns a string representing the main- and sub-type of a media type, e.g. "text/css".
// A suffix identifier will be appended after a "+" if set, e.g. "image/svg+xml".
// Hugo will register a set of default media types.
// These can be overridden by the user in the configuration,
// by defining a media type with the same Type.
func (m Type) Type() string {
	// Examples are
	// image/svg+xml
	// text/css
	return m.MainType + "/" + m.SubType
}

// Types is a slice of media types.
type Types []Type

const defaultDelimiter = "."

var HTMLType = newMediaType("text", "html")

func newMediaType(main, sub string) Type {
	t := Type{
		MainType:  main,
		SubType:   sub,
		Delimiter: defaultDelimiter}
	return t
}

// DefaultTypes is the default media types supported by Hugo.
var DefaultTypes = Types{
	HTMLType,
}

// DecodeTypes takes a list of media type configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeTypes() Types {
	var m Types

	// remove duplications
	// Maps type string to Type. Type string is the full application/svg+xml.
	mmm := make(map[string]Type)
	for _, dt := range DefaultTypes {
		mmm[dt.Type()] = dt
	}

	for _, v := range mmm {
		m = append(m, v)
	}

	return m
}

// Format represents an output representation, usually to a file on disk.
type Format struct {
	// The Name is used as an identifier. Internal output formats (i.e. HTML and RSS)
	// can be overridden by providing a new definition for those types.
	Name string `json:"name"`

	MediaType Type `json:"-"`

	// The base output file name used when not using "ugly URLs", defaults to "index".
	BaseName string `json:"baseName"`
}

// Formats is a slice of Format.
type Formats []Format

// GetByName gets a format by its identifier name.
func (formats Formats) GetByName(
	name string) (f Format, found bool) {
	for _, ff := range formats {
		if strings.EqualFold(name, ff.Name) {
			f = ff
			found = true
			return
		}
	}
	return
}

// HTMLFormat An ordered list of built-in output formats.
var HTMLFormat = Format{
	Name:      "HTML",
	MediaType: HTMLType,
	BaseName:  "index",
}

// DefaultFormats contains the default output formats supported by Hugo.
var DefaultFormats = Formats{
	HTMLFormat,
}

// DecodeFormats takes a list of output format configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeFormats(mediaTypes Types) Formats {
	// Format could be modified by mediaTypes configuration
	// just make it simple for example
	fmt.Println(mediaTypes)

	f := make(Formats, len(DefaultFormats))
	copy(f, DefaultFormats)

	return f
}

func createSiteOutputFormats(
	allFormats Formats) map[string]Formats {
	defaultOutputFormats :=
		createDefaultOutputFormats(allFormats)
	return defaultOutputFormats
}

const (
	KindPage = "page"
	kind404  = "404"
)

func createDefaultOutputFormats(
	allFormats Formats) map[string]Formats {
	htmlOut, _ := allFormats.GetByName(HTMLFormat.Name)

	m := map[string]Formats{
		KindPage: {htmlOut},
		kind404:  {htmlOut},
	}

	return m
}
