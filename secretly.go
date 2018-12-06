package secretly

import (
	"errors"
	"strings"
)

type Backend interface {
	Get(path string) (string, error)
	Put(path, val string) error
}

var Backends = map[string]Backend{}

var (
	ErrUnknownBackend      = errors.New("unknown backend")
	ErrNoSuchSecret        = errors.New("no such secret")
	ErrAmbigiousSecret     = errors.New("ambigious secret, provide a subpath after ;")
	ErrInvalidSecret       = errors.New("invalid secret data detected")
	ErrInvalidSecretAccess = errors.New("invalid secret access")
)

func AddBackend(name string, backend Backend) {
	Backends[name] = backend
}

func Lookup(id string) (Backend, string, error) {
	colon := strings.IndexByte(id, ':')
	if colon == -1 {
		return nil, "", ErrUnknownBackend
	}

	be, ok := Backends[id[:colon]]
	if !ok {
		return nil, "", ErrUnknownBackend
	}

	return be, id[colon+1:], nil
}

func Get(id string) (string, error) {
	be, path, err := Lookup(id)
	if err != nil {
		return "", err
	}

	return be.Get(path)
}

func Put(id, value string) error {
	be, path, err := Lookup(id)
	if err != nil {
		return err
	}

	return be.Put(path, value)
}
