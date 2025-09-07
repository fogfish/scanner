# Smart Text Processing

> "Finally, semantic chunking that understands meaning!" — Text processing, done intelligently.

Text scanners for Go that go beyond simple line-by-line processing. Built with semantic understanding at its core, making text chunking intelligent and context-aware.

[![Version](https://img.shields.io/github/v/tag/fogfish/scanner)](https://github.com/fogfish/scanner/releases)
[![Documentation](https://pkg.go.dev/badge/github.com/fogfish/scanner)](https://pkg.go.dev/github.com/fogfish/scanner)
[![Build Status](https://github.com/fogfish/scanner/workflows/build/badge.svg)](https://github.com/fogfish/scanner/actions/)
[![Git Hub](http://img.shields.io/github/stars/fogfish/scanner.svg?style=social&label=stars)](http://github.com/fogfish/scanner)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/scanner/badge.svg?branch=main)](https://coveralls.io/github/fogfish/scanner?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/fogfish/scanner)](https://goreportcard.com/report/github.com/fogfish/scanner)

Traditional text processing treats all chunks equally, but meaning isn't uniform across text. The `scanner` library uses embedding-based semantic analysis to group related content together, making it perfect for RAG systems, document analysis, and intelligent text processing pipelines.

## Quick Start: Semantic Chunking

Here's how to chunk text by semantic similarity instead of arbitrary boundaries:

```go
package main

import (
  "context"
  "strings"

  "github.com/fogfish/scanner"
)

func main() {
  api  := // create instance of embedding vector provider scanner.Embedder
  text := `
  The quick brown fox jumps over the lazy dog. 
  This is a classic pangram used in typography.
  
  Machine learning has revolutionized AI.
  Neural networks can now understand language context.
  
  Climate change affects global weather patterns.
  Rising temperatures impact ecosystems worldwide.
  `

  // Break text into sentences first
  sentences := scanner.NewSentencer(
    scanner.EndOfSentence, 
    strings.NewReader(text),
  )

  // Group sentences by semantic similarity
  semantic := scanner.NewSemantic(api, sentences)
  semantic.Window(10)                         // Look at 10 sentences at a time
  semantic.Similarity(scanner.HighSimilarity) // Group highly similar content

	// Get semantically coherent chunks
	for semantic.Scan() {
		chunk := semantic.Text()
		fmt.Printf("Semantic chunk: %v\n", chunk)
		// Output will group related sentences together:
		// - Typography sentences together
		// - AI/ML sentences together  
		// - Climate sentences together
	}
}
```

This approach produces chunks where sentences actually relate to each other, rather than arbitrary splits that might separate related concepts.

## Why Semantic Chunking Matters

**Traditional chunking problems:**
- Splits related content across chunks
- Breaks context mid-conversation
- Fixed boundaries ignore meaning
- Poor retrieval in RAG systems

**Semantic chunking benefits:**
- Keeps related content together
- Maintains semantic coherence
- Context-aware boundaries
- Better embedding similarity for retrieval

Perfect for:
- **RAG Systems**: Better retrieval through coherent chunks
- **Document Analysis**: Group related paragraphs and concepts
- **Content Summarization**: Preserve topic boundaries
- **Text Classification**: Maintain semantic integrity

## The Scanner Toolkit

Beyond semantic chunking, the library provides a complete text processing toolkit:

| Scanner       | Purpose                      | Use Case                       |
| ------------- | ---------------------------- | ------------------------------ |
| **Semantic**  | Groups by meaning similarity | RAG, document analysis         |
| **Sentencer** | Splits by punctuation        | Natural sentence boundaries    |
| **Slicer**    | Fixed delimiter splitting    | CSV, structured data           |
| **Chunker**   | Fixed-size chunks            | Token limits, simple splitting |
| **Sorter**    | Semantic sorting of data     | Organizing similar items       |
| **Identity**  | Entire input as one chunk    | Small documents                |

All scanners implement the familiar `bufio.Scanner` interface:

```go
for scanner.Scan() {
  text := scanner.Text()
  // Process chunk
}
```

## Similarity Control

Fine-tune semantic grouping with built-in similarity functions:

```go
semantic.Similarity(scanner.HighSimilarity)   // Very similar content (0.0-0.2)
semantic.Similarity(scanner.MediumSimilarity) // Related content (0.2-0.5)
semantic.Similarity(scanner.WeakSimilarity)   // Loosely related (0.5-0.8)

// Custom similarity threshold
semantic.Similarity(scanner.RangeSimilarity(0.1, 0.3))

// Custom similarity logic
semantic.Similarity(scanner.CosineSimilarity(func(d float32) bool {
    return d < 0.25 // Custom threshold
}))
```

## Algorithm Behavior

Control how chunks grow:

```go
// Compare new sentences to the first sentence in chunk (stable reference)
semantic.SimilarityWith(scanner.SIMILARITY_WITH_HEAD)

// Compare new sentences to the last added sentence (evolving reference)  
semantic.SimilarityWith(scanner.SIMILARITY_WITH_TAIL)
```

## Getting Started

The library requires Go 1.24 or later.

```bash
go get -u github.com/fogfish/scanner
```

Compatible with any embedding provider - OpenAI, Cohere, local models, or custom implementations. Just implement the simple `Embedder` interface:

```go
type Embedder interface {
    Embedding(ctx context.Context, text string) ([]float32, int, error)
}
```

## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)  
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

```bash
git clone https://github.com/fogfish/scanner
cd scanner
go test ./...
```

## License

[See LICENSE](LICENSE)
