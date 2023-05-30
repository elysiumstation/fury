package commands_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/elysiumstation/fury/commands"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
	"github.com/stretchr/testify/assert"
)

func TestCheckAnnounceNode(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilAnnounceNodeFails)
	t.Run("Submitting a node registration without fury pub key fails", testAnnounceNodeWithoutFuryPubKeyFails)
	t.Run("Submitting a node registration with valid fury pub key succeeds", testAnnounceNodeWithValidFuryPubKeySucceeds)
	t.Run("Submitting a node registration with invalid encoding of fury pub key succeeds", testAnnounceNodeWithInvalidEncodingOfFuryPubKeyFails)
	t.Run("Submitting a node registration without ethereum pub key fails", testAnnounceNodeWithoutEthereumAddressFails)
	t.Run("Submitting a node registration with ethereum address succeeds", testAnnounceNodeWithEthereumAddressSucceeds)
	t.Run("Submitting a node registration without chain address fails", testAnnounceNodeWithoutChainPubKeyFails)
	t.Run("Submitting a node registration with chain pub key succeeds", testAnnounceNodeWithChainPubKeySucceeds)
	t.Run("Submitting a node registration with empty signatures fails", testAnnounceNodeWithEmptySignaturesFails)
	t.Run("Submitting a node registration with nonhex signatures fails", testAnnounceNodeWithNonhexSignaturesFails)
}

func testNilAnnounceNodeFails(t *testing.T) {
	err := checkAnnounceNode(nil)

	assert.Error(t, err)
}

func testAnnounceNodeWithoutFuryPubKeyFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{})
	assert.Contains(t, err.Get("announce_node.fury_pub_key"), commands.ErrIsRequired)
}

func testAnnounceNodeWithValidFuryPubKeySucceeds(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		FuryPubKey: hex.EncodeToString([]byte("0xDEADBEEF")),
	})
	assert.NotContains(t, err.Get("announce_node.fury_pub_key"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("announce_node.fury_pub_key"), commands.ErrShouldBeHexEncoded)
}

func testAnnounceNodeWithInvalidEncodingOfFuryPubKeyFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		FuryPubKey: "invalid-hex-encoding",
	})
	assert.Contains(t, err.Get("announce_node.fury_pub_key"), commands.ErrShouldBeAValidFuryPubkey)
}

func testAnnounceNodeWithoutEthereumAddressFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{})
	assert.Contains(t, err.Get("announce_node.ethereum_address"), commands.ErrIsRequired)
}

func testAnnounceNodeWithEthereumAddressSucceeds(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		EthereumAddress: "0xDEADBEEF",
	})
	assert.NotContains(t, err.Get("announce_node.ethereum_address"), commands.ErrIsRequired)
}

func testAnnounceNodeWithoutChainPubKeyFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{})
	assert.Contains(t, err.Get("announce_node.chain_pub_key"), commands.ErrIsRequired)
}

func testAnnounceNodeWithChainPubKeySucceeds(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		ChainPubKey: "0xDEADBEEF",
	})
	assert.NotContains(t, err.Get("announce_node.chain_pub_key"), commands.ErrIsRequired)
}

func testAnnounceNodeWithEmptySignaturesFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{})
	assert.Contains(t, err.Get("announce_node.ethereum_signature"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("announce_node.fury_signature"), commands.ErrIsRequired)
}

func testAnnounceNodeWithNonhexSignaturesFails(t *testing.T) {
	err := checkAnnounceNode(&commandspb.AnnounceNode{
		FurySignature: &commandspb.Signature{
			Value: "hello",
		},
		EthereumSignature: &commandspb.Signature{
			Value: "helloagain",
		},
	})
	fmt.Println(err)
	assert.Contains(t, err.Get("announce_node.ethereum_signature.value"), commands.ErrShouldBeHexEncoded)
	assert.Contains(t, err.Get("announce_node.fury_signature.value"), commands.ErrShouldBeHexEncoded)
}

func checkAnnounceNode(cmd *commandspb.AnnounceNode) commands.Errors {
	err := commands.CheckAnnounceNode(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
