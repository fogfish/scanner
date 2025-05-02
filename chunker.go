package scanner

import (
	"strings"
)

type Chunker struct {
	Scanner
	size int
	sbuf strings.Builder
}

func NewChunker(size int, s Scanner) *Chunker {
	return &Chunker{
		Scanner: s,
		size:    size,
	}
}

func (s *Chunker) Scan() bool {
	s.sbuf.Reset()

	for s.Scanner.Scan() {
		s.sbuf.WriteString(s.Scanner.Text())
		if s.sbuf.Len() > s.size {
			return true
		}
	}

	return s.sbuf.Len() > 0
}

func (s *Chunker) Text() string { return s.sbuf.String() }
