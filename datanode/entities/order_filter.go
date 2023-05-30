package entities

import "github.com/elysiumstation/fury/protos/fury"

type OrderFilter struct {
	Statuses         []fury.Order_Status
	Types            []fury.Order_Type
	TimeInForces     []fury.Order_TimeInForce
	Reference        *string
	DateRange        *DateRange
	ExcludeLiquidity bool
	LiveOnly         bool
	PartyIDs         []string
	MarketIDs        []string
}
