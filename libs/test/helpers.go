package test

import (
	"path/filepath"

	vgrand "github.com/elysiumstation/fury/libs/rand"
)

func RandomPath() string {
	return filepath.Join("/tmp", "fury_tests", vgrand.RandomStr(10))
}
