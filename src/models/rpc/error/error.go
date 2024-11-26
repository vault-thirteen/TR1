package rme

import (
	"errors"
	"net"
	"net/http"
)

// RPC error codes.
const (
	Code_FeatureIsNotImplemented                = 1
	Code_Database                               = 2
	Code_Authorisation                          = 3
	Code_Permission                             = 4
	Code_RPCCall                                = 5
	Code_MailerError                            = 6
	Code_EmailRecipientIsNotSet                 = 7
	Code_EmailSubjectIsNotSet                   = 8
	Code_EmailMessageIsNotSet                   = 9
	Code_CaptchaError                           = 10
	Code_CaptchaCreationError                   = 11
	Code_CaptchaTaskIdIsNotSet                  = 12
	Code_CaptchaAnswerIsNotSet                  = 13
	Code_CaptchaCheckError                      = 14
	Code_CaptchaAnswerIsWrong                   = 15
	Code_CaptchaIsNotFound                      = 16
	Code_NameIsNotSet                           = 17
	Code_EmailIsNotSet                          = 18
	Code_PasswordIsNotSet                       = 19
	Code_UserNameIsUsed                         = 20
	Code_UserEmailIsUsed                        = 21
	Code_RegistrationRequestWithUserNameExists  = 22
	Code_RegistrationRequestWithUserEmailExists = 23
	Code_UserEmailIsInvalid                     = 24
	Code_UserNameIsTooLong                      = 25
	Code_UserPasswordIsTooLong                  = 26
	Code_UserPasswordIsNotAllowed               = 27
	Code_RequestIdGenerator                     = 28
	Code_VerificationCodeGenerator              = 29
	Code_RequestIdIsNotSet                      = 30
	Code_VerificationCodeIsNotSet               = 31
	Code_BPP                                    = 32
	Code_LogInRequestWithUserEmailExists        = 33
	Code_AuthDataIsNotSet                       = 34
	Code_PasswordIsWrong                        = 35
	Code_JWT                                    = 36
	Code_VerificationCodeIsWrong                = 37
	Code_SessionIsNotFound                      = 38
	Code_NotAuthorised                          = 39
	Code_NotSure                                = 40
	Code_TokenIsExpired                         = 41
	Code_EmailChangeRequestWithNewEmailExists   = 42
	Code_PasswordChangeRequestWithUserIdExists  = 43
	Code_PageIsNotSet                           = 44
)

// RPC error messages.
const (
	MsgMsg_FeatureIsNotImplemented             = "feature is not implemented"
	Msg_Database                               = "database error"
	Msg_Authorisation                          = "authorisation error"
	Msg_Permission                             = "permission error"
	Msg_RPCCall                                = "RPC call error"
	MsgF_MailerError                           = "mailer error: %s"
	Msg_EmailRecipientIsNotSet                 = "e-mail recipient is not set"
	Msg_EmailSubjectIsNotSet                   = "e-mail subject is not set"
	Msg_EmailMessageIsNotSet                   = "e-mail message is not set"
	MsgF_CaptchaError                          = "captcha error: %s"
	MsgF_CaptchaCreationError                  = "captcha creation error: %s"
	Msg_CaptchaTaskIdIsNotSet                  = "captcha task ID is not set"
	Msg_CaptchaAnswerIsNotSet                  = "captcha answer is not set"
	Msg_CaptchaCheckError                      = "captcha check error: %s"
	Msg_CaptchaAnswerIsWrong                   = "captcha answer is wrong"
	Msg_CaptchaIsNotFound                      = "captcha is not found"
	Msg_NameIsNotSet                           = "name is not set"
	Msg_EmailIsNotSet                          = "e-mail is not set"
	Msg_PasswordIsNotSet                       = "password is not set"
	Msg_UserNameIsUsed                         = "user name is used"
	Msg_UserEmailIsUsed                        = "user e-mail is used"
	Msg_RegistrationRequestWithUserNameExists  = "registration request with username exists"
	Msg_RegistrationRequestWithUserEmailExists = "registration request with user e-mail exists"
	Msg_UserEmailIsInvalid                     = "user e-mail is invalid"
	Msg_UserNameIsTooLong                      = "user name is too long"
	Msg_UserPasswordIsTooLong                  = "user password is too long"
	Msg_UserPasswordIsNotAllowed               = "user password is not allowed"
	Msg_RequestIdGenerator                     = "request ID generator error"
	Msg_VerificationCodeGenerator              = "verification code generator error"
	Msg_RequestIdIsNotSet                      = "request ID is not set"
	Msg_VerificationCodeIsNotSet               = "verification code is not set"
	MsgF_BPP                                   = "BPP error: %s"
	Msg_LogInRequestWithUserEmailExists        = "log-in request with user e-mail exists"
	Msg_AuthDataIsNotSet                       = "auth data is not set"
	Msg_PasswordIsWrong                        = "password is wrong"
	MsgF_JWT                                   = "JWT error: %s"
	Msg_VerificationCodeIsWrong                = "verification code is wrong"
	Msg_SessionIsNotFound                      = "session is not found"
	Msg_NotAuthorised                          = "not authorised"
	Msg_NotSure                                = "not sure"
	Msg_TokenIsExpired                         = "token is expired"
	Msg_EmailChangeRequestWithNewEmailExists   = "e-mail change request with new e-mail exists"
	Msg_PasswordChangeRequestWithUserIdExists  = "password change request with user ID exists"
	Msg_PageIsNotSet                           = "page is not set"
)

func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		Code_FeatureIsNotImplemented:                http.StatusInternalServerError,
		Code_Database:                               http.StatusInternalServerError,
		Code_Authorisation:                          http.StatusUnauthorized,
		Code_Permission:                             http.StatusForbidden,
		Code_RPCCall:                                http.StatusInternalServerError,
		Code_MailerError:                            http.StatusInternalServerError,
		Code_EmailRecipientIsNotSet:                 http.StatusBadRequest,
		Code_EmailSubjectIsNotSet:                   http.StatusBadRequest,
		Code_EmailMessageIsNotSet:                   http.StatusBadRequest,
		Code_CaptchaError:                           http.StatusInternalServerError,
		Code_CaptchaCreationError:                   http.StatusInternalServerError,
		Code_CaptchaTaskIdIsNotSet:                  http.StatusBadRequest,
		Code_CaptchaAnswerIsNotSet:                  http.StatusBadRequest,
		Code_CaptchaCheckError:                      http.StatusInternalServerError,
		Code_CaptchaAnswerIsWrong:                   http.StatusForbidden,
		Code_CaptchaIsNotFound:                      http.StatusNotFound,
		Code_NameIsNotSet:                           http.StatusBadRequest,
		Code_EmailIsNotSet:                          http.StatusBadRequest,
		Code_PasswordIsNotSet:                       http.StatusBadRequest,
		Code_UserNameIsUsed:                         http.StatusConflict,
		Code_UserEmailIsUsed:                        http.StatusConflict,
		Code_RegistrationRequestWithUserNameExists:  http.StatusConflict,
		Code_RegistrationRequestWithUserEmailExists: http.StatusConflict,
		Code_UserEmailIsInvalid:                     http.StatusBadRequest,
		Code_UserNameIsTooLong:                      http.StatusBadRequest,
		Code_UserPasswordIsTooLong:                  http.StatusBadRequest,
		Code_UserPasswordIsNotAllowed:               http.StatusBadRequest,
		Code_RequestIdGenerator:                     http.StatusInternalServerError,
		Code_RequestIdIsNotSet:                      http.StatusBadRequest,
		Code_VerificationCodeIsNotSet:               http.StatusBadRequest,
		Code_BPP:                                    http.StatusInternalServerError,
		Code_LogInRequestWithUserEmailExists:        http.StatusConflict,
		Code_AuthDataIsNotSet:                       http.StatusBadRequest,
		Code_PasswordIsWrong:                        http.StatusForbidden,
		Code_JWT:                                    http.StatusInternalServerError,
		Code_VerificationCodeIsWrong:                http.StatusForbidden,
		Code_SessionIsNotFound:                      http.StatusForbidden,
		Code_NotAuthorised:                          http.StatusUnauthorized,
		Code_NotSure:                                http.StatusBadRequest,
		Code_TokenIsExpired:                         http.StatusForbidden,
		Code_EmailChangeRequestWithNewEmailExists:   http.StatusConflict,
		Code_PasswordChangeRequestWithUserIdExists:  http.StatusConflict,
		Code_PageIsNotSet:                           http.StatusBadRequest,
	}
}

const (
	ErrF_DatabaseNetwork = "database network error: %v"
)

func IsNetworkError(err error) (isNetworkError bool) {
	var nerr net.Error
	return errors.As(err, &nerr)
}
