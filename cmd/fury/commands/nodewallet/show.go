// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package nodewallet

import (
	vgjson "github.com/elysiumstation/fury/libs/json"
	"github.com/elysiumstation/fury/paths"

	"github.com/elysiumstation/fury/core/config"
	"github.com/elysiumstation/fury/core/nodewallets"
	"github.com/elysiumstation/fury/core/nodewallets/registry"
	"github.com/elysiumstation/fury/logging"

	"github.com/jessevdk/go-flags"
)

type showCmd struct {
	Config nodewallets.Config
}

func (opts *showCmd) Execute(_ []string) error {
	log := logging.NewLoggerFromConfig(logging.NewDefaultConfig())
	defer log.AtExit()

	registryPass, err := rootCmd.PassphraseFile.Get("node wallet", false)
	if err != nil {
		return err
	}

	furyPaths := paths.New(rootCmd.FuryHome)

	_, conf, err := config.EnsureNodeConfig(furyPaths)
	if err != nil {
		return err
	}

	opts.Config = conf.NodeWallet

	if _, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse(); err != nil {
		return err
	}

	registryLoader, err := registry.NewLoader(furyPaths, registryPass)
	if err != nil {
		return err
	}

	registry, err := registryLoader.Get(registryPass)
	if err != nil {
		return err
	}

	if err = vgjson.PrettyPrint(registry); err != nil {
		return err
	}
	return nil
}
