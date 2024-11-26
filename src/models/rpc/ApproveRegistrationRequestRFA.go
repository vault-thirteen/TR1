package rm

type ApproveRegistrationRequestRFAParams struct {
	CommonParams
	UserEmail string `json:"userEmail"`
}

type ApproveRegistrationRequestRFAResult struct {
	CommonResult
	Success
}
