package c

import (
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/BytePackedPassword"
	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) StartLogIn(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.StartLogInParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.StartLogInResult
	r, re = c.startLogIn(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) startLogIn(p *rm.StartLogInParams) (result *rm.StartLogInResult, re *jrm1.RpcError) {
	// Access check.
	{
		re = c.mustBeNoAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}
	}

	// Check parameters.
	{
		if len(p.UserEmail) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailIsNotSet, rme.Msg_EmailIsNotSet, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())
	var user *cm.User

	// Check for existing user & log-in request with e-mail.
	{
		var exists bool
		exists, err = dbC.ExistsLogInRequestWithUserEmail(p.UserEmail)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if exists {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_LogInRequestWithUserEmailExists, rme.Msg_LogInRequestWithUserEmailExists, p.UserEmail)
		}

		user = &cm.User{Email: p.UserEmail}
		err = dbC.GetUserByEmailAbleToLogIn(user)
		if err != nil {
			return nil, c.databaseError(err)
		}
	}

	// Start logging in.
	{
		var requestId *string
		requestId, re = c.createRequestIdForLogIn()
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

		re = c.sendVerificationCodeForLogIn(p.UserEmail, *verificationCode)
		if re != nil {
			return nil, re
		}

		var pwdSalt []byte
		pwdSalt, err = bpp.GenerateRandomSalt()
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		var lir = cm.LogInRequest{
			UserEmail: p.UserEmail,

			UserId:           user.Id,
			RequestId:        *requestId,
			UserIPAB:         p.Auth.UserIPAB,
			CaptchaId:        captchaData.TaskId,
			VerificationCode: *verificationCode,
			AuthData:         pwdSalt,
		}
		err = dbC.CreateLogInRequest(lir)
		if err != nil {
			return nil, c.databaseError(err)
		}

		result = &rm.StartLogInResult{
			RequestId: *requestId,
			CaptchaId: captchaData.TaskId,
			AuthData:  pwdSalt,
		}

		return result, nil
	}
}
