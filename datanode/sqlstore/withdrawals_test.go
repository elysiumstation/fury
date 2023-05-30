// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package sqlstore_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/elysiumstation/fury/datanode/entities"
	"github.com/elysiumstation/fury/datanode/sqlstore"
	"github.com/elysiumstation/fury/datanode/sqlstore/helpers"
	"github.com/elysiumstation/fury/protos/fury"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawals(t *testing.T) {
	t.Run("Upsert should insert withdrawals if one doesn't exist for the block", testAddWithdrawalForNewBlock)
	t.Run("Upsert should update withdrawals if one already exists for the block", testUpdateWithdrawalForBlockIfExists)
	t.Run("Upsert should insert withdrawal updates if the same withdrawal id is inserted in a different block", testInsertWithdrawalUpdatesIfNewBlock)
	t.Run("GetByID should retrieve the latest state of the withdrawal with the given ID", testWithdrawalsGetByID)
	t.Run("GetByParty should retrieve the latest state of all withdrawals for a given party", testWithdrawalsGetByParty)
	t.Run("GetByTxHash", testWithdrawalsGetByTxHash)
}

func TestWithdrawalsPagination(t *testing.T) {
	t.Run("should return all withdrawals if no pagination is specified", testWithdrawalsPaginationNoPagination)
	t.Run("should return the first page of results if first is provided", testWithdrawalsPaginationFirst)
	t.Run("should return the last page of results if last is provided", testWithdrawalsPaginationLast)
	t.Run("should return the specified page of results if first and after are provided", testWithdrawalsPaginationFirstAfter)
	t.Run("should return the specified page of results if last and before are provided", testWithdrawalsPaginationLastBefore)

	t.Run("should return all withdrawals if no pagination is specified - newest first", testWithdrawalsPaginationNoPaginationNewestFirst)
	t.Run("should return the first page of results if first is provided - newest first", testWithdrawalsPaginationFirstNewestFirst)
	t.Run("should return the last page of results if last is provided - newest first", testWithdrawalsPaginationLastNewestFirst)
	t.Run("should return the specified page of results if first and after are provided - newest first", testWithdrawalsPaginationFirstAfterNewestFirst)
	t.Run("should return the specified page of results if last and before are provided - newest first", testWithdrawalsPaginationLastBeforeNewestFirst)

	t.Run("should return all withdrawals between dates if no pagination is specified", testWithdrawalsPaginationBetweenDatesNoPagination)
	t.Run("should return the first page of results between dates if first is provided", testWithdrawalsPaginationBetweenDatesFirst)
	t.Run("should return the last page of results between dates if last is provided", testWithdrawalsPaginationBetweenDatesLast)
	t.Run("should return the specified page of results between dates if first and after are provided", testWithdrawalsPaginationBetweenDatesFirstAfter)
	t.Run("should return the specified page of results between dates if last and before are provided", testWithdrawalsPaginationBetweenDatesLastBefore)

	t.Run("should return all withdrawals between dates if no pagination is specified - newest first", testWithdrawalsPaginationBetweenDatesNoPaginationNewestFirst)
	t.Run("should return the first page of results between dates if first is provided - newest first", testWithdrawalsPaginationBetweenDatesFirstNewestFirst)
	t.Run("should return the last page of results between dates if last is provided - newest first", testWithdrawalsPaginationBetweenDatesLastNewestFirst)
	t.Run("should return the specified page of results between dates if first and after are provided - newest first", testWithdrawalsPaginationBetweenDatesFirstAfterNewestFirst)
	t.Run("should return the specified page of results between dates if last and before are provided - newest first", testWithdrawalsPaginationBetweenDatesLastBeforeNewestFirst)
}

func setupWithdrawalStoreTests(t *testing.T) (*sqlstore.Blocks, *sqlstore.Withdrawals, sqlstore.Connection) {
	t.Helper()
	bs := sqlstore.NewBlocks(connectionSource)
	ws := sqlstore.NewWithdrawals(connectionSource)
	return bs, ws, connectionSource.Connection
}

func testAddWithdrawalForNewBlock(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, ws, conn := setupWithdrawalStoreTests(t)

	var rowCount int

	err := conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	require.NoError(t, err)
	assert.Equal(t, 0, rowCount)

	block := addTestBlock(t, ctx, bs)
	withdrawalProto := getTestWithdrawal(testID, testID, testID, testAmount, testID, block.FuryTime)

	withdrawal, err := entities.WithdrawalFromProto(withdrawalProto, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")
	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)
	err = conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowCount)
}

func testUpdateWithdrawalForBlockIfExists(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, ws, conn := setupWithdrawalStoreTests(t)

	var rowCount int

	err := conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	require.NoError(t, err)
	assert.Equal(t, 0, rowCount)

	block := addTestBlock(t, ctx, bs)
	withdrawalProto := getTestWithdrawal(testID, testID, testID, testAmount, testID, block.FuryTime)

	withdrawal, err := entities.WithdrawalFromProto(withdrawalProto, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)
	err = conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowCount)

	withdrawal.Status = entities.WithdrawalStatus(fury.Withdrawal_STATUS_FINALIZED)

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)
	err = conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	assert.NoError(t, err)
	var status entities.WithdrawalStatus
	err = pgxscan.Get(ctx, conn, &status, `select status from withdrawals where id = $1 and fury_time = $2`, withdrawal.ID, withdrawal.FuryTime)
	assert.NoError(t, err)
	assert.Equal(t, entities.WithdrawalStatusFinalized, status)
}

func testInsertWithdrawalUpdatesIfNewBlock(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, ws, conn := setupWithdrawalStoreTests(t)

	var rowCount int

	err := conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	require.NoError(t, err)
	assert.Equal(t, 0, rowCount)

	source := &testBlockSource{bs, time.Now()}
	block := source.getNextBlock(t, ctx)
	withdrawalProto := getTestWithdrawal(testID, testID, testID, testAmount, testID, block.FuryTime)

	withdrawal, err := entities.WithdrawalFromProto(withdrawalProto, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)
	err = conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowCount)

	block = source.getNextBlock(t, ctx)
	withdrawalProto.Status = fury.Withdrawal_STATUS_FINALIZED
	withdrawal, err = entities.WithdrawalFromProto(withdrawalProto, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)
	err = conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	assert.NoError(t, err)
	var status entities.WithdrawalStatus
	err = pgxscan.Get(ctx, conn, &status, `select status from withdrawals where id = $1 and fury_time = $2`, withdrawal.ID, withdrawal.FuryTime)
	assert.NoError(t, err)
	assert.Equal(t, entities.WithdrawalStatusFinalized, status)
}

func testWithdrawalsGetByID(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, ws, conn := setupWithdrawalStoreTests(t)

	var rowCount int

	err := conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	require.NoError(t, err)
	assert.Equal(t, 0, rowCount)

	source := &testBlockSource{bs, time.Now()}
	block := source.getNextBlock(t, ctx)
	withdrawalProto := getTestWithdrawal(testID, testID, testID, testAmount, testID, block.FuryTime)

	withdrawal, err := entities.WithdrawalFromProto(withdrawalProto, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)
	err = conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowCount)

	block = source.getNextBlock(t, ctx)
	withdrawalProto.Status = fury.Withdrawal_STATUS_FINALIZED
	withdrawal, err = entities.WithdrawalFromProto(withdrawalProto, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)

	got, err := ws.GetByID(ctx, withdrawalProto.Id)
	assert.NoError(t, err)
	withdrawal.CreatedTimestamp = withdrawal.CreatedTimestamp.Truncate(time.Microsecond)
	withdrawal.WithdrawnTimestamp = withdrawal.WithdrawnTimestamp.Truncate(time.Microsecond)
	assert.Equal(t, *withdrawal, got)
}

func testWithdrawalsGetByParty(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, ws, conn := setupWithdrawalStoreTests(t)

	var rowCount int

	err := conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	require.NoError(t, err)
	assert.Equal(t, 0, rowCount)

	source := &testBlockSource{bs, time.Now()}
	block := source.getNextBlock(t, ctx)
	withdrawalProto1 := getTestWithdrawal(testID, testID, testID, testAmount, testID, block.FuryTime)
	withdrawalProto1.Id = "deadbeef01"

	withdrawalProto2 := getTestWithdrawal(testID, testID, testID, testAmount, testID, block.FuryTime)
	withdrawalProto2.Id = "deadbeef02"

	want := make([]entities.Withdrawal, 0)

	withdrawal, err := entities.WithdrawalFromProto(withdrawalProto1, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)

	block = source.getNextBlock(t, ctx)
	withdrawalProto1.Status = fury.Withdrawal_STATUS_FINALIZED
	withdrawal, err = entities.WithdrawalFromProto(withdrawalProto1, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)

	withdrawal.CreatedTimestamp = withdrawal.CreatedTimestamp.Truncate(time.Microsecond)
	withdrawal.WithdrawnTimestamp = withdrawal.WithdrawnTimestamp.Truncate(time.Microsecond)

	want = append(want, *withdrawal)

	block = source.getNextBlock(t, ctx)
	withdrawal, err = entities.WithdrawalFromProto(withdrawalProto2, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)

	block = source.getNextBlock(t, ctx)
	withdrawal, err = entities.WithdrawalFromProto(withdrawalProto2, generateTxHash(), block.FuryTime)
	withdrawalProto2.Status = fury.Withdrawal_STATUS_FINALIZED
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)

	withdrawal.CreatedTimestamp = withdrawal.CreatedTimestamp.Truncate(time.Microsecond)
	withdrawal.WithdrawnTimestamp = withdrawal.WithdrawnTimestamp.Truncate(time.Microsecond)

	want = append(want, *withdrawal)

	got, _, _ := ws.GetByParty(ctx, withdrawalProto1.PartyId, false, entities.CursorPagination{}, entities.DateRange{})

	assert.Equal(t, want, got)
}

func testWithdrawalsGetByTxHash(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, ws, conn := setupWithdrawalStoreTests(t)

	var rowCount int

	err := conn.QueryRow(ctx, `select count(*) from withdrawals`).Scan(&rowCount)
	require.NoError(t, err)
	assert.Equal(t, 0, rowCount)

	block := addTestBlock(t, ctx, bs)
	withdrawalProto1 := getTestWithdrawal(testID, testID, testID, testAmount, testID, block.FuryTime)
	withdrawalProto1.Id = "deadbeef01"

	withdrawal, err := entities.WithdrawalFromProto(withdrawalProto1, generateTxHash(), block.FuryTime)
	require.NoError(t, err, "Converting withdrawal proto to database entity")

	err = ws.Upsert(ctx, withdrawal)
	require.NoError(t, err)

	withdrawals, err := ws.GetByTxHash(ctx, withdrawal.TxHash)
	require.NoError(t, err)
	require.Equal(t, *withdrawal, withdrawals[0])
}

func getTestWithdrawal(id, party, asset, amount, txHash string, ts time.Time) *fury.Withdrawal {
	return &fury.Withdrawal{
		Id:                 id,
		PartyId:            party,
		Amount:             amount,
		Asset:              asset,
		Status:             fury.Withdrawal_STATUS_OPEN,
		Ref:                "deadbeef",
		TxHash:             txHash,
		CreatedTimestamp:   ts.UnixNano(),
		WithdrawnTimestamp: ts.UnixNano(),
		Ext: &fury.WithdrawExt{
			Ext: &fury.WithdrawExt_Erc20{
				Erc20: &fury.Erc20WithdrawExt{
					ReceiverAddress: "0x1234",
				},
			},
		},
	}
}

func addWithdrawals(ctx context.Context, t *testing.T, bs *sqlstore.Blocks, ws *sqlstore.Withdrawals) []entities.Withdrawal {
	t.Helper()
	furyTime := time.Now().Truncate(time.Microsecond)
	amount := int64(1000)
	withdrawals := make([]entities.Withdrawal, 0, 10)
	for i := 0; i < 10; i++ {
		addTestBlockForTime(t, ctx, bs, furyTime)

		withdrawalProto := getTestWithdrawal(fmt.Sprintf("deadbeef%02d", i+1), testID, testID,
			strconv.FormatInt(amount, 10), helpers.GenerateID(), furyTime)
		withdrawal, err := entities.WithdrawalFromProto(withdrawalProto, generateTxHash(), furyTime)
		require.NoError(t, err, "Converting withdrawal proto to database entity")
		err = ws.Upsert(ctx, withdrawal)
		withdrawals = append(withdrawals, *withdrawal)
		require.NoError(t, err)

		furyTime = furyTime.Add(time.Second)
		amount += 100
	}

	return withdrawals
}

func testWithdrawalsPaginationNoPagination(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)

	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	assert.Equal(t, testWithdrawals, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[0].FuryTime,
			ID:       testWithdrawals[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[9].FuryTime,
			ID:       testWithdrawals[9].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)

	first := int32(3)
	pagination, err := entities.NewCursorPagination(&first, nil, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[:3]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[0].FuryTime,
			ID:       testWithdrawals[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[2].FuryTime,
			ID:       testWithdrawals[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationLast(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)

	last := int32(3)
	pagination, err := entities.NewCursorPagination(nil, nil, &last, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[7:]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[7].FuryTime,
			ID:       testWithdrawals[7].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[9].FuryTime,
			ID:       testWithdrawals[9].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationFirstAfter(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)

	first := int32(3)
	after := testWithdrawals[2].Cursor().Encode()
	pagination, err := entities.NewCursorPagination(&first, &after, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[3:6]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[3].FuryTime,
			ID:       testWithdrawals[3].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[5].FuryTime,
			ID:       testWithdrawals[5].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationLastBefore(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)

	last := int32(3)
	before := entities.NewCursor(entities.WithdrawalCursor{
		FuryTime: testWithdrawals[7].FuryTime,
		ID:       testWithdrawals[7].ID,
	}.String()).Encode()
	pagination, err := entities.NewCursorPagination(nil, nil, &last, &before, false)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[4:7]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[4].FuryTime,
			ID:       testWithdrawals[4].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[6].FuryTime,
			ID:       testWithdrawals[6].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationNoPaginationNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := entities.ReverseSlice(addWithdrawals(ctx, t, bs, ws))

	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, true)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	assert.Equal(t, testWithdrawals, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[0].FuryTime,
			ID:       testWithdrawals[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[9].FuryTime,
			ID:       testWithdrawals[9].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationFirstNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := entities.ReverseSlice(addWithdrawals(ctx, t, bs, ws))

	first := int32(3)
	pagination, err := entities.NewCursorPagination(&first, nil, nil, nil, true)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[:3]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[0].FuryTime,
			ID:       testWithdrawals[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[2].FuryTime,
			ID:       testWithdrawals[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationLastNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := entities.ReverseSlice(addWithdrawals(ctx, t, bs, ws))

	last := int32(3)
	pagination, err := entities.NewCursorPagination(nil, nil, &last, nil, true)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[7:]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[7].FuryTime,
			ID:       testWithdrawals[7].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[9].FuryTime,
			ID:       testWithdrawals[9].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationFirstAfterNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := entities.ReverseSlice(addWithdrawals(ctx, t, bs, ws))

	first := int32(3)
	after := testWithdrawals[2].Cursor().Encode()
	pagination, err := entities.NewCursorPagination(&first, &after, nil, nil, true)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[3:6]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[3].FuryTime,
			ID:       testWithdrawals[3].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[5].FuryTime,
			ID:       testWithdrawals[5].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationLastBeforeNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := entities.ReverseSlice(addWithdrawals(ctx, t, bs, ws))

	last := int32(3)
	before := entities.NewCursor(entities.WithdrawalCursor{
		FuryTime: testWithdrawals[7].FuryTime,
		ID:       testWithdrawals[7].ID,
	}.String()).Encode()
	pagination, err := entities.NewCursorPagination(nil, nil, &last, &before, true)
	require.NoError(t, err)
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{})

	require.NoError(t, err)
	want := testWithdrawals[4:7]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[4].FuryTime,
			ID:       testWithdrawals[4].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: testWithdrawals[6].FuryTime,
			ID:       testWithdrawals[6].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesNoPagination(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := testWithdrawals[2:8]

	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, false)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[5].FuryTime,
			ID:       want[5].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := testWithdrawals[2:8]

	first := int32(3)
	pagination, err := entities.NewCursorPagination(&first, nil, nil, nil, false)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[:3]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesLast(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := testWithdrawals[2:8]

	last := int32(3)
	pagination, err := entities.NewCursorPagination(nil, nil, &last, nil, false)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[3:]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesFirstAfter(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := testWithdrawals[2:8]

	first := int32(3)
	after := want[1].Cursor().Encode()
	pagination, err := entities.NewCursorPagination(&first, &after, nil, nil, false)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[2:5]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesLastBefore(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := testWithdrawals[2:8]

	last := int32(3)
	before := entities.NewCursor(entities.WithdrawalCursor{
		FuryTime: want[4].FuryTime,
		ID:       want[4].ID,
	}.String()).Encode()
	pagination, err := entities.NewCursorPagination(nil, nil, &last, &before, false)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[1:4]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesNoPaginationNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)

	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, true)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	want := entities.ReverseSlice(testWithdrawals[2:8])

	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[5].FuryTime,
			ID:       want[5].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesFirstNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := entities.ReverseSlice(testWithdrawals[2:8])

	first := int32(3)
	pagination, err := entities.NewCursorPagination(&first, nil, nil, nil, true)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[:3]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: false,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesLastNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := entities.ReverseSlice(testWithdrawals[2:8])

	last := int32(3)
	pagination, err := entities.NewCursorPagination(nil, nil, &last, nil, true)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[3:]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesFirstAfterNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := entities.ReverseSlice(testWithdrawals[2:8])

	first := int32(3)
	after := want[1].Cursor().Encode()
	pagination, err := entities.NewCursorPagination(&first, &after, nil, nil, true)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[2:5]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}

func testWithdrawalsPaginationBetweenDatesLastBeforeNewestFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, ws, _ := setupWithdrawalStoreTests(t)

	testWithdrawals := addWithdrawals(ctx, t, bs, ws)
	want := entities.ReverseSlice(testWithdrawals[2:8])

	last := int32(3)
	before := entities.NewCursor(entities.WithdrawalCursor{
		FuryTime: want[4].FuryTime,
		ID:       want[4].ID,
	}.String()).Encode()
	pagination, err := entities.NewCursorPagination(nil, nil, &last, &before, true)
	require.NoError(t, err)
	startDate := testWithdrawals[2].FuryTime
	endDate := testWithdrawals[8].FuryTime
	got, pageInfo, err := ws.GetByParty(ctx, testID, false, pagination, entities.DateRange{
		Start: &startDate,
		End:   &endDate,
	})

	require.NoError(t, err)
	want = want[1:4]
	assert.Equal(t, want, got)
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[0].FuryTime,
			ID:       want[0].ID,
		}.String()).Encode(),
		EndCursor: entities.NewCursor(entities.WithdrawalCursor{
			FuryTime: want[2].FuryTime,
			ID:       want[2].ID,
		}.String()).Encode(),
	}, pageInfo)
}
