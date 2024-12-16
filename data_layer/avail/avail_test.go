package avail_test

import (
	"fmt"
	"testing"

	"github.com/dymensionxyz/roller/data_layer/avail"
	"github.com/dymensionxyz/roller/utils/roller"
	"github.com/stretchr/testify/assert"
)

// Test NewAvail function
func TestNewAvail(t *testing.T) {
	rootDir := "/tmp/test_avail"
	a := avail.NewAvail(rootDir)

	assert.NotNil(t, a)
	assert.Equal(t, rootDir, a.GetRootDirectory(), "Root directory should be set correctly")
	assert.Equal(t, avail.DefaultRPCEndpoint, a.RpcEndpoint, "RPC endpoint should be set to default")
}

// Test GetPrivateKey
func TestGetPrivateKey(t *testing.T) {
	a := &avail.Avail{Mnemonic: "test-mnemonic"}
	privateKey, err := a.GetPrivateKey()

	assert.NoError(t, err)
	assert.Equal(t, "test-mnemonic", privateKey, "Mnemonic should match the expected value")
}

// Test CheckDABalance
func TestCheckDABalance(t *testing.T) {
	a := &avail.Avail{AccAddress: "test-address"}
	// mockBalance := availtypes.U128{Int: big.NewInt(10)}                        // Mocked balance
	// a.GetDAAccData = func() (availtypes.U128, error) { return mockBalance, nil } // Mock method
	keys, err := a.GetDAAccData(roller.RollappConfig{})
	fmt.Println("keys and errortr....", keys, err)

	notFundedAddresses, err := a.CheckDABalance()

	assert.NoError(t, err)
	assert.Len(t, notFundedAddresses, 1, "Should return one address needing funding")
	assert.Equal(t, "test-address", notFundedAddresses[0].Address)
}
