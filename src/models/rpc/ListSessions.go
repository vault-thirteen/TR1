package rm

type ListSessionsParams struct {
	CommonParams
	PageRequested
}

type ListSessionsResult struct {
	CommonResult
	ItemsPaginated
}
