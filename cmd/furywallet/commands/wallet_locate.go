package cmd

import (
	"fmt"
	"io"

	"github.com/elysiumstation/fury/cmd/furywallet/commands/cli"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/flags"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/printer"
	"github.com/elysiumstation/fury/paths"
	"github.com/elysiumstation/fury/wallet/wallets"

	"github.com/spf13/cobra"
)

var (
	locateWalletsLong = cli.LongDesc(`
		Locate the folder in which all the wallet files are stored.
	`)

	locateWalletsExample = cli.Examples(`
		# Locate wallet files
		{{.Software}} locate
	`)
)

type LocateWalletsResponse struct {
	Path string `json:"path"`
}

type LocateWalletsHandler func() (*LocateWalletsResponse, error)

func NewCmdLocateWallets(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func() (*LocateWalletsResponse, error) {
		furyPaths := paths.New(rf.Home)

		walletStore, err := wallets.InitialiseStoreFromPaths(furyPaths, false)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise networks store: %w", err)
		}

		return &LocateWalletsResponse{
			Path: walletStore.GetWalletsPath(),
		}, nil
	}

	return BuildCmdLocateWallets(w, h, rf)
}

func BuildCmdLocateWallets(w io.Writer, handler LocateWalletsHandler, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "locate",
		Short:   "Locate the folder containing the wallet files",
		Long:    locateWalletsLong,
		Example: locateWalletsExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			resp, err := handler()
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintLocateWalletsResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	return cmd
}

func PrintLocateWalletsResponse(w io.Writer, resp *LocateWalletsResponse) {
	p := printer.NewInteractivePrinter(w)

	str := p.String()
	defer p.Print(str)

	str.Text("Wallet files are located at: ").SuccessText(resp.Path).NextLine()
}
