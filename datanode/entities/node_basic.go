package entities

import (
	"time"

	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
	"github.com/elysiumstation/fury/protos/fury"
)

type NodeBasic struct {
	ID              NodeID
	PubKey          FuryPublicKey       `db:"fury_pub_key"`
	TmPubKey        TendermintPublicKey `db:"tendermint_pub_key"`
	EthereumAddress EthereumAddress
	InfoURL         string
	Location        string
	Status          NodeStatus
	Name            string
	AvatarURL       string
	TxHash          TxHash
	FuryTime        time.Time
}

func (n NodeBasic) ToProto() *v2.NodeBasic {
	return &v2.NodeBasic{
		Id:              n.ID.String(),
		PubKey:          n.PubKey.String(),
		TmPubKey:        n.TmPubKey.String(),
		EthereumAddress: n.EthereumAddress.String(),
		InfoUrl:         n.InfoURL,
		Location:        n.Location,
		Status:          fury.NodeStatus(n.Status),
		Name:            n.Name,
		AvatarUrl:       n.AvatarURL,
	}
}
