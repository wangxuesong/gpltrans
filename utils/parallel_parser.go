package utils

import (
	"errors"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"procinspect/pkg/parser"
	"procinspect/pkg/semantic"
)

type (
	ParseRequest struct {
		Source string
		Index  int
		Start  int
	}
	ParseResult struct {
		Index   int
		Start   int
		AstFunc func(int) (*semantic.Script, error)
		Error   error
		Source  string
	}

	ParallelParser struct {
		workers *WorkerPool

		source        string
		updateHandler func(delta int, total int)
		requests      []*ParseRequest
		total         int
		results       []*ParseResult
	}
)

func NewParallelParser(src string) *ParallelParser {
	return &ParallelParser{source: src}
}

func (p ParallelParser) WithUpdateHandler(handler func(delta, total int)) *ParallelParser {
	p.updateHandler = handler
	return &p
}

func (p *ParallelParser) Parse() (*semantic.Script, error) {
	p.prepareRequest()

	p.results = make([]*ParseResult, len(p.requests))
	p.parseSyntax()

	if err := p.processSyntaxError(); err != nil {
		return nil, err
	}

	return p.parseSemantic()
}

func (p *ParallelParser) parseSyntax() {
	numWorkers := runtime.GOMAXPROCS(0) - 1
	// TODO: clean work pool
	pool := NewWorkerPool(numWorkers, len(p.requests))
	parseChan := make(chan *ParseResult)
	for _, req := range p.requests {
		tmpReq := *req
		pool.Submit(func() {
			parseChan <- parseBlock(&tmpReq)
		})
	}
	for _, _ = range p.requests {
		result := <-parseChan
		p.results[result.Index] = result
		if p.updateHandler != nil {
			p.updateHandler(strings.Count(result.Source, "\n"), p.total)
		}
	}
	close(parseChan)
}

func (p *ParallelParser) prepareRequest() {
	p.requests = make([]*ParseRequest, 0)
	re := regexp.MustCompile(`\r\n`)
	source := p.source
	source = re.ReplaceAllString(source, "\n")
	// split source by /
	regex := regexp.MustCompile(`(?m)^/$`)
	blocks := regex.Split(source, -1)
	start := 0
	offset := 0
	for i, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		p.requests = append(p.requests, &ParseRequest{
			Source: block,
			Index:  i,
			Start:  start + offset,
		})
		start += strings.Count(block, "\n")
		offset = 0
	}
	p.total = start
}

func (p *ParallelParser) parseSemantic() (*semantic.Script, error) {
	script := &semantic.Script{}
	var ee []error
	for _, result := range p.results {
		if result.AstFunc == nil {
			continue
		}
		s, err := result.AstFunc(result.Start)
		if err != nil {
			joinErr, ok := err.(interface{ Unwrap() []error })
			if ok {
				errs := joinErr.Unwrap()
				for _, e := range errs {
					ee = append(ee, e)
				}
			} else {
				ee = append(ee, err)
			}
		}
		s = fixLineNumber(s, result.Start)
		script = appendScript(script, s)
	}
	return script, errors.Join(ee...)
}

func (p *ParallelParser) processSyntaxError() error {
	var err error
	for _, result := range p.results {
		if result.Error != nil {
			joinErr, ok := result.Error.(interface{ Unwrap() []error })
			if ok {
				errs := joinErr.Unwrap()
				for _, e := range errs {
					err = errors.Join(err, e)
				}
			} else {
				err = errors.Join(err, result.Error)
			}
		}
	}
	return err
}

type fixLineVisitor struct {
	semantic.StubNodeVisitor
	start int
}

func newFixLineVisitor(start int) *fixLineVisitor {
	v := &fixLineVisitor{start: start}
	v.StubNodeVisitor.NodeVisitor = v
	return v
}

func (v *fixLineVisitor) VisitChildren(node semantic.AstNode) (err error) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		node = nil
		return
	}
	for _, child := range semantic.GetChildren(node) {
		_ = child.Accept(v)
	}
	if n, ok := node.(interface {
		semantic.Node
		semantic.SetPosition
	}); ok {
		if n != nil {
			n.SetLine(n.Line() + v.start)
		}
	}
	return
}
func fixLineNumber(script *semantic.Script, start int) *semantic.Script {
	v := newFixLineVisitor(start)
	script.Accept(v)
	return script
}
func appendScript(script *semantic.Script, s *semantic.Script) *semantic.Script {
	for _, stmt := range s.Statements {
		script.Statements = append(script.Statements, stmt)
	}
	return script
}
func parseBlock(r *ParseRequest) *ParseResult {
	result := &ParseResult{
		Index: r.Index,
		Start: r.Start,
	}
	result.Source = r.Source
	result.AstFunc, result.Error = parser.ParseSql(r.Source)

	return result
}
