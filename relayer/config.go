package relayer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pterm/pterm"
	yaml "gopkg.in/yaml.v2"

	"github.com/dymensionxyz/roller/cmd/consts"
	"github.com/dymensionxyz/roller/utils"
	"github.com/dymensionxyz/roller/utils/roller"
)

func CreatePath(rlpCfg roller.RollappConfig) error {
	relayerHome := filepath.Join(rlpCfg.Home, consts.ConfigDirName.Relayer)
	pterm.Info.Printf("creating new ibc path from %s to %s\n", rlpCfg.HubData.ID, rlpCfg.RollappID)

	newPathCmd := exec.Command(
		consts.Executables.Relayer,
		"paths",
		"new",
		rlpCfg.HubData.ID,
		rlpCfg.RollappID,
		consts.DefaultRelayerPath,
		"--home",
		relayerHome,
	)
	if err := newPathCmd.Run(); err != nil {
		return err
	}

	return nil
}

func DeletePath(rlpCfg roller.RollappConfig) error {
	relayerHome := filepath.Join(rlpCfg.Home, consts.ConfigDirName.Relayer)
	pterm.Info.Printf("removing ibc path from %s to %s\n", rlpCfg.HubData.ID, rlpCfg.RollappID)

	newPathCmd := exec.Command(
		consts.Executables.Relayer,
		"paths",
		"delete",
		consts.DefaultRelayerPath,
		"--home",
		relayerHome,
	)
	if err := newPathCmd.Run(); err != nil {
		return err
	}

	return nil
}

type ChainConfig struct {
	ID            string
	RPC           string
	Denom         string
	AddressPrefix string
	GasPrices     string
}

func UpdateRlyConfigValue(
	rlpCfg roller.RollappConfig,
	keyPath []string,
	newValue interface{},
) error {
	rlyConfigPath := filepath.Join(
		rlpCfg.Home,
		consts.ConfigDirName.Relayer,
		"config",
		"config.yaml",
	)
	data, err := os.ReadFile(rlyConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load %s: %v", rlyConfigPath, err)
	}
	var rlyCfg map[interface{}]interface{}
	err = yaml.Unmarshal(data, &rlyCfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %v", err)
	}
	if err := utils.SetNestedValue(rlyCfg, keyPath, newValue); err != nil {
		return err
	}
	newData, err := yaml.Marshal(rlyCfg)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %v", err)
	}
	// nolint:gofumpt
	return os.WriteFile(rlyConfigPath, newData, 0o644)
}

func ReadRlyConfig(homeDir string) (map[interface{}]interface{}, error) {
	rlyConfigPath := filepath.Join(homeDir, consts.ConfigDirName.Relayer, "config", "config.yaml")
	data, err := os.ReadFile(rlyConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %v", rlyConfigPath, err)
	}
	var rlyCfg map[interface{}]interface{}
	err = yaml.Unmarshal(data, &rlyCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}
	return rlyCfg, nil
}

func WriteRlyConfig(homeDir string, rlyCfg map[interface{}]interface{}) error {
	rlyConfigPath := filepath.Join(homeDir, consts.ConfigDirName.Relayer, "config", "config.yaml")
	data, err := yaml.Marshal(rlyCfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	// nolint:gofumpt
	return os.WriteFile(rlyConfigPath, data, 0o644)
}
