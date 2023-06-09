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

package types

import (
	"fmt"

	furypb "github.com/elysiumstation/fury/protos/fury"
)

type ProposalTermsUpdateNetworkParameter struct {
	UpdateNetworkParameter *UpdateNetworkParameter
}

func (a ProposalTermsUpdateNetworkParameter) String() string {
	return fmt.Sprintf(
		"updateNetworkParameter(%s)",
		reflectPointerToString(a.UpdateNetworkParameter),
	)
}

func (a ProposalTermsUpdateNetworkParameter) IntoProto() *furypb.ProposalTerms_UpdateNetworkParameter {
	return &furypb.ProposalTerms_UpdateNetworkParameter{
		UpdateNetworkParameter: a.UpdateNetworkParameter.IntoProto(),
	}
}

func (a ProposalTermsUpdateNetworkParameter) isPTerm() {}

func (a ProposalTermsUpdateNetworkParameter) oneOfProto() interface{} {
	return a.IntoProto()
}

func (a ProposalTermsUpdateNetworkParameter) GetTermType() ProposalTermsType {
	return ProposalTermsTypeUpdateNetworkParameter
}

func (a ProposalTermsUpdateNetworkParameter) DeepClone() proposalTerm {
	if a.UpdateNetworkParameter == nil {
		return &ProposalTermsUpdateNetworkParameter{}
	}
	return &ProposalTermsUpdateNetworkParameter{
		UpdateNetworkParameter: a.UpdateNetworkParameter.DeepClone(),
	}
}

func NewUpdateNetworkParameterFromProto(
	p *furypb.ProposalTerms_UpdateNetworkParameter,
) *ProposalTermsUpdateNetworkParameter {
	var updateNP *UpdateNetworkParameter
	if p.UpdateNetworkParameter != nil {
		updateNP = &UpdateNetworkParameter{}

		if p.UpdateNetworkParameter.Changes != nil {
			updateNP.Changes = NetworkParameterFromProto(p.UpdateNetworkParameter.Changes)
		}
	}

	return &ProposalTermsUpdateNetworkParameter{
		UpdateNetworkParameter: updateNP,
	}
}

type UpdateNetworkParameter struct {
	Changes *NetworkParameter
}

func (n UpdateNetworkParameter) IntoProto() *furypb.UpdateNetworkParameter {
	return &furypb.UpdateNetworkParameter{
		Changes: n.Changes.IntoProto(),
	}
}

func (n UpdateNetworkParameter) String() string {
	return fmt.Sprintf(
		"changes(%s)",
		reflectPointerToString(n.Changes),
	)
}

func (n UpdateNetworkParameter) DeepClone() *UpdateNetworkParameter {
	if n.Changes == nil {
		return &UpdateNetworkParameter{}
	}
	return &UpdateNetworkParameter{
		Changes: n.Changes.DeepClone(),
	}
}
