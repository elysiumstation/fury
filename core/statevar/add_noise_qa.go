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

//go:build qa
// +build qa

package statevar

import (
	"math/rand"

	"github.com/elysiumstation/fury/libs/num"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/protos/fury"
	furypb "github.com/elysiumstation/fury/protos/fury"
)

// AddNoise is a function used in qa build to add noise to the state variables within their tolerance to instrument consensus seeking.
func (sv *StateVariable) AddNoise(kvb []*furypb.KeyValueBundle) []*furypb.KeyValueBundle {
	for _, kvt := range kvb {
		tol, _ := num.DecimalFromString(kvt.Tolerance)
		switch v := kvt.Value.Value.(type) {
		case *fury.StateVarValue_ScalarVal:
			random := rand.Float64() * tol.InexactFloat64() / 2.0
			if sv.log.GetLevel() <= logging.DebugLevel {
				sv.log.Debug("adding random noise", logging.String("key-name", kvt.Key), logging.Float64("randomness", random))
			}
			val, _ := num.DecimalFromString(v.ScalarVal.Value)
			val = val.Add(num.DecimalFromFloat(random))
			kvt.Value.Value = &furypb.StateVarValue_ScalarVal{
				ScalarVal: &furypb.ScalarValue{
					Value: val.String(),
				},
			}

		case *fury.StateVarValue_VectorVal:
			vec := make([]num.Decimal, 0, len(v.VectorVal.Value))
			for i, entry := range v.VectorVal.Value {
				random := rand.Float64() * tol.InexactFloat64() / 2.0
				if sv.log.GetLevel() <= logging.DebugLevel {
					sv.log.Debug("adding random noise", logging.String("key-name", kvt.Key), logging.Int("index", i), logging.Float64("randomness", random))
				}
				value, _ := num.DecimalFromString(entry)
				vec = append(vec, value.Add(num.DecimalFromFloat(random)))
			}
			vecAsString := make([]string, 0, len(vec))
			for _, v := range vec {
				vecAsString = append(vecAsString, v.String())
			}
			kvt.Value.Value = &furypb.StateVarValue_VectorVal{
				VectorVal: &furypb.VectorValue{
					Value: vecAsString,
				},
			}
		}
	}
	return kvb
}
