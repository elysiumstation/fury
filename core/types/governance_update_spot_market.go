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

	"github.com/elysiumstation/fury/libs/num"
	furypb "github.com/elysiumstation/fury/protos/fury"
)

type ProposalTermsUpdateSpotMarket struct {
	UpdateSpotMarket *UpdateSpotMarket
}

func (a ProposalTermsUpdateSpotMarket) String() string {
	return fmt.Sprintf(
		"updateSpotMarket(%s)",
		reflectPointerToString(a.UpdateSpotMarket),
	)
}

func (a ProposalTermsUpdateSpotMarket) IntoProto() *furypb.ProposalTerms_UpdateSpotMarket {
	return &furypb.ProposalTerms_UpdateSpotMarket{
		UpdateSpotMarket: a.UpdateSpotMarket.IntoProto(),
	}
}

func (a ProposalTermsUpdateSpotMarket) isPTerm() {}

func (a ProposalTermsUpdateSpotMarket) oneOfProto() interface{} {
	return a.IntoProto()
}

func (a ProposalTermsUpdateSpotMarket) GetTermType() ProposalTermsType {
	return ProposalTermsTypeUpdateSpotMarket
}

func (a ProposalTermsUpdateSpotMarket) DeepClone() proposalTerm {
	if a.UpdateSpotMarket == nil {
		return &ProposalTermsUpdateSpotMarket{}
	}
	return &ProposalTermsUpdateSpotMarket{
		UpdateSpotMarket: a.UpdateSpotMarket.DeepClone(),
	}
}

func UpdateSpotMarketFromProto(p *furypb.ProposalTerms_UpdateSpotMarket) (*ProposalTermsUpdateSpotMarket, error) {
	var updateSpotMarket *UpdateSpotMarket
	if p.UpdateSpotMarket != nil {
		updateSpotMarket = &UpdateSpotMarket{}
		updateSpotMarket.MarketID = p.UpdateSpotMarket.MarketId
		if p.UpdateSpotMarket.Changes != nil {
			var err error
			updateSpotMarket.Changes, err = UpdateSpotMarketConfigurationFromProto(p.UpdateSpotMarket.Changes)
			if err != nil {
				return nil, err
			}
		}
	}
	return &ProposalTermsUpdateSpotMarket{
		UpdateSpotMarket: updateSpotMarket,
	}, nil
}

type UpdateSpotMarket struct {
	MarketID string
	Changes  *UpdateSpotMarketConfiguration
}

func (n UpdateSpotMarket) String() string {
	return fmt.Sprintf(
		"marketID(%s) changes(%s)",
		n.MarketID,
		reflectPointerToString(n.Changes),
	)
}

func (n UpdateSpotMarket) IntoProto() *furypb.UpdateSpotMarket {
	var changes *furypb.UpdateSpotMarketConfiguration
	if n.Changes != nil {
		changes = n.Changes.IntoProto()
	}
	return &furypb.UpdateSpotMarket{
		MarketId: n.MarketID,
		Changes:  changes,
	}
}

func (n UpdateSpotMarket) DeepClone() *UpdateSpotMarket {
	cpy := UpdateSpotMarket{
		MarketID: n.MarketID,
	}
	if n.Changes != nil {
		cpy.Changes = n.Changes.DeepClone()
	}
	return &cpy
}

type UpdateSpotMarketConfiguration struct {
	Instrument                *UpdateInstrumentConfiguration
	Metadata                  []string
	PriceMonitoringParameters *PriceMonitoringParameters
	TargetStakeParameters     *TargetStakeParameters
	RiskParameters            updateRiskParams
}

func (n UpdateSpotMarketConfiguration) String() string {
	return fmt.Sprintf(
		"instrument(%s) metadata(%v) priceMonitoring(%s) targetStakeParameters(%s) risk(%s)",
		reflectPointerToString(n.Instrument),
		MetadataList(n.Metadata).String(),
		reflectPointerToString(n.PriceMonitoringParameters),
		reflectPointerToString(n.TargetStakeParameters),
		reflectPointerToString(n.RiskParameters),
	)
}

func (n UpdateSpotMarketConfiguration) DeepClone() *UpdateSpotMarketConfiguration {
	cpy := &UpdateSpotMarketConfiguration{
		Metadata: make([]string, len(n.Metadata)),
	}
	cpy.Metadata = append(cpy.Metadata, n.Metadata...)
	if n.Instrument != nil {
		cpy.Instrument = n.Instrument.DeepClone()
	}
	if n.PriceMonitoringParameters != nil {
		cpy.PriceMonitoringParameters = n.PriceMonitoringParameters.DeepClone()
	}
	if n.TargetStakeParameters != nil {
		cpy.TargetStakeParameters = n.TargetStakeParameters.DeepClone()
	}
	if n.RiskParameters != nil {
		cpy.RiskParameters = n.RiskParameters.DeepClone()
	}
	return cpy
}

func (n UpdateSpotMarketConfiguration) IntoProto() *furypb.UpdateSpotMarketConfiguration {
	riskParams := n.RiskParameters.updateRiskParamsIntoProto()
	md := make([]string, 0, len(n.Metadata))
	md = append(md, n.Metadata...)

	var priceMonitoring *furypb.PriceMonitoringParameters
	if n.PriceMonitoringParameters != nil {
		priceMonitoring = n.PriceMonitoringParameters.IntoProto()
	}
	targetStakeParameters := n.TargetStakeParameters.IntoProto()

	r := &furypb.UpdateSpotMarketConfiguration{
		Metadata:                  md,
		PriceMonitoringParameters: priceMonitoring,
		TargetStakeParameters:     targetStakeParameters,
	}
	switch rp := riskParams.(type) {
	case *furypb.UpdateSpotMarketConfiguration_Simple:
		r.RiskParameters = rp
	case *furypb.UpdateSpotMarketConfiguration_LogNormal:
		r.RiskParameters = rp
	}
	return r
}

func UpdateSpotMarketConfigurationFromProto(p *furypb.UpdateSpotMarketConfiguration) (*UpdateSpotMarketConfiguration, error) {
	md := make([]string, 0, len(p.Metadata))
	md = append(md, p.Metadata...)
	var priceMonitoring *PriceMonitoringParameters
	if p.PriceMonitoringParameters != nil {
		priceMonitoring = PriceMonitoringParametersFromProto(p.PriceMonitoringParameters)
	}
	targetStakeParameters := TargetStakeParametersFromProto(p.TargetStakeParameters)

	r := &UpdateSpotMarketConfiguration{
		Metadata:                  md,
		PriceMonitoringParameters: priceMonitoring,
		TargetStakeParameters:     targetStakeParameters,
	}
	if p.RiskParameters != nil {
		switch rp := p.RiskParameters.(type) {
		case *furypb.UpdateSpotMarketConfiguration_Simple:
			r.RiskParameters = UpdateSpotMarketConfigurationSimpleFromProto(rp)
		case *furypb.UpdateSpotMarketConfiguration_LogNormal:
			r.RiskParameters = UpdateSpotMarketConfigurationLogNormalFromProto(rp)
		}
	}
	return r, nil
}

type UpdateSpotMarketConfigurationSimple struct {
	Simple *SimpleModelParams
}

func (n UpdateSpotMarketConfigurationSimple) String() string {
	return fmt.Sprintf(
		"simple(%s)",
		reflectPointerToString(n.Simple),
	)
}

func (n UpdateSpotMarketConfigurationSimple) updateRiskParamsIntoProto() interface{} {
	return n.IntoProto()
}

func (n UpdateSpotMarketConfigurationSimple) DeepClone() updateRiskParams {
	if n.Simple == nil {
		return &UpdateSpotMarketConfigurationSimple{}
	}
	return &UpdateSpotMarketConfigurationSimple{
		Simple: n.Simple.DeepClone(),
	}
}

func (n UpdateSpotMarketConfigurationSimple) IntoProto() *furypb.UpdateSpotMarketConfiguration_Simple {
	return &furypb.UpdateSpotMarketConfiguration_Simple{
		Simple: n.Simple.IntoProto(),
	}
}

func UpdateSpotMarketConfigurationSimpleFromProto(p *furypb.UpdateSpotMarketConfiguration_Simple) *UpdateSpotMarketConfigurationSimple {
	return &UpdateSpotMarketConfigurationSimple{
		Simple: SimpleModelParamsFromProto(p.Simple),
	}
}

type UpdateSpotMarketConfigurationLogNormal struct {
	LogNormal *LogNormalRiskModel
}

func (n UpdateSpotMarketConfigurationLogNormal) String() string {
	return fmt.Sprintf(
		"logNormal(%s)",
		reflectPointerToString(n.LogNormal),
	)
}

func (n UpdateSpotMarketConfigurationLogNormal) updateRiskParamsIntoProto() interface{} {
	return n.IntoProto()
}

func (n UpdateSpotMarketConfigurationLogNormal) DeepClone() updateRiskParams {
	if n.LogNormal == nil {
		return &UpdateSpotMarketConfigurationLogNormal{}
	}
	return &UpdateSpotMarketConfigurationLogNormal{
		LogNormal: n.LogNormal.DeepClone(),
	}
}

func (n UpdateSpotMarketConfigurationLogNormal) IntoProto() *furypb.UpdateSpotMarketConfiguration_LogNormal {
	return &furypb.UpdateSpotMarketConfiguration_LogNormal{
		LogNormal: n.LogNormal.IntoProto(),
	}
}

func UpdateSpotMarketConfigurationLogNormalFromProto(p *furypb.UpdateSpotMarketConfiguration_LogNormal) *UpdateSpotMarketConfigurationLogNormal {
	return &UpdateSpotMarketConfigurationLogNormal{
		LogNormal: &LogNormalRiskModel{
			RiskAversionParameter: num.DecimalFromFloat(p.LogNormal.RiskAversionParameter),
			Tau:                   num.DecimalFromFloat(p.LogNormal.Tau),
			Params:                LogNormalParamsFromProto(p.LogNormal.Params),
		},
	}
}
