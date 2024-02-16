package utils

import (
	"procinspect/pkg/semantic"
)

func ParseSql(source string) (*semantic.Script, error) {
	parser := NewParallelParser(source)
	return parser.Parse()
}
