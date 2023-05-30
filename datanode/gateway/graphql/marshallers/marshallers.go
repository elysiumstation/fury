package marshallers

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
	"github.com/elysiumstation/fury/protos/fury"
	furypb "github.com/elysiumstation/fury/protos/fury"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
	datapb "github.com/elysiumstation/fury/protos/fury/data/v1"
	eventspb "github.com/elysiumstation/fury/protos/fury/events/v1"

	"github.com/99designs/gqlgen/graphql"
)

var ErrUnimplemented = errors.New("unmarshaller not implemented as this API is query only")

func MarshalAccountType(t fury.AccountType) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(t.String())))
	})
}

func UnmarshalAccountType(v interface{}) (fury.AccountType, error) {
	s, ok := v.(string)
	if !ok {
		return fury.AccountType_ACCOUNT_TYPE_UNSPECIFIED, fmt.Errorf("expected account type to be a string")
	}

	t, ok := fury.AccountType_value[s]
	if !ok {
		return fury.AccountType_ACCOUNT_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert AccountType from GraphQL to Proto: %v", s)
	}

	return fury.AccountType(t), nil
}

func MarshalSide(s fury.Side) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalSide(v interface{}) (fury.Side, error) {
	s, ok := v.(string)
	if !ok {
		return fury.Side_SIDE_UNSPECIFIED, fmt.Errorf("expected account type to be a string")
	}

	side, ok := fury.Side_value[s]
	if !ok {
		return fury.Side_SIDE_UNSPECIFIED, fmt.Errorf("failed to convert AccountType from GraphQL to Proto: %v", s)
	}

	return fury.Side(side), nil
}

func MarshalProposalState(s fury.Proposal_State) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProposalState(v interface{}) (fury.Proposal_State, error) {
	s, ok := v.(string)
	if !ok {
		return fury.Proposal_STATE_UNSPECIFIED, fmt.Errorf("expected proposal state to be a string")
	}

	side, ok := fury.Proposal_State_value[s]
	if !ok {
		return fury.Proposal_STATE_UNSPECIFIED, fmt.Errorf("failed to convert ProposalState from GraphQL to Proto: %v", s)
	}

	return fury.Proposal_State(side), nil
}

func MarshalTransferType(t fury.TransferType) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(t.String())))
	})
}

func UnmarshalTransferType(v interface{}) (fury.TransferType, error) {
	s, ok := v.(string)
	if !ok {
		return fury.TransferType_TRANSFER_TYPE_UNSPECIFIED, fmt.Errorf("expected transfer type to be a string")
	}

	t, ok := fury.TransferType_value[s]
	if !ok {
		return fury.TransferType_TRANSFER_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert TransferType from GraphQL to Proto: %v", s)
	}

	return fury.TransferType(t), nil
}

func MarshalTransferStatus(s eventspb.Transfer_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalTransferStatus(v interface{}) (eventspb.Transfer_Status, error) {
	return eventspb.Transfer_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalDispatchMetric(s fury.DispatchMetric) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalDispatchMetric(v interface{}) (fury.DispatchMetric, error) {
	return fury.DispatchMetric_DISPATCH_METRIC_UNSPECIFIED, ErrUnimplemented
}

func MarshalNodeStatus(s fury.NodeStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalNodeStatus(v interface{}) (fury.NodeStatus, error) {
	return fury.NodeStatus_NODE_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalAssetStatus(s fury.Asset_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalAssetStatus(v interface{}) (fury.Asset_Status, error) {
	return fury.Asset_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalNodeSignatureKind(s commandspb.NodeSignatureKind) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalNodeSignatureKind(v interface{}) (commandspb.NodeSignatureKind, error) {
	return commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_UNSPECIFIED, ErrUnimplemented
}

func MarshalOracleSpecStatus(s furypb.DataSourceSpec_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOracleSpecStatus(v interface{}) (furypb.DataSourceSpec_Status, error) {
	return furypb.DataSourceSpec_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalPropertyKeyType(s datapb.PropertyKey_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalPropertyKeyType(v interface{}) (datapb.PropertyKey_Type, error) {
	return datapb.PropertyKey_TYPE_UNSPECIFIED, ErrUnimplemented
}

func MarshalConditionOperator(s datapb.Condition_Operator) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalConditionOperator(v interface{}) (datapb.Condition_Operator, error) {
	return datapb.Condition_OPERATOR_UNSPECIFIED, ErrUnimplemented
}

func MarshalVoteValue(s fury.Vote_Value) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalVoteValue(v interface{}) (fury.Vote_Value, error) {
	return fury.Vote_VALUE_UNSPECIFIED, ErrUnimplemented
}

func MarshalAuctionTrigger(s fury.AuctionTrigger) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalAuctionTrigger(v interface{}) (fury.AuctionTrigger, error) {
	return fury.AuctionTrigger_AUCTION_TRIGGER_UNSPECIFIED, ErrUnimplemented
}

func MarshalStakeLinkingStatus(s eventspb.StakeLinking_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalStakeLinkingStatus(v interface{}) (eventspb.StakeLinking_Status, error) {
	return eventspb.StakeLinking_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalStakeLinkingType(s eventspb.StakeLinking_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalStakeLinkingType(v interface{}) (eventspb.StakeLinking_Type, error) {
	return eventspb.StakeLinking_TYPE_UNSPECIFIED, ErrUnimplemented
}

func MarshalWithdrawalStatus(s fury.Withdrawal_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalWithdrawalStatus(v interface{}) (fury.Withdrawal_Status, error) {
	return fury.Withdrawal_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalDepositStatus(s fury.Deposit_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalDepositStatus(v interface{}) (fury.Deposit_Status, error) {
	return fury.Deposit_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalOrderStatus(s fury.Order_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderStatus(v interface{}) (fury.Order_Status, error) {
	s, ok := v.(string)
	if !ok {
		return fury.Order_STATUS_UNSPECIFIED, fmt.Errorf("exoected order status to be a string")
	}

	t, ok := fury.Order_Status_value[s]
	if !ok {
		return fury.Order_STATUS_UNSPECIFIED, fmt.Errorf("failed to convert order status from GraphQL to Proto: %v", s)
	}

	return fury.Order_Status(t), nil
}

func MarshalOrderTimeInForce(s fury.Order_TimeInForce) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderTimeInForce(v interface{}) (fury.Order_TimeInForce, error) {
	s, ok := v.(string)
	if !ok {
		return fury.Order_TIME_IN_FORCE_UNSPECIFIED, fmt.Errorf("expected order time in force to be a string")
	}

	t, ok := fury.Order_TimeInForce_value[s]
	if !ok {
		return fury.Order_TIME_IN_FORCE_UNSPECIFIED, fmt.Errorf("failed to convert TimeInForce from GraphQL to Proto: %v", s)
	}

	return fury.Order_TimeInForce(t), nil
}

func MarshalPeggedReference(s fury.PeggedReference) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalPeggedReference(v interface{}) (fury.PeggedReference, error) {
	return fury.PeggedReference_PEGGED_REFERENCE_UNSPECIFIED, ErrUnimplemented
}

func MarshalProposalRejectionReason(s fury.ProposalError) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProposalRejectionReason(v interface{}) (fury.ProposalError, error) {
	return fury.ProposalError_PROPOSAL_ERROR_UNSPECIFIED, ErrUnimplemented
}

func MarshalOrderRejectionReason(s fury.OrderError) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderRejectionReason(v interface{}) (fury.OrderError, error) {
	return fury.OrderError_ORDER_ERROR_UNSPECIFIED, ErrUnimplemented
}

func MarshalOrderType(s fury.Order_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalOrderType(v interface{}) (fury.Order_Type, error) {
	s, ok := v.(string)
	if !ok {
		return fury.Order_TYPE_UNSPECIFIED, fmt.Errorf("expected order type to be a string")
	}

	t, ok := fury.Order_Type_value[s]
	if !ok {
		return fury.Order_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert OrderType from GraphQL to Proto: %v", s)
	}

	return fury.Order_Type(t), nil
}

func MarshalMarketState(s fury.Market_State) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalMarketState(v interface{}) (fury.Market_State, error) {
	return fury.Market_STATE_UNSPECIFIED, ErrUnimplemented
}

func MarshalMarketTradingMode(s fury.Market_TradingMode) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalMarketTradingMode(v interface{}) (fury.Market_TradingMode, error) {
	return fury.Market_TRADING_MODE_UNSPECIFIED, ErrUnimplemented
}

func MarshalInterval(s fury.Interval) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalInterval(v interface{}) (fury.Interval, error) {
	s, ok := v.(string)
	if !ok {
		return fury.Interval_INTERVAL_UNSPECIFIED, fmt.Errorf("expected interval in force to be a string")
	}

	t, ok := fury.Interval_value[s]
	if !ok {
		return fury.Interval_INTERVAL_UNSPECIFIED, fmt.Errorf("failed to convert Interval from GraphQL to Proto: %v", s)
	}

	return fury.Interval(t), nil
}

func MarshalProposalType(s v2.ListGovernanceDataRequest_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProposalType(v interface{}) (v2.ListGovernanceDataRequest_Type, error) {
	s, ok := v.(string)
	if !ok {
		return v2.ListGovernanceDataRequest_TYPE_UNSPECIFIED, fmt.Errorf("expected proposal type in force to be a string")
	}

	t, ok := v2.ListGovernanceDataRequest_Type_value[s]
	if !ok {
		return v2.ListGovernanceDataRequest_TYPE_UNSPECIFIED, fmt.Errorf("failed to convert proposal type from GraphQL to Proto: %v", s)
	}

	return v2.ListGovernanceDataRequest_Type(t), nil
}

func MarshalLiquidityProvisionStatus(s fury.LiquidityProvision_Status) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalLiquidityProvisionStatus(v interface{}) (fury.LiquidityProvision_Status, error) {
	return fury.LiquidityProvision_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalTradeType(s fury.Trade_Type) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalTradeType(v interface{}) (fury.Trade_Type, error) {
	return fury.Trade_TYPE_UNSPECIFIED, ErrUnimplemented
}

func MarshalValidatorStatus(s fury.ValidatorNodeStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalValidatorStatus(v interface{}) (fury.ValidatorNodeStatus, error) {
	return fury.ValidatorNodeStatus_VALIDATOR_NODE_STATUS_UNSPECIFIED, ErrUnimplemented
}

func MarshalProtocolUpgradeProposalStatus(s eventspb.ProtocolUpgradeProposalStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalProtocolUpgradeProposalStatus(v interface{}) (eventspb.ProtocolUpgradeProposalStatus, error) {
	s, ok := v.(string)
	if !ok {
		return eventspb.ProtocolUpgradeProposalStatus_PROTOCOL_UPGRADE_PROPOSAL_STATUS_UNSPECIFIED, fmt.Errorf("expected proposal type in force to be a string")
	}

	t, ok := eventspb.ProtocolUpgradeProposalStatus_value[s] // v2.ListGovernanceDataRequest_Type_value[s]
	if !ok {
		return eventspb.ProtocolUpgradeProposalStatus_PROTOCOL_UPGRADE_PROPOSAL_STATUS_UNSPECIFIED, fmt.Errorf("failed to convert proposal type from GraphQL to Proto: %v", s)
	}

	return eventspb.ProtocolUpgradeProposalStatus(t), nil
}

func MarshalPositionStatus(s fury.PositionStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(s.String())))
	})
}

func UnmarshalPositionStatus(v interface{}) (fury.PositionStatus, error) {
	s, ok := v.(string)
	if !ok {
		return fury.PositionStatus_POSITION_STATUS_UNSPECIFIED, fmt.Errorf("expected position status to be a string")
	}
	t, ok := fury.PositionStatus_value[s]
	if !ok {
		return fury.PositionStatus_POSITION_STATUS_UNSPECIFIED, fmt.Errorf("failed to convert position status to Proto: %v", s)
	}
	return fury.PositionStatus(t), nil
}
