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
	"unicode"
	"unicode/utf8"
)

type Sentencer []byte

// [bufio.SplitFunc] for sentence.
func (s Sentencer) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
		}
	}

	// Scan until end of sentence [.!?]\s+|\z
	var r rune
	for width, i := 0, start; i < len(data)-1; i += width {
		r, width = utf8.DecodeRune(data[i+1:])
		if bytes.IndexByte([]byte(s), data[i]) >= 0 {
			if unicode.IsSpace(r) {
				return i + 1, data[start : i+1], nil
			}
		}
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data[start:], nil
	}

	// Request more data.
	return 0, nil, nil
}

// Create a scanner that slices input stream by end of sentence
func NewSentencer(eos string, r io.Reader) *bufio.Scanner {
	if len(eos) == 0 {
		eos = EndOfSentence
	}

	s := bufio.NewScanner(r)
	s.Split(Sentencer(eos).Split)
	return s
}

// Default end of sentence
const EndOfSentence = ".!?"
