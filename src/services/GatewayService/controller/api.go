package c

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	rm "github.com/vault-thirteen/TR1/src/models/rpc"
	rme "github.com/vault-thirteen/TR1/src/models/rpc/error"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

// TODO
func (c *Controller) initAPI() {
	c.httpStatusCodesByRpcErrorCode = rme.GetMapOfHttpStatusCodesByRpcErrorCodes()

	c.apiHandlers = map[string]rm.RequestHandler{

		// AuthService.
		rm.ApiFunctionName_StartRegistration: c.StartRegistration,
	}
}

func (c *Controller) handleApiRequest(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	if req.Method != http.MethodPost {
		c.respondMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := hh.CheckBrowserSupportForJson(req)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}
	if !ok {
		c.respondNotAcceptable(rw)
		return
	}

	var reqBody []byte
	reqBody, err = io.ReadAll(req.Body)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}

	// Get the requested function.
	var arwoa rm.RequestWithOnlyAction
	err = json.Unmarshal(reqBody, &arwoa)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	if (arwoa.Action == nil) ||
		(arwoa.Parameters == nil) {
		c.respondBadRequest(rw)
		return
	}

	var handler rm.RequestHandler
	handler, ok = c.apiHandlers[*arwoa.Action]
	if !ok {
		c.respondNotFound(rw)
		return
	}

	var token *string
	token, err = rm.GetToken(req)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	var ar = &rm.Request{
		Action:     arwoa.Action,
		Parameters: arwoa.Parameters,
		Authorisation: &rm.Auth{
			UserIPA: clientIPA,
		},
	}

	if token != nil {
		ar.Authorisation.Token = *token
	}

	handler(ar, req, rw)
	return
}

// AuthService.

func (c *Controller) StartRegistration(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.StartRegistrationParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.StartRegistrationResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_StartRegistration, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}

//TODO: Add other handlers.

// MessageService.

//TODO: Add other handlers.
