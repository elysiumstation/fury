package commands

import (
	"encoding/hex"
	"errors"
)

const furyPubkeyLen = 64

var (
	ErrShouldBeAValidFuryPubkey = errors.New("should be a valid fury public key")
	ErrShouldBeAValidFuryID     = errors.New("should be a valid fury ID")
)

// IsFuryPubkey check if a string is a valid fury public key.
// A fury public key is a string of 64 characters containing only hexadecimal characters.
func IsFuryPubkey(pk string) bool {
	pkLen := len(pk)
	_, err := hex.DecodeString(pk)
	return pkLen == furyPubkeyLen && err == nil
}
