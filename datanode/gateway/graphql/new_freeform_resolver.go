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

package gql

import (
	"context"

	"github.com/elysiumstation/fury/protos/fury"
)

type newFreeformResolver FuryResolverRoot

func (r *newFreeformResolver) DoNotUse(_ context.Context, _ *fury.NewFreeform) (*bool, error) {
	shouldUse := false
	return &shouldUse, nil
}
