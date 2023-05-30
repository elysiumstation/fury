package cmd

import (
	"fmt"
	"io"

	"github.com/elysiumstation/fury/cmd/furywallet/commands/cli"
	"github.com/spf13/cobra"
)

const disclaimerTxt = `The Fury Wallet is an application that allows users to, among other things,  (i)  access
other Fury applications;  (ii)  manage multiple wallets and keys; and  (iii) sign
transactions on the Fury network.  It is free, public and open source software.

The Fury Wallet is purely non-custodial application, meaning users never lose custody,
possession, or control of their digital assets at any time. Users are solely responsible for
the custody of the cryptographic private keys to their Fury Wallet and and should never
share their wallet credentials or seed phrase with anyone.

The Fury Wallet relies on emerging technologies that are subject to increased risk
through users misuse of things such as public/private key cryptography or failing to
properly update or run software to accommodate upgrades.  The developers of the
Fury Wallet do not operate or run the Fury Blockchain or any other blockchain.  Digital
tokens present market volatility risk, technical software risk, regulatory risk and
cybersecurity risk.  Software upgrades may contain bugs or security vulnerabilities that
might result in loss of functionality or assets.

The Fury Wallet is provided “as is”.  The developers of the Fury Wallet make no
representations or warranties of any kind, whether express or implied, statutory
or otherwise regarding the Fury Wallet.  They disclaim all warranties of
merchantability, quality, fitness for purpose.  They disclaim all warranties that the
Fury Wallet  is free of harmful components or errors.

No developer of the Fury Wallet accepts any responsibility for, or liability to users
in connection with their use of the Fury Wallet.  Users are solely responsible for
any associated wallet and no developer of the Fury Wallet is liable for any acts or
omissions by users in connection with or as a result of their Fury Wallet or other
associated wallet being compromised.
`

var disclaimerLong = cli.LongDesc(`
		Prints the disclaimer of the Fury Wallet.
	`)

type DisclaimerHandler func(home string, f *DisclaimerFlags) error

func NewCmdDisclaimer(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdDisclaimer(w, Disclaimer, rf)
}

func BuildCmdDisclaimer(w io.Writer, handler DisclaimerHandler, rf *RootFlags) *cobra.Command {
	f := &DisclaimerFlags{}

	cmd := &cobra.Command{
		Use:   "disclaimer",
		Short: "Prints the disclaimer of the Fury Wallet",
		Long:  disclaimerLong,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := handler(rf.Home, f); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

type DisclaimerFlags struct{}

func Disclaimer(home string, f *DisclaimerFlags) error {
	fmt.Print(disclaimerTxt)
	return nil
}
