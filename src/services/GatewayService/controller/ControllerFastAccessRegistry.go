package c

import (
	"github.com/vault-thirteen/Simple-File-Server"
	"github.com/vault-thirteen/TR1/src/models/rpc/Client"
	"github.com/vault-thirteen/TR1/src/services/common/components/HttpServerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
)

type ControllerFastAccessRegistry struct {
	systemSettings *ccse.CommonConfigurationServiceEntry

	authServiceClient    *rmc.Client
	messageServiceClient *rmc.Client

	pageSize                                  int
	messageEditTime                           int
	isDeveloperMode                           bool
	devModeHttpHeaderAccessControlAllowOrigin string
	clientIPAddressSource_CustomHeader        string
	sessionMaxDuration                        int

	rcc        *rcc.RpcClientComponent
	httpServer *hsc.HttpServerComponent
	fileServer *sfs.SimpleFileServer
}
