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
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:8]
}

func (s *Shortener) ToBase62(num int) string {
	base62Chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if num == 0 {
		return "0"
	}

	base62Str := ""
	base := 62

	for num > 0 {
		remainder := num % base
		base62Str = string(base62Chars[remainder]) + base62Str
		num /= base
	}

	return base62Str
}
