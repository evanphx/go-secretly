package file

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/pbkdf2"

	"github.com/evanphx/go-secretly"
)

type f struct{}

const pemType = "SECRETLY ENCODED FILE"

func (_ *f) Put(path, val string) error {
	parts := strings.SplitN(path, "|", 2)

	path = parts[0]
	pass := path

	if len(parts) == 2 {
		pass = parts[1]
	}

	nonce := make([]byte, chacha20poly1305.NonceSize)

	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	key := pbkdf2.Key([]byte(pass), nonce, 4096, chacha20poly1305.KeySize, sha256.New)

	c, err := chacha20poly1305.New(key)
	if err != nil {
		return err
	}

	data := c.Seal(nil, nonce, []byte(val), nil)

	o, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer o.Close()

	var b pem.Block
	b.Bytes = data
	b.Type = pemType
	b.Headers = map[string]string{
		"S-Nonce": base64.StdEncoding.EncodeToString(nonce),
	}

	return pem.Encode(o, &b)
}

func (_ *f) Get(path string) (string, error) {
	parts := strings.SplitN(path, "|", 2)

	path = parts[0]
	pass := path

	if len(parts) == 2 {
		pass = parts[1]
	}

	i, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", secretly.ErrNoSuchSecret
		}

		return "", err
	}

	data, err := ioutil.ReadAll(i)
	if err != nil {
		return "", err
	}

	blk, _ := pem.Decode(data)

	if blk == nil {
		return "", secretly.ErrInvalidSecret
	}

	if blk.Type != pemType {
		return "", secretly.ErrInvalidSecret
	}

	nonce, err := base64.StdEncoding.DecodeString(blk.Headers["S-Nonce"])
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(pass), nonce, 4096, chacha20poly1305.KeySize, sha256.New)

	c, err := chacha20poly1305.New(key)
	if err != nil {
		return "", err
	}

	raw, err := c.Open(nil, nonce, blk.Bytes, nil)
	if err != nil {
		return "", secretly.ErrInvalidSecretAccess
	}

	return string(raw), nil
}

func init() {
	secretly.AddBackend("file", &f{})
}
