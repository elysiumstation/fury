package ethcall

import "github.com/elysiumstation/fury/core/config/encoding"

type Config struct {
	Level encoding.LogLevel `long:"log-level"`
}
