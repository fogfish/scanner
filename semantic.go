//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/scanner
//

package scanner

import (
	"context"
	"fmt"
)

// Semantic provides a convenient solution for semantic chunking.
// Successive calls to the Semantic.Scan method will step through the context
// windows of a file and grouping sentences semantically. The context window
// is defined by number sentences, use Window method to change
// default 32 sentences value.
//
// The specification of a sentence is defined by the Scanner interface, which
// is compatible with [bufio.NewScanner]. Use a Split function of type SplitFunc
// within [bufio.NewScanner] to control sentence breakdown.
//
// The package provides [NewSentencer] utility that breaks the input into
// sentences using punctuation runes. Redefine Use Split function of
// [bufio.NewScanner] to define own algorithms.
//
// The scanner uses embeddings to determine similarity. Use Similarity method
// to change the default high cosine similarity to own implementation.
// The module provides high, medium, weak and dissimilarity functions based on
// cosine distance.
//
// Scanning stops unrecoverably at EOF or the first I/O error.
type Semantic struct {
	embed                 Embedder
	confSimilarity        func([]float32, []float32) bool
	confWindowInSentences int
	confSimilarityWith    SimilarityWith
	scanner               Scanner
	err                   error
	eof                   bool
	window                []vector
	cursor                []string
}

type vector struct {
	text string
	vf32 []float32
}

// Creates new instance of Scanner to read from io.Reader and using embedding.
func NewSemantic(embed Embedder, r Scanner) *Semantic {
	return &Semantic{
		embed:                 embed,
		confSimilarity:        HighSimilarity,
		confWindowInSentences: 32,
		confSimilarityWith:    SIMILARITY_WITH_TAIL,
		scanner:               r,
		window:                make([]vector, 0),
	}
}

// Similarity sets the similarity function for the Semantic.
// The default is HighSimilarity.
func (s *Semantic) Similarity(f func([]float32, []float32) bool) {
	s.confSimilarity = f
}

// Similarity sets the behavior to sorting algorithms.
//
// Using SIMILARITY_WITH_HEAD configures algorithm to sort chunk similar
// to the first element of chunk. The first element of chunk is stable during
// the chunk forming.
//
// Using SIMILARITY_WITH_TAIL configures algorithm to sort chunk similar
// to the last element of chunk. The last element is changed after new one is added to chunk.
func (s *Semantic) SimilarityWith(x SimilarityWith) {
	s.confSimilarityWith = x
}

// Widow defines the context window for similarity detection.
// The default value is 32 sentences.
func (s *Semantic) Window(n int) {
	s.confWindowInSentences = n
}

func (s *Semantic) Err() error     { return s.err }
func (s *Semantic) Text() []string { return s.cursor }

// Scan advances the Semantic through context window, sequences will be available
// through [Semantic.Text]. It returns false if there was I/O error or EOF is reached.
func (s *Semantic) Scan() bool {
	if s.err != nil {
		return false
	}

	if !s.eof {
		s.eof, s.err = s.fill()
		if s.err != nil {
			return false
		}
	}

	s.cursor = s.peek()

	return !(s.eof && len(s.cursor) == 0)
}

// fill the window
func (s *Semantic) fill() (bool, error) {
	wn := s.confWindowInSentences - len(s.window)
	for wn > 0 && s.scanner.Scan() {
		txt := s.scanner.Text()
		v32, _, err := s.embed.Embedding(context.Background(), txt)
		if err != nil {
			return false, fmt.Errorf("embedding has failed: %w, for {%s}", err, txt)
		}

		s.window = append(s.window, vector{text: txt, vf32: v32})
		wn--
	}

	if err := s.scanner.Err(); err != nil {
		return false, err
	}

	return wn != 0, nil
}

// peek similar from the window
func (s *Semantic) peek() []string {
	if len(s.window) == 0 {
		return nil
	}

	// split the window into similar (a) and non-similar (b) items
	a, b := make([]vector, 0), make([]vector, 0)
	a = append(a, s.window[0])

	for i := 1; i < len(s.window); i++ {
		var at int
		switch s.confSimilarityWith {
		case SIMILARITY_WITH_HEAD:
			at = 0
		case SIMILARITY_WITH_TAIL:
			at = len(a) - 1
		}
		ref := a[at]

		if s.confSimilarity(ref.vf32, s.window[i].vf32) {
			a = append(a, s.window[i])
		} else {
			b = append(b, s.window[i])
		}
	}

	s.window = b

	seq := make([]string, len(a))
	for i, x := range a {
		seq[i] = x.text
	}
	return seq
}
