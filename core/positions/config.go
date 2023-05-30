// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.FURY file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package positions

import (
	"github.com/elysiumstation/fury/libs/config/encoding"
	"github.com/elysiumstation/fury/logging"
)

// namedLogger is the identifier for package and should ideally match the package name
// this is simply emitted as a hierarchical label e.g. 'api.grpc'.
const namedLogger = "position"

// Config represents the configuration of the position engine.
type Config struct {
	Level                 encoding.LogLevel `long:"log-level"`
	StreamPositionVerbose encoding.Bool     `long:"log-position-update"`
}

// NewDefaultConfig creates an instance of the package specific configuration, given a
// pointer to a logger instance to be used for logging within the package.
func NewDefaultConfig() Config {
	return Config{
		Level:                 encoding.LogLevel{Level: logging.InfoLevel},
		StreamPositionVerbose: false,
	}
}