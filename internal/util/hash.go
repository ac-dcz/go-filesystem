package util

import (
	"crypto/sha1"
	"hash"
)

type Hasher struct {
	sha1 hash.Hash
}

func NewHasher() *Hasher {
	return &Hasher{
		sha1: sha1.New(),
	}
}

func (h *Hasher) Add(data []byte) *Hasher {
	h.sha1.Write(data)
	return h
}

func (h *Hasher) Sum() []byte {
	defer h.sha1.Reset()
	return h.sha1.Sum(nil)
}
