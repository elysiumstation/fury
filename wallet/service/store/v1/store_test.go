package v1_test

import (
	"testing"
	"time"

	vgencoding "github.com/elysiumstation/fury/libs/encoding"
	vgrand "github.com/elysiumstation/fury/libs/rand"
	vgtest "github.com/elysiumstation/fury/libs/test"
	"github.com/elysiumstation/fury/paths"
	"github.com/elysiumstation/fury/wallet/service"
	storeV1 "github.com/elysiumstation/fury/wallet/service/store/v1"
	v1 "github.com/elysiumstation/fury/wallet/service/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestFileStoreV1(t *testing.T) {
	t.Run("New store succeeds", testNewStoreSucceeds)
	t.Run("Saving already existing RSA keys succeeds", testFileStoreV1SaveAlreadyExistingRSAKeysSucceeds)
	t.Run("Saving RSA keys succeeds", testFileStoreV1SaveRSAKeysSucceeds)
	t.Run("Verifying non-existing RSA keys fails", testFileStoreV1VerifyingNonExistingRSAKeysFails)
	t.Run("Verifying existing RSA keys succeeds", testFileStoreV1VerifyingExistingRSAKeysSucceeds)
	t.Run("Getting non-existing RSA keys fails", testFileStoreV1GetNonExistingRSAKeysFails)
	t.Run("Getting existing RSA keys succeeds", testFileStoreV1GetExistingRSAKeysSucceeds)
	t.Run("Getting config while not being initialised succeeds", testFileStoreV1GetConfigWhileNotInitialisedSucceeds)
	t.Run("Saving config succeeds", testFileStoreV1SavingConfigSucceeds)
	t.Run("Verifying config succeeds", testFileStoreV1VerifyingConfigSucceeds)
}

func testNewStoreSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	s, err := storeV1.InitialiseStore(furyHome)

	require.NoError(t, err)
	assert.NotNil(t, s)
	vgtest.AssertDirAccess(t, rsaKeysHome(t, furyHome))
}

func testFileStoreV1SaveAlreadyExistingRSAKeysSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)
	keys := &v1.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)

	// when
	err = s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
}

func testFileStoreV1SaveRSAKeysSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)
	keys := &v1.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, publicRSAKeyFilePath(t, furyHome))
	vgtest.AssertFileAccess(t, privateRSAKeyFilePath(t, furyHome))

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, keys, returnedKeys)
}

func testFileStoreV1VerifyingNonExistingRSAKeysFails(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)

	// when
	exists, err := s.RSAKeysExists()

	// then
	assert.NoError(t, err)
	assert.False(t, exists)
}

func testFileStoreV1VerifyingExistingRSAKeysSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)
	keys := &v1.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, publicRSAKeyFilePath(t, furyHome))
	vgtest.AssertFileAccess(t, privateRSAKeyFilePath(t, furyHome))

	// when
	exists, err := s.RSAKeysExists()

	// then
	require.NoError(t, err)
	assert.True(t, exists)
}

func testFileStoreV1GetNonExistingRSAKeysFails(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)

	// when
	keys, err := s.GetRsaKeys()

	// then
	assert.Error(t, err)
	assert.Nil(t, keys)
}

func testFileStoreV1GetExistingRSAKeysSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)
	keys := &v1.RSAKeys{
		Pub:  []byte("my public key"),
		Priv: []byte("my private key"),
	}

	// when
	err := s.SaveRSAKeys(keys)

	// then
	require.NoError(t, err)
	vgtest.AssertFileAccess(t, publicRSAKeyFilePath(t, furyHome))
	vgtest.AssertFileAccess(t, privateRSAKeyFilePath(t, furyHome))

	// when
	returnedKeys, err := s.GetRsaKeys()

	// then
	require.NoError(t, err)
	assert.Equal(t, keys, returnedKeys)
}

func testFileStoreV1GetConfigWhileNotInitialisedSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)

	// when
	cfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, service.DefaultConfig(), cfg)
}

func testFileStoreV1SavingConfigSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)
	originalCfg := &service.Config{
		LogLevel: vgencoding.LogLevel{
			Level: zap.DebugLevel,
		},
		Server: service.ServerConfig{
			Port: 123456789,
			Host: vgrand.RandomStr(5),
		},
		APIV1: service.APIV1Config{
			MaximumTokenDuration: vgencoding.Duration{
				Duration: 234 * time.Second,
			},
		},
	}

	// when
	err := s.SaveConfig(originalCfg)

	// then
	require.NoError(t, err)

	// when
	cfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, originalCfg, cfg)
}

func testFileStoreV1VerifyingConfigSucceeds(t *testing.T) {
	furyHome := newFuryHome(t)

	// given
	s := initialiseFromPath(t, furyHome)
	originalCfg := &service.Config{
		LogLevel: vgencoding.LogLevel{
			Level: zap.DebugLevel,
		},
		Server: service.ServerConfig{
			Port: 123456789,
			Host: vgrand.RandomStr(5),
		},
		APIV1: service.APIV1Config{
			MaximumTokenDuration: vgencoding.Duration{
				Duration: 234 * time.Second,
			},
		},
	}

	// when
	err := s.SaveConfig(originalCfg)

	// then
	require.NoError(t, err)

	// when
	cfg, err := s.GetConfig()

	// then
	require.NoError(t, err)
	assert.Equal(t, originalCfg, cfg)
}

func initialiseFromPath(t *testing.T, furyHome *paths.CustomPaths) *storeV1.Store {
	t.Helper()
	s, err := storeV1.InitialiseStore(furyHome)
	if err != nil {
		t.Fatalf("couldn't initialise store: %v", err)
	}

	return s
}

func newFuryHome(t *testing.T) *paths.CustomPaths {
	t.Helper()
	return &paths.CustomPaths{CustomHome: t.TempDir()}
}

func rsaKeysHome(t *testing.T, furyHome *paths.CustomPaths) string {
	t.Helper()
	return furyHome.DataPathFor(paths.WalletServiceRSAKeysDataHome)
}

func publicRSAKeyFilePath(t *testing.T, furyHome *paths.CustomPaths) string {
	t.Helper()
	return furyHome.DataPathFor(paths.WalletServicePublicRSAKeyDataFile)
}

func privateRSAKeyFilePath(t *testing.T, furyHome *paths.CustomPaths) string {
	t.Helper()
	return furyHome.DataPathFor(paths.WalletServicePrivateRSAKeyDataFile)
}
