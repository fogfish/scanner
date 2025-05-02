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

type Slicer []byte

// [bufio.SplitFunc] for fixed delimiter.
func (s Slicer) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	delim := []byte(s)

	if i := bytes.Index(data, delim); i >= 0 {
		return i + len(delim), data[:i], nil
	}

	if atEOF {
		if len(data) > 0 {
			return len(data), data, nil
		}
	}

	return 0, nil, nil
}

// Create a scanner that slices input stream by fixed delimiter
func NewSlicer(delim string, r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(Slicer(delim).Split)
	return s
}
