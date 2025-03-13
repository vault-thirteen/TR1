package c

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	hm "github.com/vault-thirteen/TR1/src/models/http"
	rm "github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

// Static files.
const (
	StaticFile_IndexHtml = "index.html"
)

// URL paths.
const (
	UrlPath_Api = `/api`
)

const (
	ErrFUnknownRpcErrorCode = "unknown RPC error code: %v"
)

func (c *Controller) initGatewayRouter() {
	c.far.httpServer.SetHttpRouter(http.Handler(http.HandlerFunc(c.gatewayRouter)))
}

func (c *Controller) processInternalServerError(rw http.ResponseWriter, err error) {
	c.logError(err)
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusInternalServerError)
}
func (c *Controller) processRpcError(re *jrm1.RpcError, rw http.ResponseWriter) {
	httpStatusCode, ok := c.httpStatusCodesByRpcErrorCode[re.Code.Int()]
	if !ok {
		err := fmt.Errorf(ErrFUnknownRpcErrorCode, re.Code.Int())
		c.processInternalServerError(rw, err)
		return
	}

	switch httpStatusCode {
	case http.StatusInternalServerError:
		err := re.AsError()
		c.processInternalServerError(rw, err)
		return
	}

	c.respondWithPlainText(rw, re.AsError().Error(), httpStatusCode)
	return
}
func (c *Controller) respondWithPlainText(rw http.ResponseWriter, text string, httpStatusCode int) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.Header().Set(header.HttpHeaderContentType, hm.ContentType_PlainText)
	rw.WriteHeader(httpStatusCode)

	_, err := rw.Write([]byte(text))
	if err != nil {
		c.logError(err)
		return
	}
}
func (c *Controller) respondWithJsonObject(rw http.ResponseWriter, obj any) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.Header().Set(header.HttpHeaderContentType, hm.ContentType_Json)

	err := json.NewEncoder(rw).Encode(obj)
	if err != nil {
		c.logError(err)
		return
	}
}
func (c *Controller) respondBadRequest(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusBadRequest)
}
func (c *Controller) respondForbidden(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusForbidden)
}
func (c *Controller) respondMethodNotAllowed(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
func (c *Controller) respondNotAcceptable(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusNotAcceptable)
}
func (c *Controller) respondNotFound(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusNotFound)
}

func (c *Controller) setTokenCookie(rw http.ResponseWriter, token string) {
	var cookie = &http.Cookie{
		Name:   rm.CookieName_Token,
		Value:  token,
		MaxAge: c.far.sessionMaxDuration,

		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}

	hm.SetCookie(rw, cookie)
}
func (c *Controller) clearTokenCookie(rw http.ResponseWriter) {
	var cookie = &http.Cookie{
		Name: rm.CookieName_Token,
		//Value
		//MaxAge

		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}

	hm.UnsetCookie(rw, cookie)
}

func (c *Controller) getClientIPAddress(req *http.Request) (cipa string, err error) {
	var host string

	if len(c.far.clientIPAddressSource_CustomHeader) == 0 {
		host, _, err = net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return "", err
		}

		return host, nil
	}

	host, err = hh.GetSingleHttpHeader(req, c.far.clientIPAddressSource_CustomHeader)
	if err != nil {
		return "", err
	}

	return host, nil
}

func (c *Controller) gatewayRouter(rw http.ResponseWriter, req *http.Request) {
	clientIPA, err := c.getClientIPAddress(req)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}

	switch req.URL.Path {
	case UrlPath_Api:
		c.handleApiRequest(rw, req, clientIPA)
		return
	}

	//TODO: Add other handlers.

	var file []byte
	file, err = c.far.fileServer.GetFile(StaticFile_IndexHtml)
	if err != nil {
		c.logError(err)
		return
	}

	//TODO:Add MIME Type.
	_, err = rw.Write(file)
	if err != nil {
		c.logError(err)
		return
	}
}
