package c

//import (
//	"fmt"
//	"time"
//
//	"github.com/vault-thirteen/TR1/src/models/common"
//	"github.com/vault-thirteen/TR1/src/models/dbc"
//	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
//)

// Functions for scheduler component.

// TODO
//func (c *Controller) RemoveOutdatedSomething() (err error) {
//	sessionTtl := c.far.systemSettings.GetParameterAsInt(ccp.SessionMaxDuration)
//	edgeTime := time.Now().Add(-time.Duration(sessionTtl) * time.Second)
//	dbC := dbc.NewDbController(c.GetDb())
//	var isDebugMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode)
//
//	var s any
//	for {
//		s, err = c.getNextOutdatedSomething(edgeTime)
//		if err != nil {
//			return err
//		}
//
//		if s == nil {
//			break
//		}
//
//		if isDebugMode {
//			fmt.Println("removing outdated something. Id:", s)
//		}
//
//		// Delete session.
//		err = dbC.DeleteOldSession(s)
//		if err != nil {
//			return err
//		}
//
//		// Journaling.
//		logEvent := cm.NewLogEvent(cm.LogEvent_Type_LogOutByTimeout, -1, nil, nil)
//
//		err = dbC.CreateLogEvent(logEvent)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//func (c *Controller) getNextOutdatedSomething(edgeTime time.Time) (s any, err error) {
//	dbC := dbc.NewDbController(c.GetDb())
//
//	var ss []cm.Session
//	ss, err = dbC.GetFirstOutdatedSession(edgeTime)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(ss) != 1 {
//		return nil, nil
//	}
//
//	return &ss[0], nil
//}
