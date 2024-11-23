package rm

const (
	RpcDurationFieldName  = "dur"
	RpcRequestIdFieldName = "rid"
)

const (
	ServiceNextPingAttemptDelaySec     = 5
	ServicePingAttemptsDurationMinutes = 15
)

const (
	UrlSchemeHttp  = "http"
	UrlSchemeHttps = "https"
)

const (
	Func_Ping = "Ping"
)

const (
	ServiceShortName_Auth    = "auth"
	ServiceShortName_Captcha = "captcha" // Captcha images (Proxy).
	ServiceShortName_Gateway = "gateway"
	ServiceShortName_Mailer  = "mailer"
	ServiceShortName_Message = "message"
	ServiceShortName_RCS     = "rcs" // Captcha questions (RPC).
)

// Function names.

// Captcha service.
const (
	Func_CreateCaptcha = "CreateCaptcha"
	Func_CheckCaptcha  = "CheckCaptcha"
)

// Mailer service.
const (
	Func_SendEmailMessage = "SendEmailMessage"
)
