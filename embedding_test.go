//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/scanner
//

package scanner_test

import (
	"testing"

	"github.com/fogfish/scanner"
)

func TestHighSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float32
		expected bool
	}{
		{
			name:     "identical vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{1.0, 1.0, 1.0, 1.0},
			expected: true, // cosine distance = 0.0
		},
		{
			name:     "very similar vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{0.99, 0.99, 0.99, 0.99},
			expected: true, // cosine distance ≈ 0.02
		},
		{
			name:     "moderately similar vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{0.5, 0.5, 0.5, 0.5},
			expected: true, // cosine distance = 0.0 (same direction, different magnitude)
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{0.0, 1.0, 0.0, 0.0},
			expected: false, // cosine distance = 0.5
		},
		{
			name:     "opposite vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{-1.0, -1.0, -1.0, -1.0},
			expected: false, // cosine distance = 1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.HighSimilarity(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("HighSimilarity(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMediumSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float32
		expected bool
	}{
		{
			name:     "identical vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{1.0, 1.0, 1.0, 1.0},
			expected: false, // cosine distance = 0.0 (too similar)
		},
		{
			name:     "medium similar vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{0.5, 0.5, 0.5, 0.5},
			expected: false, // cosine distance = 0.0 (too similar)
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{0.0, 1.0, 0.0, 0.0},
			expected: true, // cosine distance = 0.5 (within range)
		},
		{
			name:     "different magnitude similar direction",
			a:        []float32{2.0, 2.0, 2.0, 2.0},
			b:        []float32{1.0, 1.0, 1.0, 1.0},
			expected: false, // cosine distance = 0.0 (same direction, different magnitude)
		},
		{
			name:     "weakly similar vectors",
			a:        []float32{1.0, 0.3, 0.3, 0.3},
			b:        []float32{0.3, 1.0, 0.3, 0.3},
			expected: false, // cosine distance in high similarity range
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.MediumSimilarity(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("MediumSimilarity(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestWeakSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float32
		expected bool
	}{
		{
			name:     "medium similar vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{0.5, 0.5, 0.5, 0.5},
			expected: false, // cosine distance ≈ 0.5 (too similar)
		},
		{
			name:     "weakly similar vectors",
			a:        []float32{1.0, 0.1, 0.1, 0.1},
			b:        []float32{0.1, 1.0, 0.1, 0.1},
			expected: false, // cosine distance in medium range
		},
		{
			name:     "vectors with weak similarity",
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{-0.1, 0.9, 0.0, 0.0},
			expected: true, // cosine distance within weak range
		},
		{
			name:     "boundary case - exactly at weak range start",
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{-0.3, 0.7, 0.0, 0.0},
			expected: true, // cosine distance within weak range
		},
		{
			name:     "almost opposite vectors",
			a:        []float32{1.0, 0.1, 0.1, 0.1},
			b:        []float32{-0.5, 0.5, 0.5, 0.5},
			expected: true, // cosine distance within weak range
		},
		{
			name:     "dissimilar vectors",
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{-1.0, 0.0, 0.0, 0.0},
			expected: false, // cosine distance = 1.0 (too dissimilar)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.WeakSimilarity(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("WeakSimilarity(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDissimilar(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float32
		expected bool
	}{
		{
			name:     "identical vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{1.0, 1.0, 1.0, 1.0},
			expected: false, // cosine distance = 0.0 (too similar)
		},
		{
			name:     "opposite vectors",
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{-1.0, -1.0, -1.0, -1.0},
			expected: true, // cosine distance = 1.0
		},
		{
			name:     "nearly opposite vectors",
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{-0.9, 0.0, 0.0, 0.0},
			expected: true, // cosine distance ≈ 0.95
		},
		{
			name:     "weakly similar vectors",
			a:        []float32{1.0, 0.1, 0.1, 0.1},
			b:        []float32{0.1, 1.0, 0.1, 0.1},
			expected: false, // cosine distance within weak range
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.Dissimilar(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Dissimilar(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestRangeSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		lo, hi   float32
		a, b     []float32
		expected bool
	}{
		{
			name:     "custom range [0.0, 0.3] - identical vectors",
			lo:       0.0,
			hi:       0.3,
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{1.0, 1.0, 1.0, 1.0},
			expected: true, // cosine distance = 0.0
		},
		{
			name:     "custom range [0.0, 0.3] - medium vectors",
			lo:       0.0,
			hi:       0.3,
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{0.5, 0.5, 0.5, 0.5},
			expected: true, // cosine distance = 0.0 (within range)
		},
		{
			name:     "custom range [0.4, 0.6] - orthogonal vectors",
			lo:       0.4,
			hi:       0.6,
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{0.0, 1.0, 0.0, 0.0},
			expected: true, // cosine distance = 0.5
		},
		{
			name:     "custom range [0.9, 1.0] - opposite vectors",
			lo:       0.9,
			hi:       1.0,
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{-1.0, -1.0, -1.0, -1.0},
			expected: true, // cosine distance = 1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rangeSim := scanner.RangeSimilarity(tt.lo, tt.hi)
			result := rangeSim(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("RangeSimilarity(%v, %v)(%v, %v) = %v, want %v", tt.lo, tt.hi, tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		assert   func(float32) bool
		a, b     []float32
		expected bool
	}{
		{
			name:     "assert distance < 0.1",
			assert:   func(d float32) bool { return d < 0.1 },
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{1.0, 1.0, 1.0, 1.0},
			expected: true, // cosine distance = 0.0
		},
		{
			name:     "assert distance > 0.4",
			assert:   func(d float32) bool { return d > 0.4 },
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{0.0, 1.0, 0.0, 0.0},
			expected: true, // cosine distance = 0.5
		},
		{
			name:     "assert distance == 0.5",
			assert:   func(d float32) bool { return d == 0.5 },
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{0.0, 1.0, 0.0, 0.0},
			expected: true, // cosine distance = 0.5
		},
		{
			name:     "assert distance > 0.9",
			assert:   func(d float32) bool { return d > 0.9 },
			a:        []float32{1.0, 1.0, 1.0, 1.0},
			b:        []float32{-1.0, -1.0, -1.0, -1.0},
			expected: true, // cosine distance = 1.0
		},
		{
			name:     "assert distance < 0.1 - false case",
			assert:   func(d float32) bool { return d < 0.1 },
			a:        []float32{1.0, 0.0, 0.0, 0.0},
			b:        []float32{0.0, 1.0, 0.0, 0.0},
			expected: false, // cosine distance = 0.5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cosineSim := scanner.CosineSimilarity(tt.assert)
			result := cosineSim(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CosineSimilarity(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCosineSimilarityPanicCases(t *testing.T) {
	t.Run("vectors with different lengths", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for vectors with different lengths")
			}
		}()
		scanner.HighSimilarity([]float32{1.0, 1.0, 1.0, 1.0}, []float32{1.0, 1.0})
	})

	t.Run("vectors with length not multiple of 4", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for vectors with length not multiple of 4")
			}
		}()
		scanner.HighSimilarity([]float32{1.0, 1.0, 1.0}, []float32{1.0, 1.0, 1.0})
	})
}
