package entities

import (
	"encoding/hex"
	"time"

	"github.com/elysiumstation/fury/core/events"
	eventspb "github.com/elysiumstation/fury/protos/fury/events/v1"
	"github.com/pkg/errors"
)

type BeginBlockEvent interface {
	events.Event
	BeginBlock() eventspb.BeginBlock
}

func BlockFromBeginBlock(b BeginBlockEvent) (*Block, error) {
	hash, err := hex.DecodeString(b.TraceID())
	if err != nil {
		return nil, errors.Wrapf(err, "Trace ID is not valid hex string, trace ID:%s", b.TraceID())
	}

	furyTime := time.Unix(0, b.BeginBlock().Timestamp)

	// Postgres only stores timestamps in microsecond resolution
	block := Block{
		FuryTime: furyTime.Truncate(time.Microsecond),
		Hash:     hash,
		Height:   b.BlockNr(),
	}
	return &block, err
}
