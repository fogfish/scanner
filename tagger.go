//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/scanner
//

package scanner

import (
	"bufio"
	"bytes"
	"io"
)

type Tagger struct {
	open, close []byte
}

// [bufio.SplitFunc] for sentence.
func (t Tagger) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	openIndex := bytes.Index(data, t.open)
	if openIndex < 0 {
		if atEOF {
			return len(data), nil, nil
		}
		return 0, nil, nil
	}

	start := openIndex + len(t.open)
	closeIndex := bytes.Index(data[start:], t.close)
	if closeIndex < 0 {
		if atEOF {
			return len(data), nil, nil
		}
		return 0, nil, nil
	}

	end := start + closeIndex
	advance = end + len(t.close)
	token = data[start:end]
	return advance, token, nil
}

// Create a scanner that slices input stream by end of sentence
func NewTagger(open, close string, r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(Tagger{open: []byte(open), close: []byte(close)}.Split)
	return s
}
