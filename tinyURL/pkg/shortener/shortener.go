package shortener

import (
	"crypto/rand"
	"encoding/base64"
)

type Shortener struct{}

func New() *Shortener {
	return &Shortener{}
}

func (s *Shortener) Generate() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
