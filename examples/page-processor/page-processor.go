package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"strconv"
)

func main() {
	s := &set{elements: []string{}}
	p := newPagesProcessor(s)
	p.Start(context.Background())
	defer func() {
		err := p.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}()

	data := []any{1, "hello", 2, 3, 4, 5, 6, "world", "happy"}
	for _, d := range data {
		err := p.Process(d)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println(s)
}

type set struct {
	elements []string
}

func (s *set) Add(element string) {
	s.elements = append(s.elements, element)
}

type pageProcessor interface {
	Process(item any) error
	Start(ctx context.Context) context.Context
	Wait() error
}

func newPagesProcessor(s *set) *pagesProcessor {
	ps := make(map[string]pageProcessor)
	ps["i"] = &intProcessor{processor{
		s:        s,
		itemChan: make(chan interface{}, 2),
	}}
	ps["s"] = &stringProcessor{processor{
		s:        s,
		itemChan: make(chan interface{}, 2),
	}}

	return &pagesProcessor{processors: ps}
}

type pagesProcessor struct {
	processors map[string]pageProcessor
}

func (p *pagesProcessor) Process(item any) error {
	switch v := item.(type) {
	// Page bundles mapped to their language.
	case int:
		err := p.processors["i"].Process(v)
		if err != nil {
			return err
		}
	case string:
		err := p.processors["s"].Process(v)
		if err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unrecognized item type in Process: %T", item))
	}

	return nil
}
func (p *pagesProcessor) Start(ctx context.Context) context.Context {
	for _, proc := range p.processors {
		ctx = proc.Start(ctx)
	}
	return ctx
}
func (p *pagesProcessor) Wait() error {
	var err error
	for _, proc := range p.processors {
		if e := proc.Wait(); e != nil {
			err = e
		}
	}
	return err
}

type processor struct {
	s         *set
	ctx       context.Context
	itemChan  chan any
	itemGroup *errgroup.Group
}

func (p *processor) Process(item any) error {
	select {
	case <-p.ctx.Done():
		return nil
	default:
		p.itemChan <- item
	}
	return nil
}
func (p *processor) Start(ctx context.Context) context.Context {
	p.itemGroup, ctx = errgroup.WithContext(ctx)
	p.ctx = ctx
	return ctx
}
func (p *processor) Wait() error {
	close(p.itemChan)
	return p.itemGroup.Wait()
}

type intProcessor struct {
	processor
}

func (i *intProcessor) Start(ctx context.Context) context.Context {
	ctx = i.processor.Start(ctx)
	i.processor.itemGroup.Go(func() error {
		for item := range i.processor.itemChan {
			if err := i.doProcess(item); err != nil {
				return err
			}
		}
		return nil
	})
	return ctx
}

func (i *intProcessor) doProcess(item any) error {
	switch v := item.(type) {
	case int:
		i.processor.s.Add(strconv.Itoa(v))
	default:
		panic(fmt.Sprintf("unrecognized item type in intProcess: %T", item))
	}
	return nil
}

type stringProcessor struct {
	processor
}

func (i *stringProcessor) Start(ctx context.Context) context.Context {
	ctx = i.processor.Start(ctx)
	i.processor.itemGroup.Go(func() error {
		for item := range i.processor.itemChan {
			if err := i.doProcess(item); err != nil {
				return err
			}
		}
		return nil
	})
	return ctx
}

func (i *stringProcessor) doProcess(item any) error {
	switch v := item.(type) {
	case string:
		i.processor.s.Add(v)
	default:
		panic(fmt.Sprintf("unrecognized item type in stringProcessor: %T", item))
	}
	return nil
}
