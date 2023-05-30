package ethcall

import (
	"fmt"

	"github.com/elysiumstation/fury/protos/fury"
	"golang.org/x/crypto/sha3"
)

type blockish interface {
	NumberU64() uint64
	Time() uint64
}

type Trigger interface {
	Trigger(prev blockish, current blockish) bool
	ToProto() *fury.EthCallTrigger
	Hash() []byte
}

func TriggerFromProto(proto *fury.EthCallTrigger) (Trigger, error) {
	if proto == nil {
		return nil, fmt.Errorf("trigger proto is nil")
	}

	switch t := proto.Trigger.(type) {
	case *fury.EthCallTrigger_TimeTrigger:
		return TimeTriggerFromProto(t.TimeTrigger), nil
	default:
		return nil, fmt.Errorf("unknown trigger type: %T", proto.Trigger)
	}
}

type TimeTrigger struct {
	Initial uint64
	Every   uint64 // 0 = don't repeat
	Until   uint64 // 0 = forever
}

func (t TimeTrigger) Trigger(prev blockish, current blockish) bool {
	// Before initial?
	if current.Time() < t.Initial {
		return false
	}

	// Crossing initial boundary?
	if prev.Time() < t.Initial && current.Time() >= t.Initial {
		return true
	}

	// After until?
	if t.Until != 0 && current.Time() > t.Until {
		return false
	}

	if t.Every == 0 {
		return false
	}
	// Somewhere in the middle..
	prevTriggerCount := (prev.Time() - t.Initial) / t.Every
	currentTriggerCount := (current.Time() - t.Initial) / t.Every
	return currentTriggerCount > prevTriggerCount
}

func (t TimeTrigger) Hash() []byte {
	hashFunc := sha3.New256()
	ident := fmt.Sprintf("timetrigger: %v/%v/%v", t.Initial, t.Every, t.Until)
	hashFunc.Write([]byte(ident))
	return hashFunc.Sum(nil)
}

func (t TimeTrigger) ToProto() *fury.EthCallTrigger {
	var initial, every, until *uint64

	if t.Initial != 0 {
		initial = &t.Initial
	}

	if t.Every != 0 {
		every = &t.Every
	}

	if t.Until != 0 {
		until = &t.Until
	}

	return &fury.EthCallTrigger{
		Trigger: &fury.EthCallTrigger_TimeTrigger{
			TimeTrigger: &fury.EthTimeTrigger{
				Initial: initial,
				Every:   every,
				Until:   until,
			},
		},
	}
}

func TimeTriggerFromProto(proto *fury.EthTimeTrigger) TimeTrigger {
	tt := TimeTrigger{}
	if proto.Initial != nil {
		tt.Initial = *proto.Initial
	}
	if proto.Every != nil {
		tt.Every = *proto.Every
	}
	if proto.Until != nil {
		tt.Until = *proto.Until
	}
	return tt
}
