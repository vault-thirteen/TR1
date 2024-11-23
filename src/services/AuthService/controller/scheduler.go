package c

import (
	"fmt"
	"time"

	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

// Functions for scheduler component.

func (c *Controller) RemoveOutdatedRegistrationRequests() (err error) {
	rrTtl := c.far.systemSettings.GetParameterAsInt(ccp.RegistrationRequestTtl)
	edgeTime := time.Now().Add(-time.Duration(rrTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var rr *cm.RegistrationRequest
	for {
		rr, err = c.getNextOutdatedRegistrationRequest(edgeTime)
		if err != nil {
			return err
		}

		if rr == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated registration request. RequestId:", rr.RequestId)
		}
		err = dbC.DeleteRegistrationRequestNRFA(rr)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedRegistrationRequest(edgeTime time.Time) (rr *cm.RegistrationRequest, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var rrs []cm.RegistrationRequest
	rrs, err = dbC.GetFirstOutdatedRegistrationRequest(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(rrs) != 1 {
		return nil, nil
	}

	return &rrs[0], nil
}

func (c *Controller) RemoveOutdatedLogInRequests() (err error) {
	lirTtl := c.far.systemSettings.GetParameterAsInt(ccp.LogInRequestTtl)
	edgeTime := time.Now().Add(-time.Duration(lirTtl) * time.Second)
	dbC := dbc.NewDbController(c.GetDb())
	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)

	var lir *cm.LogInRequest
	for {
		lir, err = c.getNextOutdatedLogInRequest(edgeTime)
		if err != nil {
			return err
		}

		if lir == nil {
			break
		}

		if isDebugMode {
			fmt.Println("removing outdated log-in request. RequestId:", lir.RequestId)
		}
		err = dbC.DeleteOldLogInRequest(lir)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Controller) getNextOutdatedLogInRequest(edgeTime time.Time) (lir *cm.LogInRequest, err error) {
	dbC := dbc.NewDbController(c.GetDb())

	var lirs []cm.LogInRequest
	lirs, err = dbC.GetFirstOutdatedLogInRequest(edgeTime)
	if err != nil {
		return nil, err
	}

	if len(lirs) != 1 {
		return nil, nil
	}

	return &lirs[0], nil
}
