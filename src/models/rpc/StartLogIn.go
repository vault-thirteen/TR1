package rm

type StartLogInParams struct {
	CommonParams

	// Fields provided by user.
	UserEmail string `json:"userEmail"`
}

type StartLogInResult struct {
	CommonResult
	RequestId string `json:"requestId"`
	CaptchaId string `json:"captchaId"`
	AuthData  []byte `json:"authData"`
}
