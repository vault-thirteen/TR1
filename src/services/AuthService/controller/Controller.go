package c

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/net"
	"github.com/vault-thirteen/TR1/src/libraries/scheduler"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/services/common/components/CaptchaComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/DatabaseComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/JwtManagerComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RequestIdGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/services/common/components/VerificationCodeGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	"github.com/vault-thirteen/auxie/number"
	"gorm.io/gorm"
)

// List of component indices of the controller must be synchronised with the
// order of components used in the application's constructor.
const (
	ComponentIndex_ErrorListenerComponent             = 0
	ComponentIndex_DatabaseComponent                  = 1
	ComponentIndex_JwtManagerComponent                = 2
	ComponentIndex_RequestIdGeneratorComponent        = 3
	ComponentIndex_RpcClientComponent                 = 4
	ComponentIndex_VerificationCodeGeneratorComponent = 5
	ComponentIndex_RpcServerComponent                 = 6
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
		c.Ping,
		c.StartRegistration,
		c.ConfirmRegistration,
		c.StartLogIn,
		c.ConfirmLogIn,
	}
}

func (c *Controller) GetScheduledFunctions() []sch.ScheduledFn {
	return []sch.ScheduledFn{
		c.RemoveOutdatedRegistrationRequests,
		c.RemoveOutdatedLogInRequests,
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

	c.far.rcsServiceClient = c.far.rcc.GetClientMap()[rm.ServiceShortName_RCS]
	c.far.mailerServiceClient = c.far.rcc.GetClientMap()[rm.ServiceShortName_Mailer]

	c.far.dbc = dc.FromAny(c.service.GetComponentByIndex(ComponentIndex_DatabaseComponent))
	c.far.db = c.far.dbc.GetGormDb()

	c.far.ridgc = rigc.FromAny(c.service.GetComponentByIndex(ComponentIndex_RequestIdGeneratorComponent))
	c.far.ridg = c.far.ridgc.GetRidg()

	c.far.vcgc = vcgc.FromAny(c.service.GetComponentByIndex(ComponentIndex_VerificationCodeGeneratorComponent))
	c.far.vcg = c.far.vcgc.GetVcg()

	c.far.jmc = jmc.FromAny(c.service.GetComponentByIndex(ComponentIndex_JwtManagerComponent))
	c.far.jwtkm = c.far.jmc.GetKeyMaker()
}

func (c *Controller) prepareDb() (err error) {
	db := c.GetDb()

	if c.far.systemSettings.GetParameterAsBool(ccp.IsDatabaseInitialisationUsed) {
		classesToInit := []any{
			&cm.LogEvent{},
			&cm.LogInRequest{},
			&cm.Password{},
			&cm.RegistrationRequest{},
			&cm.Session{},
			&cm.User{},
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

func (c *Controller) mustBeNoAuthToken(auth *rm.Auth) (re *jrm1.RpcError) {
	re = c.mustBeAuthUserIPA(auth)
	if re != nil {
		return re
	}

	if len(auth.Token) > 0 {
		return jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
	}

	return nil
}
func (c *Controller) mustBeAuthUserIPA(auth *rm.Auth) (re *jrm1.RpcError) {
	if (auth == nil) || (len(auth.UserIPA) == 0) {
		return jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
	}

	var err error
	auth.UserIPAB, err = net.ParseIPA(auth.UserIPA)
	if err != nil {
		c.logError(err)
		return jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
	}

	return nil
}

func (c *Controller) createRequestIdForLogIn() (rid *string, re *jrm1.RpcError) {
	var err error
	rid, err = c.far.ridg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_RequestIdGenerator, rme.Msg_RequestIdGenerator, nil)
	}

	return rid, nil
}
func (c *Controller) createCaptcha() (result *rm.CreateCaptchaResult, re *jrm1.RpcError) {
	var params = rm.CreateCaptchaParams{}
	result = new(rm.CreateCaptchaResult)

	var err error
	re, err = c.far.rcsServiceClient.MakeRequest(context.Background(), rm.Func_CreateCaptcha, params, result)
	if err != nil {
		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_RPCCall, rme.Msg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaError, fmt.Sprintf(rme.MsgF_CaptchaError, re.AsError().Error()), nil)
	}

	if result.IsImageDataReturned {
		// We do not return an image in a JSON message.
		err = errors.New(cc.Err_UnexpectedResponse)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaError, fmt.Sprintf(rme.MsgF_CaptchaError, err.Error()), nil)
	}

	return result, nil
}
func (c *Controller) createVerificationCode() (vc *string, re *jrm1.RpcError) {
	var err error
	var s *string
	s, err = c.far.vcg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_VerificationCodeGenerator, rme.Msg_VerificationCodeGenerator, nil)
	}

	return s, nil
}
func (c *Controller) sendVerificationCodeForReg(email string, vCode string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplateForRegVCode), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplateForRegVCode), c.far.systemSettings.GetParameterAsString(ccp.SiteName), vCode)
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendVerificationCodeForLogIn(email string, vCode string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplateForRegVCode), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplateForLogIn), vCode)
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendEmailMessage(params rm.SendEmailMessageParams) (re *jrm1.RpcError) {
	var result = new(rm.SendEmailMessageResult)

	var err error
	re, err = c.far.mailerServiceClient.MakeRequest(context.Background(), rm.Func_SendEmailMessage, params, result)
	if err != nil {
		c.logError(err)
		return jrm1.NewRpcErrorByUser(rme.Code_RPCCall, rme.Msg_RPCCall, nil)
	}
	if re != nil {
		return jrm1.NewRpcErrorByUser(rme.Code_MailerError, fmt.Sprintf(rme.MsgF_MailerError, re.AsError().Error()), err)
	}

	return nil
}
func (c *Controller) checkCaptcha(captchaId string, answer string) (isCorrect bool, re *jrm1.RpcError) {
	n, err := number.ParseUint(answer)
	if err != nil {
		return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaError, fmt.Sprintf(rme.MsgF_CaptchaError, err.Error()), answer)
	}

	var params = rm.CheckCaptchaParams{
		TaskId: captchaId,
		Value:  n,
	}

	// Fool check.
	{
		if len(captchaId) == 0 {
			return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaTaskIdIsNotSet, rme.Msg_CaptchaTaskIdIsNotSet, nil)
		}
		if len(answer) == 0 {
			return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsNotSet, rme.Msg_CaptchaAnswerIsNotSet, nil)
		}

		params.Value, err = number.ParseUint(answer)
		if err != nil {
			return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsNotSet, rme.Msg_CaptchaAnswerIsNotSet, nil)
		}
	}

	var result = new(rm.CheckCaptchaResult)
	re, err = c.far.rcsServiceClient.MakeRequest(context.Background(), rm.Func_CheckCaptcha, params, result)
	if err != nil {
		c.logError(err)
		return false, jrm1.NewRpcErrorByUser(rme.Code_RPCCall, rme.Msg_RPCCall, nil)
	}
	if re != nil {
		return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaCheckError, fmt.Sprintf(rme.Msg_CaptchaCheckError, re.AsError().Error()), nil)
	}

	return result.IsSuccess, nil
}