//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/scanner
//

package scanner_test

import (
	"strings"
	"testing"

	"github.com/fogfish/it/v2"
	"github.com/fogfish/scanner"
)

func TestSlicer(t *testing.T) {
	for input, expected := range map[string][]string{
		"Hello World!":         {"Hello World!"},
		"Hello!!World.":        {"Hello", "World."},
		"Hello!!World!!3.14!!": {"Hello", "World", "3.14"},
	} {
		s := scanner.NewSlicer("!!", strings.NewReader(input))

		seq := make([]string, 0)
		for s.Scan() {
			seq = append(seq, s.Text())
		}

		it.Then(t).Should(
			it.Seq(seq).Equal(expected...),
		)
	}
}
