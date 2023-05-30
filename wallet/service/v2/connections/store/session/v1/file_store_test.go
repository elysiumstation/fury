package v1_test

import (
	"context"
	"testing"

	vgrand "github.com/elysiumstation/fury/libs/rand"
	"github.com/elysiumstation/fury/paths"
	"github.com/elysiumstation/fury/wallet/service/v2/connections"
	v1 "github.com/elysiumstation/fury/wallet/service/v2/connections/store/session/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	t.Run("List tokens succeeds", testFileStoreListTokensSucceeds)
	t.Run("Deleting an existing token succeeds", testFileStoreDeletingExistingSessionSucceeds)
}

func testFileStoreListTokensSucceeds(t *testing.T) {
	ctx := context.Background()
	furyPaths := testHome(t)
	store := newTestFileStore(t, furyPaths)

	hostnameA := "a" + vgrand.RandomStr(5)
	hostnameB := "b" + vgrand.RandomStr(5)
	walletA := "a" + vgrand.RandomStr(5)
	walletB := "b" + vgrand.RandomStr(5)

	// given
	session1 := connections.Session{
		Hostname: hostnameA,
		Token:    connections.GenerateToken(),
		Wallet:   walletA,
	}

	// when
	err := store.TrackSession(session1)

	// then
	require.NoError(t, err)

	// given
	session2 := connections.Session{
		Hostname: hostnameB,
		Token:    connections.GenerateToken(),
		Wallet:   walletA,
	}

	// when
	err = store.TrackSession(session2)

	// then
	require.NoError(t, err)

	// given

	session3 := connections.Session{
		Hostname: hostnameA,
		Token:    connections.GenerateToken(),
		Wallet:   walletB,
	}

	// when
	err = store.TrackSession(session3)

	// then
	require.NoError(t, err)

	// when
	sessions, err := store.ListSessions(ctx)

	// then
	require.NoError(t, err)
	assert.Equal(t, []connections.Session{session1, session3, session2}, sessions)
}

func testFileStoreDeletingExistingSessionSucceeds(t *testing.T) {
	ctx := context.Background()
	furyPaths := testHome(t)
	store := newTestFileStore(t, furyPaths)

	// given
	session1 := connections.Session{
		Hostname: vgrand.RandomStr(5),
		Token:    connections.GenerateToken(),
		Wallet:   vgrand.RandomStr(5),
	}

	// when
	err := store.TrackSession(session1)

	// then
	require.NoError(t, err)

	// given
	session2 := connections.Session{
		Hostname: vgrand.RandomStr(5),
		Token:    connections.GenerateToken(),
		Wallet:   vgrand.RandomStr(5),
	}

	// when
	err = store.TrackSession(session2)

	// then
	require.NoError(t, err)

	// when
	err = store.DeleteSession(ctx, session1.Token)

	// then
	require.NoError(t, err)

	// when
	sessions, err := store.ListSessions(ctx)

	// then
	require.NoError(t, err)
	assert.Equal(t, []connections.Session{session2}, sessions)
}

type testFileStore struct {
	*v1.FileStore
}

func testHome(t *testing.T) paths.Paths {
	t.Helper()
	return paths.New(t.TempDir())
}

func newTestFileStore(t *testing.T, furyPaths paths.Paths) *testFileStore {
	t.Helper()

	tokenStore, err := v1.InitialiseStore(furyPaths)
	if err != nil {
		t.Fatalf("could not initialise the file store for tests: %v", err)
	}

	return &testFileStore{
		FileStore: tokenStore,
	}
}
