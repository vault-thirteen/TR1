package rm

type RejectRegistrationRequestRFAParams struct {
	CommonParams
	UserEmail string `json:"userEmail"`
}

type RejectRegistrationRequestRFAResult struct {
	CommonResult
	Success
}
