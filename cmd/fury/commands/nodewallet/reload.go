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
	"context"
	"fmt"
	"time"

	vgjson "github.com/elysiumstation/fury/libs/json"
	"github.com/elysiumstation/fury/paths"

	"github.com/elysiumstation/fury/core/admin"
	"github.com/elysiumstation/fury/core/config"
	"github.com/elysiumstation/fury/logging"

	"github.com/jessevdk/go-flags"
)

type reloadCmd struct {
	config.OutputFlag

	Config admin.Config

	Chain string `short:"c" long:"chain" required:"true" description:"The chain to be imported" choice:"fury" choice:"ethereum"`
}

func (opts *reloadCmd) Execute(_ []string) error {
	output, err := opts.GetOutput()
	if err != nil {
		return err
	}

	log := logging.NewLoggerFromConfig(logging.NewDefaultConfig())
	defer log.AtExit()

	furyPaths := paths.New(rootCmd.FuryHome)

	_, conf, err := config.EnsureNodeConfig(furyPaths)
	if err != nil {
		return err
	}

	opts.Config = conf.Admin

	if _, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse(); err != nil {
		return err
	}

	sc := admin.NewClient(log, opts.Config)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var resp *admin.NodeWalletReloadReply
	switch opts.Chain {
	case furyChain, ethereumChain:
		resp, err = sc.NodeWalletReload(ctx, opts.Chain)
		if err != nil {
			return fmt.Errorf("failed to reload node wallet: %w", err)
		}
	default:
		return fmt.Errorf("chain %q is not supported", opts.Chain)
	}
	if output.IsHuman() {
		fmt.Println(green("reload successful:"))

		vgjson.PrettyPrint(resp)
	} else if output.IsJSON() {
		if err := vgjson.Print(resp); err != nil {
			return err
		}
	}

	return nil
}
