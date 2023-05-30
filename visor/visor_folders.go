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

package visor

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/visor/config"
	"github.com/elysiumstation/fury/visor/github"
	"github.com/elysiumstation/fury/visor/utils"
)

var furyDataNodeStartCmdArgs = []string{"datanode", "start"}

func (v *Visor) setCurrentFolder(sourceFolder, currentFolder string) error {
	v.log.Info("Setting current folder",
		logging.String("sourceFolder", sourceFolder),
		logging.String("currentFolder", currentFolder),
	)

	runConfPath := path.Join(sourceFolder, config.RunConfigFileName)
	runConfExists, err := utils.PathExists(runConfPath)
	if err != nil {
		return err
	}

	if !runConfExists {
		return fmt.Errorf("missing run config in %q folder", runConfPath)
	}

	if err := os.RemoveAll(currentFolder); err != nil {
		return fmt.Errorf("failed to remove current folder: %w", err)
	}

	if err := os.Symlink(sourceFolder, currentFolder); err != nil {
		return fmt.Errorf("failed to set current folder as current: %w", err)
	}

	return nil
}

func (v *Visor) installUpgradeFolder(ctx context.Context, folder, releaseTag string, conf config.AutoInstallConfig) error {
	v.log.Info("Automatically installing upgrade folder")

	runConf, err := config.ParseRunConfig(v.conf.CurrentRunConfigPath())
	if err != nil {
		return err
	}

	if conf.Asset.Name == "" {
		return missingAutoInstallAssetError("fury")
	}

	if err := os.MkdirAll(folder, 0o755); err != nil {
		return fmt.Errorf("failed to create upgrade folder %q, %w", folder, err)
	}

	assetsFetcher := github.NewAssetsFetcher(
		conf.GithubRepositoryOwner,
		conf.GithubRepository,
		[]string{conf.Asset.Name},
	)

	v.log.Info("Downloading asset from Github", logging.String("asset", conf.Asset.Name))
	if err := assetsFetcher.Download(ctx, releaseTag, folder); err != nil {
		return fmt.Errorf("failed to download release assets for tag %q: %w", releaseTag, err)
	}

	runConf.Name = releaseTag
	runConf.Fury.Binary.Path = conf.Asset.GetBinaryPath()

	if runConf.DataNode != nil {
		runConf.DataNode.Binary.Path = conf.Asset.GetBinaryPath()

		if len(runConf.DataNode.Binary.Args) != 0 && runConf.DataNode.Binary.Args[0] != furyDataNodeStartCmdArgs[0] {
			runConf.DataNode.Binary.Args = append(furyDataNodeStartCmdArgs, runConf.DataNode.Binary.Args[1:]...)
		}
	}

	runConfPath := path.Join(folder, config.RunConfigFileName)
	if err := runConf.WriteToFile(runConfPath); err != nil {
		return fmt.Errorf("failed to create run config %q: %w", runConfPath, err)
	}

	return nil
}

func (v *Visor) prepareNextUpgradeFolder(ctx context.Context, releaseTag string) error {
	v.log.Debug("preparing next upgrade folder",
		logging.String("furyTagVersion", releaseTag),
	)

	upgradeFolder := v.conf.UpgradeFolder(releaseTag)
	upgradeFolderExists, err := utils.PathExists(upgradeFolder)
	if err != nil {
		return err
	}

	if !upgradeFolderExists {
		autoInstallConf := v.conf.AutoInstall()
		if !autoInstallConf.Enabled {
			return fmt.Errorf("required upgrade folder %q is missing", upgradeFolder)
		}

		if err := v.installUpgradeFolder(ctx, upgradeFolder, releaseTag, autoInstallConf); err != nil {
			return fmt.Errorf("failed to auto install folder %q for release %q: %w", upgradeFolder, releaseTag, err)
		}
	}

	if err := v.setCurrentFolder(upgradeFolder, v.conf.CurrentFolder()); err != nil {
		return fmt.Errorf("failed to set current folder to %q: %w", v.conf.CurrentFolder(), err)
	}

	return nil
}

func missingAutoInstallAssetError(asset string) error {
	return fmt.Errorf("missing required auto install %s asset definition in Visor config", asset)
}
