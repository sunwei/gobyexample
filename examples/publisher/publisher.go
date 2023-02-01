package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// publisher needs to know:
	// 1: what to publish
	// 2: where to publish

	// 1
	// src is template executed result
	// it is the source that we need to publish
	// take a look at template executor example
	// https://c.sunwei.xyz/template-executor.html
	src := &bytes.Buffer{}
	src.Write([]byte("template executed result"))

	b := &bytes.Buffer{}
	transformers := createTransformerChain()
	if err := transformers.Apply(b, src); err != nil {
		fmt.Println(err)
		return
	}

	dir, _ := os.MkdirTemp("", "hugo")
	defer os.RemoveAll(dir)

	// 2
	// targetPath is from pageState
	// this is where we need to publish
	// take a look at page state example
	// https://c.sunwei.xyz/page-state.html
	targetPath := filepath.Join(dir, "index.html")

	if err := os.WriteFile(
		targetPath,
		bytes.TrimSuffix(b.Bytes(), []byte("\n")),
		os.ModePerm); err != nil {
		panic(err)
	}

	fmt.Println("1. what to publish: ", string(b.Bytes()))
	fmt.Println("2. where to publish: ", dir)
}

func (c *Chain) Apply(to io.Writer, from io.Reader) error {
	fb := &bytes.Buffer{}
	if _, err := fb.ReadFrom(from); err != nil {
		return err
	}

	tb := &bytes.Buffer{}

	ftb := &fromToBuffer{from: fb, to: tb}
	for i, tr := range *c {
		if i > 0 {
			panic("switch from/to and reset to")
		}
		if err := tr(ftb); err != nil {
			continue
		}
	}
	_, err := ftb.to.WriteTo(to)
	return err
}

func createTransformerChain() Chain {
	transformers := NewEmpty()
	transformers = append(transformers, func(ft FromTo) error {
		content := ft.From().Bytes()
		w := ft.To()
		tc := bytes.Replace(
			content,
			[]byte("result"), []byte("transferred result"), 1)
		_, _ = w.Write(tc)
		return nil
	})
	return transformers
}

// Chain is an ordered processing chain. The next transform operation will
// receive the output from the previous.
type Chain []Transformer

// Transformer is the func that needs to be implemented by a transformation step.
type Transformer func(ft FromTo) error

// FromTo is sent to each transformation step in the chain.
type FromTo interface {
	From() BytesReader
	To() io.Writer
}

// BytesReader wraps the Bytes method, usually implemented by bytes.Buffer, and an
// io.Reader.
type BytesReader interface {
	// Bytes The slice given by Bytes is valid for use only until the next buffer modification.
	// That is, if you want to use this value outside of the current transformer step,
	// you need to take a copy.
	Bytes() []byte

	io.Reader
}

// NewEmpty creates a new slice of transformers with a capacity of 20.
func NewEmpty() Chain {
	return make(Chain, 0, 2)
}

// Implements contentTransformer
// Content is read from the from-buffer and rewritten to to the to-buffer.
type fromToBuffer struct {
	from *bytes.Buffer
	to   *bytes.Buffer
}

func (ft fromToBuffer) From() BytesReader {
	return ft.from
}

func (ft fromToBuffer) To() io.Writer {
	return ft.to
}
