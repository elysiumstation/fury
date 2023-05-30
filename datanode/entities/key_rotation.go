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

package entities

import (
	"encoding/json"
	"fmt"
	"time"

	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
	eventspb "github.com/elysiumstation/fury/protos/fury/events/v1"
)

type KeyRotation struct {
	NodeID      NodeID
	OldPubKey   FuryPublicKey
	NewPubKey   FuryPublicKey
	BlockHeight uint64
	TxHash      TxHash
	FuryTime    time.Time
}

func KeyRotationFromProto(kr *eventspb.KeyRotation, txHash TxHash, furyTime time.Time) (*KeyRotation, error) {
	return &KeyRotation{
		NodeID:      NodeID(kr.NodeId),
		OldPubKey:   FuryPublicKey(kr.OldPubKey),
		NewPubKey:   FuryPublicKey(kr.NewPubKey),
		BlockHeight: kr.BlockHeight,
		TxHash:      txHash,
		FuryTime:    furyTime,
	}, nil
}

func (kr KeyRotation) ToProto() *eventspb.KeyRotation {
	return &eventspb.KeyRotation{
		NodeId:      kr.NodeID.String(),
		OldPubKey:   kr.OldPubKey.String(),
		NewPubKey:   kr.NewPubKey.String(),
		BlockHeight: kr.BlockHeight,
	}
}

func (kr KeyRotation) Cursor() *Cursor {
	cursor := KeyRotationCursor{
		FuryTime:  kr.FuryTime,
		NodeID:    kr.NodeID,
		OldPubKey: kr.OldPubKey,
		NewPubKey: kr.NewPubKey,
	}

	return NewCursor(cursor.String())
}

func (kr KeyRotation) ToProtoEdge(_ ...any) (*v2.KeyRotationEdge, error) {
	return &v2.KeyRotationEdge{
		Node:   kr.ToProto(),
		Cursor: kr.Cursor().Encode(),
	}, nil
}

type KeyRotationCursor struct {
	FuryTime  time.Time     `json:"fury_time"`
	NodeID    NodeID        `json:"node_id"`
	OldPubKey FuryPublicKey `json:"old_pub_key"`
	NewPubKey FuryPublicKey `json:"new_pub_key"`
}

func (c KeyRotationCursor) String() string {
	bs, err := json.Marshal(c)
	// This should never fail so if it does, we should panic
	if err != nil {
		panic(fmt.Errorf("could not marshal key rotation cursor: %w", err))
	}

	return string(bs)
}

func (c *KeyRotationCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}

	return json.Unmarshal([]byte(cursorString), c)
}
