package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/elysiumstation/fury/cmd/furywallet/commands/cli"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/flags"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/printer"
	coreversion "github.com/elysiumstation/fury/version"
	wversion "github.com/elysiumstation/fury/wallet/version"
	"github.com/spf13/cobra"
)

var (
	softwareVersionLong = cli.LongDesc(`
		Get the version of the software and checks if its compatibility with the
		registered networks.

		This is NOT related to the wallet version. To get information about the wallet,
		use the "info" command.
	`)

	softwareVersionExample = cli.Examples(`
		# Get the version of the software
		{{.Software}} software version
	`)
)

type GetSoftwareVersionHandler func() *wversion.GetSoftwareVersionResponse

func NewCmdSoftwareVersion(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdGetSoftwareVersion(w, wversion.GetSoftwareVersionInfo, rf)
}

func BuildCmdGetSoftwareVersion(w io.Writer, handler GetSoftwareVersionHandler, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Get the version of the software",
		Long:    softwareVersionLong,
		Example: softwareVersionExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			resp := handler()

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintGetSoftwareVersionResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	return cmd
}

func PrintGetSoftwareVersionResponse(w io.Writer, resp *wversion.GetSoftwareVersionResponse) {
	p := printer.NewInteractivePrinter(w)

	str := p.String()
	defer p.Print(str)

	if wversion.IsUnreleased() {
		str.CrossMark().DangerText("You are running an unreleased version of the software (").DangerText(coreversion.Get()).DangerText(").").NextLine()
		str.Pad().DangerText("Use it at your own risk!").NextSection()
	}

	str.Text("Software version:").NextLine().WarningText(resp.Version).NextSection()
	str.Text("Git hash:").NextLine().WarningText(resp.GitHash).NextSection()

	str.RedArrow().DangerText("Important").NextLine()
	str.Text("The software version is NOT related to the key derivation version of your wallets.").NextLine()
	str.Bold("The software managing the wallets should not be confused with the wallets themselves.").NextLine()
	str.Text("To get the key derivation version of a wallet, see the following command:").NextSection()
	str.Code(fmt.Sprintf("%s describe --help", os.Args[0])).NextLine()

	str.BlueArrow().InfoText("Check the network compatibility").NextLine()
	str.Text("To determine if this software is compatible with the registered networks, use the following command:").NextSection()
	str.Code(fmt.Sprintf("%s software compatibility", os.Args[0])).NextLine()
}
