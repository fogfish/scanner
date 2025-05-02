package scanner

import "io"

func NewIdentity(r io.Reader) *Identity {
	return &Identity{Reader: r}
}

type Identity struct {
	io.Reader
	err error
	txt []byte
}

func (r *Identity) Err() error   { return r.err }
func (r *Identity) Text() string { return string(r.txt) }
func (r *Identity) Scan() bool {
	if r.err != nil || r.txt != nil {
		return false
	}

	r.txt, r.err = io.ReadAll(r.Reader)
	return r.err == nil
}
