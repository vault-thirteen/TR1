package c

import (
	"context"
	"fmt"
	"log"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/scheduler"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/services/common/components/DatabaseComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	"gorm.io/gorm"
)

// List of component indices of the controller must be synchronised with the
// order of components used in the application's constructor.
const (
	ComponentIndex_ErrorListenerComponent = 0
	ComponentIndex_DatabaseComponent      = 1
	ComponentIndex_RpcClientComponent     = 2
	ComponentIndex_RpcServerComponent     = 3
	ComponentIndex_SchedulerComponent     = 4
)

type Controller struct {
	cfg        *cm.Configuration
	errorsChan *chan error
	service    *cm.Service
	far        ControllerFastAccessRegistry
}

func NewController() (c *Controller) {
	errorsChan := make(chan error, 1)

	return &Controller{
		errorsChan: &errorsChan,
	}
}

func (c *Controller) GetRpcFunctions() []jrm1.RpcFunction {
	return []jrm1.RpcFunction{
		// Ping.
		c.Ping,

		// Forums.
		c.AddForum,
		c.ListForums,
		c.GetForum,
		c.ChangeForumName,
		c.MoveForumUp,
		c.MoveForumDown,
		c.DeleteForum,
	}
}

func (c *Controller) GetScheduledFunctions() []sch.ScheduledFn {
	return []sch.ScheduledFn{
		//c.RemoveOutdatedSomething,
	}
}

func (c *Controller) GetErrorsChan() (errorsChan *chan error) {
	return c.errorsChan
}

func (c *Controller) LinkWithService(service interfaces.IService) (err error) {
	c.cfg = (service.GetConfiguration()).(*cm.Configuration)
	c.service = service.(*cm.Service)
	c.initFAR()

	err = c.prepareDb()
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) initFAR() {
	c.far = ControllerFastAccessRegistry{}

	c.far.systemSettings = c.cfg.GetComponent(cm.Component_System, cm.Protocol_None)
	c.far.messageSettings = c.cfg.GetComponent(cm.Component_Message, cm.Protocol_None)

	c.far.rcc = rcc.FromAny(c.service.GetComponentByIndex(ComponentIndex_RpcClientComponent))

	c.far.authServiceClient = c.far.rcc.GetClientMap()[rm.ServiceShortName_Auth]

	c.far.pageSize = c.far.systemSettings.GetParameterAsInt(ccp.PageSize)

	c.far.dbc = dc.FromAny(c.service.GetComponentByIndex(ComponentIndex_DatabaseComponent))
	c.far.db = c.far.dbc.GetGormDb()
}

func (c *Controller) prepareDb() (err error) {
	db := c.GetDb()

	if c.far.systemSettings.GetParameterAsBool(ccp.IsDatabaseInitialisationUsed) {
		classesToInit := []any{
			&cm.Forum{},
		}

		for _, cti := range classesToInit {
			err = db.AutoMigrate(cti)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) GetDb() (gormDb *gorm.DB) {
	return c.far.db
}

func (c *Controller) logError(err error) {
	if err == nil {
		return
	}

	if c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode) {
		log.Println(err)
	}
}
func (c *Controller) databaseError(err error) (re *jrm1.RpcError) {
	c.processDatabaseError(err)
	return jrm1.NewRpcErrorByUser(rme.Code_Database, rme.Msg_Database, err)
}
func (c *Controller) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if rme.IsNetworkError(err) {
		log.Println(fmt.Sprintf(rme.ErrF_DatabaseNetwork, err.Error()))
		*(c.errorsChan) <- err
	} else {
		c.logError(err)
	}

	return
}

func (c *Controller) getSelfRoles(params rm.GetSelfRolesParams) (user *cm.User, re *jrm1.RpcError) {
	var result = new(rm.GetSelfRolesResult)

	var err error
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_GetSelfRoles, params, result)
	if err != nil {
		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_RPCCall, rme.Msg_RPCCall, nil)
	}
	if re != nil {
		return nil, re
	}

	return result.User, nil
}
