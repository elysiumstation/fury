package helpers

import (
	vgcrypto "github.com/elysiumstation/fury/libs/crypto"
)

// GenerateID generates a 256 bit pseudo-random hash ID.
func GenerateID() string {
	return vgcrypto.RandomHash()
}
