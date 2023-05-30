package cmd

import (
	"fmt"
	"io"

	"github.com/elysiumstation/fury/cmd/furywallet/commands/cli"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/flags"
	"github.com/elysiumstation/fury/cmd/furywallet/commands/printer"
	"github.com/elysiumstation/fury/paths"
	svcStoreV1 "github.com/elysiumstation/fury/wallet/service/store/v1"

	"github.com/spf13/cobra"
)

var (
	locateServiceConfigLong = cli.LongDesc(`
		Locate the wallet service configuration file.
	`)

	locateServiceConfigExample = cli.Examples(`
		# Locate the wallet service configuration file
		{{.Software}} service config locate
	`)
)

type LocateServiceConfigsResponse struct {
	Path string `json:"path"`
}

type LocateServiceConfigsHandler func() (*LocateServiceConfigsResponse, error)

func NewCmdLocateServiceConfig(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func() (*LocateServiceConfigsResponse, error) {
		furyPaths := paths.New(rf.Home)

		svcConfig, err := svcStoreV1.InitialiseStore(furyPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise service store: %w", err)
		}

		return &LocateServiceConfigsResponse{
			Path: svcConfig.GetServiceConfigsPath(),
		}, nil
	}

	return BuildCmdLocateServiceConfigs(w, h, rf)
}

func BuildCmdLocateServiceConfigs(w io.Writer, handler LocateServiceConfigsHandler, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "locate",
		Short:   " Locate the wallet service configuration file",
		Long:    locateServiceConfigLong,
		Example: locateServiceConfigExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			resp, err := handler()
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintLocateServiceConfigsResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	return cmd
}

func PrintLocateServiceConfigsResponse(w io.Writer, resp *LocateServiceConfigsResponse) {
	p := printer.NewInteractivePrinter(w)

	str := p.String()
	defer p.Print(str)

	str.Text("The service configuration file is located at: ").SuccessText(resp.Path).NextLine()
}
