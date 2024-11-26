package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

func (c *Controller) StartRegistration(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.StartRegistrationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.StartRegistrationResult
	r, re = c.startRegistration(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) startRegistration(p *rm.StartRegistrationParams) (result *rm.StartRegistrationResult, re *jrm1.RpcError) {
	// Access check.
	{
		re = c.mustBeNoAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}
	}

	// Check parameters.
	{
		if len(p.UserName) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_NameIsNotSet, rme.Msg_NameIsNotSet, nil)
		}
		if len(p.UserEmail) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailIsNotSet, rme.Msg_EmailIsNotSet, nil)
		}
		if len(p.UserPassword) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsNotSet, rme.Msg_PasswordIsNotSet, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())

	// Check for existing user with name & e-mail.
	{
		var isFree bool
		isFree, err = dbC.IsUserNameFree(p.UserName)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if !isFree {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserNameIsUsed, rme.Msg_UserNameIsUsed, p.UserName)
		}

		isFree, err = dbC.IsUserEmailFree(p.UserEmail)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if !isFree {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserEmailIsUsed, rme.Msg_UserEmailIsUsed, p.UserEmail)
		}
	}

	// Check for existing registration request with name & e-mail.
	{
		var exists bool
		exists, err = dbC.ExistsRegistrationRequestWithUserName(p.UserName)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if exists {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_RegistrationRequestWithUserNameExists, rme.Msg_RegistrationRequestWithUserNameExists, p.UserName)
		}

		exists, err = dbC.ExistsRegistrationRequestWithUserEmail(p.UserEmail)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if exists {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_RegistrationRequestWithUserEmailExists, rme.Msg_RegistrationRequestWithUserEmailExists, p.UserEmail)
		}
	}

	// Other checks.
	{
		if !cm.IsUserEmailValid(p.UserEmail) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserEmailIsInvalid, rme.Msg_UserEmailIsInvalid, p.UserEmail)
		}

		if len([]byte(p.UserName)) > c.far.systemSettings.GetParameterAsInt(ccp.UserNameMaxLenInBytes) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserNameIsTooLong, rme.Msg_UserNameIsTooLong, p.UserName)
		}

		if len([]byte(p.UserPassword)) > c.far.systemSettings.GetParameterAsInt(ccp.UserPasswordMaxLenInBytes) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserPasswordIsTooLong, rme.Msg_UserPasswordIsTooLong, nil)
		}

		ok := cm.IsUserPasswordAllowed(p.UserPassword)
		if !ok {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserPasswordIsNotAllowed, rme.Msg_UserPasswordIsNotAllowed, nil)
		}
	}

	// Start user registration.
	{
		var requestId *string
		requestId, re = c.createRequestId()
		if re != nil {
			return nil, re
		}

		var captchaData *rm.CreateCaptchaResult
		captchaData, re = c.createCaptcha()
		if re != nil {
			return nil, re
		}

		var verificationCode *string
		verificationCode, re = c.createVerificationCode()
		if re != nil {
			return nil, re
		}

		re = c.sendVerificationCodeForReg(p.UserEmail, *verificationCode)
		if re != nil {
			return nil, re
		}

		var rr = cm.RegistrationRequest{
			UserName:     p.UserName,
			UserEmail:    p.UserEmail,
			UserPassword: p.UserPassword,

			RequestId:        *requestId,
			UserIPAB:         p.Auth.UserIPAB,
			CaptchaId:        captchaData.TaskId,
			VerificationCode: *verificationCode,
		}
		err = dbC.CreateRegistrationRequest(rr)
		if err != nil {
			return nil, c.databaseError(err)
		}

		result = &rm.StartRegistrationResult{
			RequestId: *requestId,
			CaptchaId: captchaData.TaskId,
		}

		return result, nil
	}
}
