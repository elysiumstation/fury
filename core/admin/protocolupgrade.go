package admin

import (
	"net/http"

	"github.com/elysiumstation/fury/core/types"
)

type ReadyForUpgradeArgs struct{}

type ProtocolUpgradeAdminService struct {
	protocolUpgradeService ProtocolUpgradeService
}

func NewProtocolUpgradeService(protocolUpgradeService ProtocolUpgradeService) *ProtocolUpgradeAdminService {
	return &ProtocolUpgradeAdminService{
		protocolUpgradeService: protocolUpgradeService,
	}
}

func (p *ProtocolUpgradeAdminService) UpgradeStatus(r *http.Request, args *ReadyForUpgradeArgs, reply *types.UpgradeStatus) error {
	*reply = p.protocolUpgradeService.GetUpgradeStatus()
	return nil
}
