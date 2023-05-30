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

//lint:file-ignore ST1003 Ignore underscores in names, this is straight copied from the proto package to ease introducing the domain types

package types

import (
	"errors"
	"fmt"

	"github.com/elysiumstation/fury/libs/crypto"
	"github.com/elysiumstation/fury/libs/num"
	furypb "github.com/elysiumstation/fury/protos/fury"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
)

type WithdrawExt struct {
	Ext isWithdrawExtExt
}

func (x *WithdrawExt) String() string {
	return fmt.Sprintf(
		"ext(%s)",
		reflectPointerToString(x.Ext),
	)
}

func (x *WithdrawExt) IntoProto() *furypb.WithdrawExt {
	if x == nil {
		return nil
	}
	switch st := x.Ext.(type) {
	case *WithdrawExtErc20:
		return &furypb.WithdrawExt{
			Ext: st.IntoProto(),
		}
	default:
		return nil
	}
}

func (x *WithdrawExt) GetErc20() *WithdrawExtErc20 {
	switch st := x.Ext.(type) {
	case *WithdrawExtErc20:
		return st
	default:
		return nil
	}
}

func WithdrawExtFromProto(extProto *furypb.WithdrawExt) *WithdrawExt {
	if extProto == nil {
		return nil
	}
	var src isWithdrawExtExt
	switch st := extProto.Ext.(type) {
	case *furypb.WithdrawExt_Erc20:
		src = WithdrawExtErc20FromProto(st)
	}
	return &WithdrawExt{
		Ext: src,
	}
}

type isWithdrawExtExt interface {
	isWithdrawExtExt()
	String() string
}

type WithdrawExtErc20 struct {
	Erc20 *Erc20WithdrawExt
}

func (x *WithdrawExtErc20) isWithdrawExtExt() {}

func (x *WithdrawExtErc20) String() string {
	return fmt.Sprintf(
		"erc20(%s)",
		reflectPointerToString(x.Erc20),
	)
}

func (x *WithdrawExtErc20) IntoProto() *furypb.WithdrawExt_Erc20 {
	return &furypb.WithdrawExt_Erc20{
		Erc20: x.Erc20.IntoProto(),
	}
}

func (x *WithdrawExtErc20) GetReceiverAddress() string {
	return x.Erc20.ReceiverAddress
}

func WithdrawExtErc20FromProto(erc20 *furypb.WithdrawExt_Erc20) *WithdrawExtErc20 {
	return &WithdrawExtErc20{
		Erc20: Erc20WithdrawExtFromProto(erc20.Erc20),
	}
}

type Erc20WithdrawExt struct {
	ReceiverAddress string
}

func (x *Erc20WithdrawExt) String() string {
	return fmt.Sprintf("receiverAddress(%s)", x.ReceiverAddress)
}

func (x *Erc20WithdrawExt) IntoProto() *furypb.Erc20WithdrawExt {
	return &furypb.Erc20WithdrawExt{
		ReceiverAddress: x.ReceiverAddress,
	}
}

func Erc20WithdrawExtFromProto(erc20 *furypb.Erc20WithdrawExt) *Erc20WithdrawExt {
	return &Erc20WithdrawExt{
		ReceiverAddress: crypto.EthereumChecksumAddress(erc20.ReceiverAddress),
	}
}

type WithdrawalStatus = furypb.Withdrawal_Status

const (
	// WithdrawalStatusUnspecified Default value, always invalid.
	WithdrawalStatusUnspecified WithdrawalStatus = 0
	// WithdrawalStatusOpen The withdrawal is open and being processed by the network.
	WithdrawalStatusOpen WithdrawalStatus = 1
	// WithdrawalStatusRejected The withdrawal have been rejected.
	WithdrawalStatusRejected WithdrawalStatus = 2
	// WithdrawalStatusFinalized The withdrawal went through and is fully finalised, the funds are removed from the
	// Fury network and are unlocked on the foreign chain bridge, for example, on the Ethereum network.
	WithdrawalStatusFinalized WithdrawalStatus = 3
)

type Withdrawal struct {
	// ID Unique identifier for the withdrawal
	ID string
	// PartyID Unique party identifier of the user initiating the withdrawal
	PartyID string
	// Amount The amount to be withdrawn
	Amount *num.Uint
	// Asset The asset we want to withdraw funds from
	Asset string
	// Status The status of the withdrawal
	Status WithdrawalStatus
	// Ref The reference which is used by the foreign chain
	// to refer to this withdrawal
	Ref string
	// TxHash The hash of the foreign chain for this transaction
	TxHash string
	// CreationDate Timestamp for when the network started to process this withdrawal
	CreationDate int64
	// WithdrawalDate Timestamp for when the withdrawal was finalised by the network
	WithdrawalDate int64
	// Ext Foreign chain specifics
	Ext *WithdrawExt
}

func (w *Withdrawal) IntoProto() *furypb.Withdrawal {
	return &furypb.Withdrawal{
		Id:                 w.ID,
		PartyId:            w.PartyID,
		Amount:             num.UintToString(w.Amount),
		Asset:              w.Asset,
		Status:             w.Status,
		Ref:                w.Ref,
		TxHash:             w.TxHash,
		CreatedTimestamp:   w.CreationDate,
		WithdrawnTimestamp: w.WithdrawalDate,
		Ext:                w.Ext.IntoProto(),
	}
}

func WithdrawalFromProto(w *furypb.Withdrawal) *Withdrawal {
	amt, _ := num.UintFromString(w.Amount, 10)
	return &Withdrawal{
		ID:             w.Id,
		PartyID:        w.PartyId,
		Amount:         amt,
		Asset:          w.Asset,
		Status:         w.Status,
		Ref:            w.Ref,
		TxHash:         w.TxHash,
		CreationDate:   w.CreatedTimestamp,
		WithdrawalDate: w.WithdrawnTimestamp,
		Ext:            WithdrawExtFromProto(w.Ext),
	}
}

type DepositStatus = furypb.Deposit_Status

const (
	// DepositStatusUnspecified Default value, always invalid.
	DepositStatusUnspecified DepositStatus = 0
	// DepositStatusOpen The deposit is being processed by the network.
	DepositStatusOpen DepositStatus = 1
	// DepositStatusCancelled The deposit has been cancelled by the network.
	DepositStatusCancelled DepositStatus = 2
	// DepositStatusFinalized The deposit has been finalised and accounts have been updated.
	DepositStatusFinalized DepositStatus = 3
)

// Deposit represent a deposit on to the Fury network.
type Deposit struct {
	// ID Unique identifier for the deposit
	ID string
	// Status of the deposit
	Status DepositStatus
	// Party identifier of the user initiating the deposit
	PartyID string
	// Asset The Fury asset targeted by this deposit
	Asset string
	// Amount The amount to be deposited
	Amount *num.Uint
	// TxHash The hash of the transaction from the foreign chain
	TxHash string
	// Timestamp for when the Fury account was updated with the deposit
	CreditDate int64
	// Timestamp for when the deposit was created on the Fury network
	CreationDate int64
}

func (d *Deposit) IntoProto() *furypb.Deposit {
	return &furypb.Deposit{
		Id:                d.ID,
		Status:            d.Status,
		PartyId:           d.PartyID,
		Asset:             d.Asset,
		Amount:            num.UintToString(d.Amount),
		TxHash:            d.TxHash,
		CreditedTimestamp: d.CreditDate,
		CreatedTimestamp:  d.CreationDate,
	}
}

func (d *Deposit) String() string {
	return fmt.Sprintf(
		"ID(%s) status(%s) partyID(%s) asset(%s) amount(%s) txHash(%s) creditDate(%v) creationDate(%v)",
		d.ID,
		d.Status.String(),
		d.PartyID,
		d.Asset,
		uintPointerToString(d.Amount),
		d.TxHash,
		d.CreditDate,
		d.CreationDate,
	)
}

func DepositFromProto(d *furypb.Deposit) *Deposit {
	amt, _ := num.UintFromString(d.Amount, 10)
	return &Deposit{
		ID:           d.Id,
		Status:       d.Status,
		PartyID:      d.PartyId,
		Asset:        d.Asset,
		Amount:       amt,
		TxHash:       d.TxHash,
		CreditDate:   d.CreditedTimestamp,
		CreationDate: d.CreatedTimestamp,
	}
}

type ChainEventERC20 struct {
	ERC20 *ERC20Event
}

func NewChainEventERC20FromProto(p *commandspb.ChainEvent_Erc20) (*ChainEventERC20, error) {
	c := ChainEventERC20{}
	var err error
	c.ERC20, err = NewERC20Event(p.Erc20)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c ChainEventERC20) IntoProto() *commandspb.ChainEvent_Erc20 {
	p := &commandspb.ChainEvent_Erc20{
		Erc20: c.ERC20.IntoProto(),
	}
	return p
}

func (c ChainEventERC20) String() string {
	return fmt.Sprintf(
		"erc20(%s)",
		reflectPointerToString(c.ERC20),
	)
}

type BuiltinAssetDeposit struct {
	// A Fury network internal asset identifier
	FuryAssetID string
	// A Fury party identifier (pub-key)
	PartyID string
	// The amount to be deposited
	Amount *num.Uint
}

func NewBuiltinAssetDepositFromProto(p *furypb.BuiltinAssetDeposit) (*BuiltinAssetDeposit, error) {
	amount := num.UintZero()
	if len(p.Amount) > 0 {
		var overflowed bool
		amount, overflowed = num.UintFromString(p.Amount, 10)
		if overflowed {
			return nil, errors.New("invalid amount")
		}
	}
	return &BuiltinAssetDeposit{
		FuryAssetID: p.FuryAssetId,
		PartyID:     p.PartyId,
		Amount:      amount,
	}, nil
}

func (b BuiltinAssetDeposit) IntoProto() *furypb.BuiltinAssetDeposit {
	return &furypb.BuiltinAssetDeposit{
		FuryAssetId: b.FuryAssetID,
		PartyId:     b.PartyID,
		Amount:      num.UintToString(b.Amount),
	}
}

func (b BuiltinAssetDeposit) String() string {
	return fmt.Sprintf(
		"party(%s) furyAssetID(%s) amount(%s)",
		b.PartyID,
		b.FuryAssetID,
		uintPointerToString(b.Amount),
	)
}

func (b BuiltinAssetDeposit) GetFuryAssetID() string {
	return b.FuryAssetID
}

type BuiltinAssetWithdrawal struct {
	// A Fury network internal asset identifier
	FuryAssetID string
	// A Fury network party identifier (pub-key)
	PartyID string
	// The amount to be withdrawn
	Amount *num.Uint
}

func NewBuiltinAssetWithdrawalFromProto(p *furypb.BuiltinAssetWithdrawal) (*BuiltinAssetWithdrawal, error) {
	amount := num.UintZero()
	if len(p.Amount) > 0 {
		var overflowed bool
		amount, overflowed = num.UintFromString(p.Amount, 10)
		if overflowed {
			return nil, errors.New("invalid amount")
		}
	}
	return &BuiltinAssetWithdrawal{
		FuryAssetID: p.FuryAssetId,
		PartyID:     p.PartyId,
		Amount:      amount,
	}, nil
}

func (b BuiltinAssetWithdrawal) IntoProto() *furypb.BuiltinAssetWithdrawal {
	return &furypb.BuiltinAssetWithdrawal{
		FuryAssetId: b.FuryAssetID,
		PartyId:     b.PartyID,
		Amount:      num.UintToString(b.Amount),
	}
}

func (b BuiltinAssetWithdrawal) String() string {
	return fmt.Sprintf(
		"partyID(%s) furyAssetID(%s) amount(%s)",
		b.PartyID,
		b.FuryAssetID,
		uintPointerToString(b.Amount),
	)
}

func (b BuiltinAssetWithdrawal) GetFuryAssetID() string {
	return b.FuryAssetID
}

type ChainEventBuiltin struct {
	Builtin *BuiltinAssetEvent
}

func NewChainEventBuiltinFromProto(p *commandspb.ChainEvent_Builtin) (*ChainEventBuiltin, error) {
	c := ChainEventBuiltin{}
	var err error
	c.Builtin, err = NewBuiltinAssetEventFromProto(p.Builtin)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c ChainEventBuiltin) IntoProto() *commandspb.ChainEvent_Builtin {
	ceb := &commandspb.ChainEvent_Builtin{
		Builtin: c.Builtin.IntoProto(),
	}
	return ceb
}

func (c ChainEventBuiltin) String() string {
	return fmt.Sprintf(
		"builtin(%s)",
		reflectPointerToString(c.Builtin),
	)
}

type BuiltinAssetEvent struct {
	// Types that are valid to be assigned to Action:
	//	*BuiltinAssetEvent_Deposit
	//	*BuiltinAssetEvent_Withdrawal
	Action builtinAssetEventAction
}

type builtinAssetEventAction interface {
	isBuiltinAssetEvent()
	oneOfProto() interface{}
	String() string
}

func NewBuiltinAssetEventFromProto(p *furypb.BuiltinAssetEvent) (*BuiltinAssetEvent, error) {
	var (
		ae  = &BuiltinAssetEvent{}
		err error
	)
	switch e := p.Action.(type) {
	case *furypb.BuiltinAssetEvent_Deposit:
		ae.Action, err = NewBuiltinAssetEventDeposit(e)
		return ae, err
	case *furypb.BuiltinAssetEvent_Withdrawal:
		ae.Action, err = NewBuiltinAssetEventWithdrawal(e)
		return ae, err
	default:
		return nil, errors.New("unknown asset event type")
	}
}

func (c BuiltinAssetEvent) IntoProto() *furypb.BuiltinAssetEvent {
	action := c.Action.oneOfProto()
	ceb := &furypb.BuiltinAssetEvent{}
	switch a := action.(type) {
	case *furypb.BuiltinAssetEvent_Deposit:
		ceb.Action = a
	case *furypb.BuiltinAssetEvent_Withdrawal:
		ceb.Action = a
	}
	return ceb
}

func (c BuiltinAssetEvent) String() string {
	return fmt.Sprintf(
		"action(%s)",
		reflectPointerToString(c.Action),
	)
}

type BuiltinAssetEventDeposit struct {
	Deposit *BuiltinAssetDeposit
}

func (b BuiltinAssetEventDeposit) String() string {
	return fmt.Sprintf(
		"deposit(%s)",
		reflectPointerToString(b.Deposit),
	)
}

func NewBuiltinAssetEventDeposit(p *furypb.BuiltinAssetEvent_Deposit) (*BuiltinAssetEventDeposit, error) {
	dep, err := NewBuiltinAssetDepositFromProto(p.Deposit)
	if err != nil {
		return nil, err
	}
	return &BuiltinAssetEventDeposit{
		Deposit: dep,
	}, nil
}

func (b BuiltinAssetEventDeposit) IntoProto() *furypb.BuiltinAssetEvent_Deposit {
	p := &furypb.BuiltinAssetEvent_Deposit{
		Deposit: b.Deposit.IntoProto(),
	}
	return p
}

func (b BuiltinAssetEventDeposit) isBuiltinAssetEvent() {}

func (b BuiltinAssetEventDeposit) oneOfProto() interface{} {
	return b.IntoProto()
}

type BuiltinAssetEventWithdrawal struct {
	Withdrawal *BuiltinAssetWithdrawal
}

func (b BuiltinAssetEventWithdrawal) String() string {
	return fmt.Sprintf(
		"withdrawal(%s)",
		reflectPointerToString(b.Withdrawal),
	)
}

func NewBuiltinAssetEventWithdrawal(p *furypb.BuiltinAssetEvent_Withdrawal) (*BuiltinAssetEventWithdrawal, error) {
	withdrawal, err := NewBuiltinAssetWithdrawalFromProto(p.Withdrawal)
	if err != nil {
		return nil, err
	}
	return &BuiltinAssetEventWithdrawal{
		Withdrawal: withdrawal,
	}, nil
}

func (b BuiltinAssetEventWithdrawal) IntoProto() *furypb.BuiltinAssetEvent_Withdrawal {
	p := &furypb.BuiltinAssetEvent_Withdrawal{
		Withdrawal: b.Withdrawal.IntoProto(),
	}
	return p
}

func (b BuiltinAssetEventWithdrawal) isBuiltinAssetEvent() {}

func (b BuiltinAssetEventWithdrawal) oneOfProto() interface{} {
	return b.IntoProto()
}

type ERC20Event struct {
	// Index of the transaction
	Index uint64
	// The block in which the transaction was added
	Block uint64
	// The action
	//
	// Types that are valid to be assigned to Action:
	//	*ERC20EventAssetList
	//	*ERC20EventAssetDelist
	//	*ERC20EventDeposit
	//	*ERC20EventWithdrawal
	//	*ERC20EventAssetLimitsUpdated
	//	*ERC20BridgeStopped
	//	*ERC20BridgeRemoved
	Action erc20EventAction
}

type erc20EventAction interface {
	isErc20EventAction()
	oneOfProto() interface{}
	String() string
}

func NewERC20Event(p *furypb.ERC20Event) (*ERC20Event, error) {
	e := ERC20Event{
		Index: p.Index,
		Block: p.Block,
	}

	var err error
	switch a := p.Action.(type) {
	case *furypb.ERC20Event_AssetDelist:
		e.Action = NewERC20EventAssetDelist(a)
		return &e, nil
	case *furypb.ERC20Event_AssetList:
		e.Action = NewERC20EventAssetList(a)
		return &e, nil
	case *furypb.ERC20Event_Deposit:
		e.Action, err = NewERC20EventDeposit(a)
		if err != nil {
			return nil, err
		}
		return &e, nil
	case *furypb.ERC20Event_Withdrawal:
		e.Action = NewERC20EventWithdrawal(a)
		return &e, nil
	case *furypb.ERC20Event_AssetLimitsUpdated:
		e.Action = NewERC20EventAssetLimitsUpdated(a)
		return &e, nil
	case *furypb.ERC20Event_BridgeStopped:
		e.Action = NewERC20EventBridgeStopped(a)
		return &e, nil
	case *furypb.ERC20Event_BridgeResumed:
		e.Action = NewERC20EventBridgeResumed(a)
		return &e, nil
	default:
		return nil, errors.New("unknown erc20 event type")
	}
}

func (e ERC20Event) IntoProto() *furypb.ERC20Event {
	p := &furypb.ERC20Event{
		Index: e.Index,
		Block: e.Block,
	}

	switch a := e.Action.(type) {
	case *ERC20EventAssetDelist:
		p.Action = a.IntoProto()
	case *ERC20EventAssetList:
		p.Action = a.IntoProto()
	case *ERC20EventDeposit:
		p.Action = a.IntoProto()
	case *ERC20EventWithdrawal:
		p.Action = a.IntoProto()
	default:
		return nil
	}

	return p
}

func (e ERC20Event) String() string {
	return fmt.Sprintf(
		"index(%v) block(%v) action(%s)",
		e.Index,
		e.Block,
		reflectPointerToString(e.Action),
	)
}

type ERC20EventAssetDelist struct {
	AssetDelist *ERC20AssetDelist
}

func (e ERC20EventAssetDelist) String() string {
	return fmt.Sprintf(
		"assetDelist(%s)",
		reflectPointerToString(e.AssetDelist),
	)
}

func (ERC20EventAssetDelist) isErc20EventAction() {}

func (e ERC20EventAssetDelist) oneOfProto() interface{} {
	return e.AssetDelist.IntoProto()
}

func NewERC20EventAssetDelist(p *furypb.ERC20Event_AssetDelist) *ERC20EventAssetDelist {
	return &ERC20EventAssetDelist{
		AssetDelist: NewERC20AssetDelistFromProto(p.AssetDelist),
	}
}

func (e ERC20EventAssetDelist) IntoProto() *furypb.ERC20Event_AssetDelist {
	return &furypb.ERC20Event_AssetDelist{
		AssetDelist: e.AssetDelist.IntoProto(),
	}
}

type ERC20AssetDelist struct {
	// The Fury network internal identifier of the asset
	FuryAssetID string
}

func NewERC20AssetDelistFromProto(p *furypb.ERC20AssetDelist) *ERC20AssetDelist {
	return &ERC20AssetDelist{
		FuryAssetID: p.FuryAssetId,
	}
}

func (e ERC20AssetDelist) IntoProto() *furypb.ERC20AssetDelist {
	return &furypb.ERC20AssetDelist{
		FuryAssetId: e.FuryAssetID,
	}
}

func (e ERC20AssetDelist) String() string {
	return fmt.Sprintf("furyAssetID(%s)", e.FuryAssetID)
}

type ERC20EventAssetList struct {
	AssetList *ERC20AssetList
}

func (ERC20EventAssetList) isErc20EventAction() {}

func (e ERC20EventAssetList) oneOfProto() interface{} {
	return e.AssetList.IntoProto()
}

func (e ERC20EventAssetList) String() string {
	return fmt.Sprintf(
		"assetList(%s)",
		reflectPointerToString(e.AssetList),
	)
}

func NewERC20EventAssetList(p *furypb.ERC20Event_AssetList) *ERC20EventAssetList {
	return &ERC20EventAssetList{
		AssetList: NewERC20AssetListFromProto(p.AssetList),
	}
}

func (e ERC20EventAssetList) IntoProto() *furypb.ERC20Event_AssetList {
	return &furypb.ERC20Event_AssetList{
		AssetList: e.AssetList.IntoProto(),
	}
}

type ERC20AssetList struct {
	// The Fury network internal identifier of the asset
	FuryAssetID string
	// ethereum address of the asset
	AssetSource string
}

func NewERC20AssetListFromProto(p *furypb.ERC20AssetList) *ERC20AssetList {
	return &ERC20AssetList{
		FuryAssetID: p.FuryAssetId,
		AssetSource: p.AssetSource,
	}
}

func (e ERC20AssetList) IntoProto() *furypb.ERC20AssetList {
	return &furypb.ERC20AssetList{
		FuryAssetId: e.FuryAssetID,
	}
}

func (e ERC20AssetList) String() string {
	return fmt.Sprintf(
		"furyAssetID(%s)",
		e.FuryAssetID,
	)
}

func (e ERC20AssetList) GetFuryAssetID() string {
	return e.FuryAssetID
}

type ERC20EventWithdrawal struct {
	Withdrawal *ERC20Withdrawal
}

func (ERC20EventWithdrawal) isErc20EventAction() {}

func (e ERC20EventWithdrawal) oneOfProto() interface{} {
	return e.Withdrawal.IntoProto()
}

func (e ERC20EventWithdrawal) String() string {
	return fmt.Sprintf(
		"withdrawal(%s)",
		reflectPointerToString(e.Withdrawal),
	)
}

func NewERC20EventWithdrawal(p *furypb.ERC20Event_Withdrawal) *ERC20EventWithdrawal {
	return &ERC20EventWithdrawal{
		Withdrawal: NewERC20WithdrawalFromProto(p.Withdrawal),
	}
}

func (e ERC20EventWithdrawal) IntoProto() *furypb.ERC20Event_Withdrawal {
	return &furypb.ERC20Event_Withdrawal{
		Withdrawal: e.Withdrawal.IntoProto(),
	}
}

type ERC20Withdrawal struct {
	// The Fury network internal identifier of the asset
	FuryAssetID string
	// The target Ethereum wallet address
	TargetEthereumAddress string
	// The reference nonce used for the transaction
	ReferenceNonce string
}

func NewERC20WithdrawalFromProto(p *furypb.ERC20Withdrawal) *ERC20Withdrawal {
	return &ERC20Withdrawal{
		FuryAssetID:           p.FuryAssetId,
		TargetEthereumAddress: p.TargetEthereumAddress,
		ReferenceNonce:        p.ReferenceNonce,
	}
}

func (e ERC20Withdrawal) IntoProto() *furypb.ERC20Withdrawal {
	return &furypb.ERC20Withdrawal{
		FuryAssetId:           e.FuryAssetID,
		TargetEthereumAddress: e.TargetEthereumAddress,
		ReferenceNonce:        e.ReferenceNonce,
	}
}

func (e ERC20Withdrawal) String() string {
	return fmt.Sprintf(
		"furyAssetID(%s) referenceNonce(%s) targetEthereumAddress(%s)",
		e.FuryAssetID,
		e.ReferenceNonce,
		e.TargetEthereumAddress,
	)
}

func (e ERC20Withdrawal) GetFuryAssetID() string {
	return e.FuryAssetID
}

type ERC20EventDeposit struct {
	Deposit *ERC20Deposit
}

func (e ERC20EventDeposit) String() string {
	return fmt.Sprintf(
		"deposit(%s)",
		reflectPointerToString(e.Deposit),
	)
}

func (ERC20EventDeposit) isErc20EventAction() {}

func (e ERC20EventDeposit) oneOfProto() interface{} {
	return e.Deposit.IntoProto()
}

func NewERC20EventDeposit(p *furypb.ERC20Event_Deposit) (*ERC20EventDeposit, error) {
	e := ERC20EventDeposit{}
	var err error
	e.Deposit, err = NewERC20DepositFromProto(p.Deposit)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (e ERC20EventDeposit) IntoProto() *furypb.ERC20Event_Deposit {
	p := furypb.ERC20Event_Deposit{
		Deposit: e.Deposit.IntoProto(),
	}
	return &p
}

type ERC20Deposit struct {
	// The fury network internal identifier of the asset
	FuryAssetID string
	// The Ethereum wallet that initiated the deposit
	SourceEthereumAddress string
	// The Fury party identifier (pub-key) which is the target of the deposit
	TargetPartyID string
	// The amount to be deposited
	Amount *num.Uint
}

func NewERC20DepositFromProto(p *furypb.ERC20Deposit) (*ERC20Deposit, error) {
	e := ERC20Deposit{
		FuryAssetID:           p.FuryAssetId,
		SourceEthereumAddress: p.SourceEthereumAddress,
		TargetPartyID:         p.TargetPartyId,
	}
	if len(p.Amount) > 0 {
		var failed bool
		e.Amount, failed = num.UintFromString(p.Amount, 10)
		if failed {
			return nil, fmt.Errorf("failed to convert numerical string to Uint: %v", p.Amount)
		}
	}
	return &e, nil
}

func (e ERC20Deposit) IntoProto() *furypb.ERC20Deposit {
	return &furypb.ERC20Deposit{
		FuryAssetId:           e.FuryAssetID,
		SourceEthereumAddress: e.SourceEthereumAddress,
		TargetPartyId:         e.TargetPartyID,
		Amount:                num.UintToString(e.Amount),
	}
}

func (e ERC20Deposit) String() string {
	return fmt.Sprintf(
		"furyAssetID(%s) targetPartyID(%s) amount(%s) sourceEthereumAddress(%s)",
		e.FuryAssetID,
		e.TargetPartyID,
		uintPointerToString(e.Amount),
		e.SourceEthereumAddress,
	)
}

func (e ERC20Deposit) GetFuryAssetID() string {
	return e.FuryAssetID
}

type ERC20EventAssetLimitsUpdated struct {
	AssetLimitsUpdated *ERC20AssetLimitsUpdated
}

func (ERC20EventAssetLimitsUpdated) isErc20EventAction() {}

func (e ERC20EventAssetLimitsUpdated) oneOfProto() interface{} {
	return e.AssetLimitsUpdated.IntoProto()
}

func (e ERC20EventAssetLimitsUpdated) String() string {
	return fmt.Sprintf(
		"assetLimitsUpdated(%s)",
		reflectPointerToString(e.AssetLimitsUpdated),
	)
}

func NewERC20EventAssetLimitsUpdated(p *furypb.ERC20Event_AssetLimitsUpdated) *ERC20EventAssetLimitsUpdated {
	return &ERC20EventAssetLimitsUpdated{
		AssetLimitsUpdated: NewERC20AssetLimitsUpdatedFromProto(p.AssetLimitsUpdated),
	}
}

func (e ERC20EventAssetLimitsUpdated) IntoProto() *furypb.ERC20Event_AssetLimitsUpdated {
	return &furypb.ERC20Event_AssetLimitsUpdated{
		AssetLimitsUpdated: e.AssetLimitsUpdated.IntoProto(),
	}
}

type ERC20AssetLimitsUpdated struct {
	FuryAssetID           string
	SourceEthereumAddress string
	LifetimeLimits        *num.Uint
	WithdrawThreshold     *num.Uint
}

func NewERC20AssetLimitsUpdatedFromProto(p *furypb.ERC20AssetLimitsUpdated) *ERC20AssetLimitsUpdated {
	lifetimeLimits, _ := num.UintFromString(p.LifetimeLimits, 10)
	withdrawThreshold, _ := num.UintFromString(p.WithdrawThreshold, 10)
	return &ERC20AssetLimitsUpdated{
		FuryAssetID:           p.FuryAssetId,
		SourceEthereumAddress: p.SourceEthereumAddress,
		LifetimeLimits:        lifetimeLimits,
		WithdrawThreshold:     withdrawThreshold,
	}
}

func (e ERC20AssetLimitsUpdated) IntoProto() *furypb.ERC20AssetLimitsUpdated {
	return &furypb.ERC20AssetLimitsUpdated{
		FuryAssetId:           e.FuryAssetID,
		SourceEthereumAddress: e.SourceEthereumAddress,
		LifetimeLimits:        num.UintToString(e.LifetimeLimits),
		WithdrawThreshold:     num.UintToString(e.WithdrawThreshold),
	}
}

func (e ERC20AssetLimitsUpdated) String() string {
	return fmt.Sprintf(
		"furyAssetID(%s) sourceEthereumAddress(%s) lifetimeLimits(%s) withdrawThreshold(%s)",
		e.FuryAssetID,
		e.SourceEthereumAddress,
		uintPointerToString(e.LifetimeLimits),
		uintPointerToString(e.WithdrawThreshold),
	)
}

func (e ERC20AssetLimitsUpdated) GetFuryAssetID() string {
	return e.FuryAssetID
}

type ERC20EventBridgeStopped struct {
	BridgeStopped bool
}

func (ERC20EventBridgeStopped) isErc20EventAction() {}

func (e ERC20EventBridgeStopped) oneOfProto() interface{} {
	return e.IntoProto()
}

func (e ERC20EventBridgeStopped) String() string {
	return fmt.Sprintf(
		"bridgeStopped(%v)",
		e.BridgeStopped,
	)
}

func NewERC20EventBridgeStopped(p *furypb.ERC20Event_BridgeStopped) *ERC20EventBridgeStopped {
	return &ERC20EventBridgeStopped{
		BridgeStopped: p.BridgeStopped,
	}
}

func (e ERC20EventBridgeStopped) IntoProto() *furypb.ERC20Event_BridgeStopped {
	return &furypb.ERC20Event_BridgeStopped{
		BridgeStopped: e.BridgeStopped,
	}
}

type ERC20EventBridgeResumed struct {
	BridgeResumed bool
}

func (ERC20EventBridgeResumed) isErc20EventAction() {}

func (e ERC20EventBridgeResumed) oneOfProto() interface{} {
	return e.IntoProto()
}

func (e ERC20EventBridgeResumed) String() string {
	return fmt.Sprintf(
		"bridgeResumed(%v)",
		e.BridgeResumed,
	)
}

func NewERC20EventBridgeResumed(p *furypb.ERC20Event_BridgeResumed) *ERC20EventBridgeResumed {
	return &ERC20EventBridgeResumed{
		BridgeResumed: p.BridgeResumed,
	}
}

func (e ERC20EventBridgeResumed) IntoProto() *furypb.ERC20Event_BridgeResumed {
	return &furypb.ERC20Event_BridgeResumed{
		BridgeResumed: e.BridgeResumed,
	}
}
