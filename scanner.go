//
// Copyright (C) 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/scanner
//

package scanner

// Scanner is an interface similar to [bufio.Scanner].
// It defines core functionality defined by this library.
type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}
