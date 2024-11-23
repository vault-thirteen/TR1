package main

import (
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/services/AuthService/controller"
	"github.com/vault-thirteen/TR1/src/services/common"
	"github.com/vault-thirteen/TR1/src/services/common/components/DatabaseComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/ErrorListenerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/JwtManagerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RequestIdGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RpcServerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/SchedulerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/VerificationCodeGeneratorComponent"
)

func main() {
	// Order of components must be synchronised with a list of component
	// indices of the controller.
	var serviceComponents = []interfaces.IServiceComponent{
		&elc.ErrorListenerComponent{},
		&dc.DatabaseComponent{},
		&jmc.JwtManagerComponent{},
		&rigc.RequestIdGeneratorComponent{},
		&rcc.RpcClientComponent{},
		&vcgc.VerificationCodeGeneratorComponent{},
		&rsc.RpcServerComponent{},
		&shc.SchedulerComponent{},
	}

	var controller interfaces.IController
	controller = c.NewController()

	app, err := cm.NewApplication(common.ServiceName_AuthService, serviceComponents, controller)
	mustBeNoError(err)

	err = app.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
