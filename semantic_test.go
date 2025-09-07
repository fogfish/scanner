//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/embeddings
//

package scanner_test

import (
	"context"
	"strings"
	"testing"

	"github.com/fogfish/it/v2"
	"github.com/fogfish/scanner"
)

func TestScanner(t *testing.T) {
	text := "a. bb. c. ddd. ff."

	s := scanner.NewSemantic(
		embed{},
		scanner.NewSentencer(scanner.EndOfSentence, strings.NewReader(text)),
	)
	s.Similarity(similar)
	s.Window(3)

	it.Then(t).Should(
		it.True(s.Scan()),
		it.Seq(s.Text()).Equal("a.", "c."),
		it.True(s.Scan()),
		it.Seq(s.Text()).Equal("bb.", "ff."),
		it.True(s.Scan()),
		it.Seq(s.Text()).Equal("ddd."),
	)

	it.Then(t).ShouldNot(
		it.True(s.Scan()),
	)
}

//------------------------------------------------------------------------------

type embed struct{}

func (embed) UsedTokens() int { return 0 }
func (embed) Embedding(ctx context.Context, text string) ([]float32, int, error) {
	return []float32{float32(len(text))}, 0, nil
}

func similar(a, b []float32) bool { return a[0] == b[0] }
