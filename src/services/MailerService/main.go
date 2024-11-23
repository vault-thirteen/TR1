package main

import (
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/services/MailerService/controller"
	"github.com/vault-thirteen/TR1/src/services/common"
	"github.com/vault-thirteen/TR1/src/services/common/components/ErrorListenerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/MailerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RpcServerComponent"
)

func main() {
	// Order of components must be synchronised with a list of component
	// indices of the controller.
	var serviceComponents = []interfaces.IServiceComponent{
		&elc.ErrorListenerComponent{},
		&mc.MailerComponent{},
		&rsc.RpcServerComponent{},
	}

	var controller interfaces.IController
	controller = c.NewController()

	app, err := cm.NewApplication(common.ServiceName_MailerService, serviceComponents, controller)
	mustBeNoError(err)

	err = app.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
