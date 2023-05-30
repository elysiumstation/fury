package snapshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractIntervalFromViewDefinition(t *testing.T) {
	viewDefinition := ` SELECT balances.account_id,
	time_bucket('01:00:00'::interval, balances.fury_time) AS bucket,
		last(balances.balance, balances.fury_time) AS balance,
		last(balances.tx_hash, balances.fury_time) AS tx_hash,
		last(balances.fury_time, balances.fury_time) AS fury_time
	FROM balances
	GROUP BY balances.account_id, (time_bucket('01:00:00'::interval, balances.fury_time));`

	interval, err := extractIntervalFromViewDefinition(viewDefinition)
	require.NoError(t, err)
	assert.Equal(t, "01:00:00", interval)

	viewDefinition = ` SELECT balances.account_id,
	time_bucket('1 day'::interval, balances.fury_time) AS bucket,
		last(balances.balance, balances.fury_time) AS balance,
		last(balances.tx_hash, balances.fury_time) AS tx_hash,
		last(balances.fury_time, balances.fury_time) AS fury_time
	FROM balances
	GROUP BY balances.account_id, (time_bucket('1 day'::interval, balances.fury_time));`

	interval, err = extractIntervalFromViewDefinition(viewDefinition)
	require.NoError(t, err)
	assert.Equal(t, "1 day", interval)
}
