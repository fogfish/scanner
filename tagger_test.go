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

func TestTagger(t *testing.T) {
	for input, expected := range map[string][]string{
		"<>Hello World!</>":         {"Hello World!"},
		"x<>Hello World!</>x":       {"Hello World!"},
		"<>Hello!</><>World.</>":    {"Hello!", "World."},
		"<>Hello!</>x<>World.</>":   {"Hello!", "World."},
		"z<>Hello!</>x<>World.</>x": {"Hello!", "World."},
		`<>Hello! World.`:           {},
		`Hello!\xWorld.<>`:          {},
		`</>Hello!\xWorld.`:         {},
		"<>Hello 3.14 World!</>":    {"Hello 3.14 World!"},
	} {
		s := scanner.NewTagger("<>", "</>", strings.NewReader(input))

		seq := make([]string, 0)
		for s.Scan() {
			seq = append(seq, s.Text())
		}

		it.Then(t).Should(
			it.Seq(seq).Equal(expected...),
		)
	}
}
