package main

import (
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/services/GatewayService/controller"
	"github.com/vault-thirteen/TR1/src/services/common/components/ErrorListenerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/HttpServerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/StaticFileServerComponent"
)

func main() {
	// Order of components must be synchronised with a list of component
	// indices of the controller.
	var serviceComponents = []interfaces.IServiceComponent{
		&elc.ErrorListenerComponent{},
		&rcc.RpcClientComponent{},
		&hsc.HttpServerComponent{},
		&sfsc.StaticFileServerComponent{},
	}

	var controller interfaces.IController
	controller = c.NewController()

	app, err := cm.NewApplication(rm.ServiceName_GatewayService, serviceComponents, controller)
	mustBeNoError(err)

	err = app.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
