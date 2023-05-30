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

package sqlsubscribers

import (
	"context"

	"github.com/pkg/errors"

	"github.com/elysiumstation/fury/core/events"
	"github.com/elysiumstation/fury/datanode/entities"
	eventspb "github.com/elysiumstation/fury/protos/fury/events/v1"
)

type EpochUpdateEvent interface {
	events.Event
	Proto() eventspb.EpochEvent
}

type EpochStore interface {
	Add(context.Context, entities.Epoch) error
}

type Epoch struct {
	subscriber
	store EpochStore
}

func NewEpoch(store EpochStore) *Epoch {
	t := &Epoch{
		store: store,
	}
	return t
}

func (es *Epoch) Types() []events.Type {
	return []events.Type{events.EpochUpdate}
}

func (es *Epoch) Push(ctx context.Context, evt events.Event) error {
	return es.consume(ctx, evt.(EpochUpdateEvent))
}

func (es *Epoch) consume(ctx context.Context, event EpochUpdateEvent) error {
	epochUpdateEvent := event.Proto()
	epoch := entities.EpochFromProto(epochUpdateEvent, entities.TxHash(event.TxHash()))
	epoch.FuryTime = es.furyTime

	return errors.Wrap(es.store.Add(ctx, epoch), "error adding epoch update")
}
