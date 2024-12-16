package avail

import (
	"os"
	"testing"

	// "github.com/dymensionxyz/roller/data_layer/avail"

	"github.com/stretchr/testify/assert"
)

func TestWriteConfigToTOML(t *testing.T) {
	path := "/tmp/test_avail_config.toml"
	cfg := Avail{
		Root:     "/tmp/test_root",
		AppID:    1,
		Mnemonic: "test-mnemonic",
	}

	err := writeConfigToTOML(path, cfg)
	defer os.Remove(path)

	assert.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestLoadConfigFromTOML(t *testing.T) {
	path := "/tmp/test_avail_config.toml"
	cfg := Avail{
		Root:     "/tmp/test_root",
		AppID:    1,
		Mnemonic: "test-mnemonic",
	}

	// Write config for testing
	err := writeConfigToTOML(path, cfg)
	assert.NoError(t, err)
	defer os.Remove(path)

	// Load config
	loadedCfg, err := loadConfigFromTOML(path)
	assert.NoError(t, err)
	assert.Equal(t, cfg, loadedCfg, "Loaded configuration should match the written configuration")
}
