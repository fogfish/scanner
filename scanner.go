package scanner

// Scanner is an interface similar to [bufio.Scanner].
// It defines core functionality defined by this library.
type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}
