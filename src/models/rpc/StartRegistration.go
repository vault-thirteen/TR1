package rm

type StartRegistrationParams struct {
	CommonParams

	// Fields provided by user.
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
}

type StartRegistrationResult struct {
	CommonResult
	RequestId string `json:"requestId"`
	CaptchaId string `json:"captchaId"`
}
