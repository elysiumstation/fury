package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/elysiumstation/fury/cmd/furywallet/commands/cli"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/flags"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/printer"
	"github.com/elysiumstation/fury/wallet/api"
	"github.com/elysiumstation/fury/wallet/wallets"

	"github.com/spf13/cobra"
)

var (
	rotateKeyLong = cli.LongDesc(`
		Build a signed key rotation transaction as a Base64 encoded string.
		Choose a public key to rotate to and target block height.

		The generated transaction can be sent using the command: "tx send".
	`)

	rotateKeyExample = cli.Examples(`
		# Build signed transaction for rotating to new key public key
		{{.Software}} key rotate --wallet WALLET --tx-height TX_HEIGHT --chain-id CHAIN_ID --target-height TARGET_HEIGHT --pubkey PUBLIC_KEY --current-pubkey CURRENT_PUBLIC_KEY
	`)
)

type RotateKeyHandler func(api.AdminRotateKeyParams, string) (api.AdminRotateKeyResult, error)

func NewCmdRotateKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(params api.AdminRotateKeyParams, passphrase string) (api.AdminRotateKeyResult, error) {
		ctx := context.Background()

		walletStore, err := wallets.InitialiseStore(rf.Home, false)
		if err != nil {
			return api.AdminRotateKeyResult{}, fmt.Errorf("could not initialise wallets store: %w", err)
		}
		defer walletStore.Close()

		if _, errDetails := api.NewAdminUnlockWallet(walletStore).Handle(ctx, api.AdminUnlockWalletParams{
			Wallet:     params.Wallet,
			Passphrase: passphrase,
		}); errDetails != nil {
			return api.AdminRotateKeyResult{}, errors.New(errDetails.Data)
		}

		rawResult, errDetails := api.NewAdminRotateKey(walletStore).Handle(context.Background(), params)
		if errDetails != nil {
			return api.AdminRotateKeyResult{}, errors.New(errDetails.Data)
		}
		return rawResult.(api.AdminRotateKeyResult), nil
	}

	return BuildCmdRotateKey(w, h, rf)
}

func BuildCmdRotateKey(w io.Writer, handler RotateKeyHandler, rf *RootFlags) *cobra.Command {
	f := RotateKeyFlags{}

	cmd := &cobra.Command{
		Use:     "rotate",
		Short:   "Build a signed key rotation transaction",
		Long:    rotateKeyLong,
		Example: rotateKeyExample,
		RunE: func(_ *cobra.Command, args []string) error {
			req, pass, err := f.Validate()
			if err != nil {
				return err
			}

			resp, err := handler(req, pass)
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintRotateKeyResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Wallet holding the master key and new public key",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)
	cmd.Flags().StringVar(&f.ToPublicKey,
		"new-pubkey",
		"",
		"A public key to rotate to. Should be generated by wallet's 'generate' command.",
	)
	cmd.Flags().StringVar(&f.ChainID,
		"chain-id",
		"",
		"The identifier of the chain on which the rotation will be done.",
	)
	cmd.Flags().StringVar(&f.FromPublicKey,
		"current-pubkey",
		"",
		"A public key to rotate from. Should be currently used public key.",
	)
	cmd.Flags().Uint64Var(&f.SubmissionBlockHeight,
		"tx-height",
		0,
		"It should be close to the current block height when the transaction is applied, with a threshold of ~ - 150 blocks.",
	)
	cmd.Flags().Uint64Var(&f.EnactmentBlockHeight,
		"target-height",
		0,
		"Height of block where the public key change will take effect",
	)

	autoCompleteWallet(cmd, rf.Home, "wallet")

	return cmd
}

type RotateKeyFlags struct {
	Wallet                string
	PassphraseFile        string
	FromPublicKey         string
	ToPublicKey           string
	ChainID               string
	SubmissionBlockHeight uint64
	EnactmentBlockHeight  uint64
}

func (f *RotateKeyFlags) Validate() (api.AdminRotateKeyParams, string, error) {
	if f.ToPublicKey == "" {
		return api.AdminRotateKeyParams{}, "", flags.MustBeSpecifiedError("new-pubkey")
	}

	if f.FromPublicKey == "" {
		return api.AdminRotateKeyParams{}, "", flags.MustBeSpecifiedError("current-pubkey")
	}

	if f.ChainID == "" {
		return api.AdminRotateKeyParams{}, "", flags.MustBeSpecifiedError("chain-id")
	}

	if f.EnactmentBlockHeight == 0 {
		return api.AdminRotateKeyParams{}, "", flags.MustBeSpecifiedError("target-height")
	}

	if f.SubmissionBlockHeight == 0 {
		return api.AdminRotateKeyParams{}, "", flags.MustBeSpecifiedError("tx-height")
	}

	if f.EnactmentBlockHeight <= f.SubmissionBlockHeight {
		return api.AdminRotateKeyParams{}, "", flags.RequireLessThanFlagError("tx-height", "target-height")
	}

	if len(f.Wallet) == 0 {
		return api.AdminRotateKeyParams{}, "", flags.MustBeSpecifiedError("wallet")
	}

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return api.AdminRotateKeyParams{}, "", err
	}

	return api.AdminRotateKeyParams{
		Wallet:                f.Wallet,
		FromPublicKey:         f.FromPublicKey,
		ToPublicKey:           f.ToPublicKey,
		ChainID:               f.ChainID,
		SubmissionBlockHeight: f.SubmissionBlockHeight,
		EnactmentBlockHeight:  f.EnactmentBlockHeight,
	}, passphrase, nil
}

func PrintRotateKeyResponse(w io.Writer, req api.AdminRotateKeyResult) {
	p := printer.NewInteractivePrinter(w)

	str := p.String()
	defer p.Print(str)

	str.CheckMark().SuccessText("Key rotation succeeded").NextSection()
	str.Text("Transaction (base64-encoded):").NextLine()
	str.WarningText(req.EncodedTransaction).NextSection()
	str.Text("Master public key used:").NextLine()
	str.WarningText(req.MasterPublicKey).NextLine()
}
