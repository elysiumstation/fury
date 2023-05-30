// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.FURY file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package services

import (
	"context"
	"sync"

	"github.com/elysiumstation/fury/core/events"
	"github.com/elysiumstation/fury/core/subscribers"
	"github.com/elysiumstation/fury/protos/fury"
)

type netLimitsEvent interface {
	events.Event
	NetworkLimits() *fury.NetworkLimits
}

type NetLimits struct {
	*subscribers.Base
	ctx    context.Context
	limits fury.NetworkLimits
	ch     chan fury.NetworkLimits
	mu     sync.RWMutex
}

func NewNetLimits(ctx context.Context) (netLimits *NetLimits) {
	defer func() { go netLimits.consume() }()
	return &NetLimits{
		Base: subscribers.NewBase(ctx, 1000, true),
		ctx:  ctx,
		ch:   make(chan fury.NetworkLimits, 100),
	}
}

func (n *NetLimits) consume() {
	defer func() { close(n.ch) }()
	for {
		select {
		case <-n.Closed():
			return
		case limits, ok := <-n.ch:
			if !ok {
				n.Halt()
				return
			}
			n.mu.Lock()
			n.limits = limits
			n.mu.Unlock()
		}
	}
}

func (n *NetLimits) Get() *fury.NetworkLimits {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.limits.DeepClone()
}

func (n *NetLimits) Push(evts ...events.Event) {
	for _, e := range evts {
		if ne, ok := e.(netLimitsEvent); ok {
			n.ch <- *ne.NetworkLimits()
		}
	}
}

func (n *NetLimits) Types() []events.Type {
	return []events.Type{events.NetworkLimitsEvent}
}
