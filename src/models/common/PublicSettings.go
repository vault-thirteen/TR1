package cm

type PublicSettings struct {
	SiteName           string `json:"siteName"`
	SiteDomain         string `json:"siteDomain"`
	SessionMaxDuration int    `json:"sessionMaxDuration"`
	MessageEditTime    int    `json:"messageEditTime"`
	PageSize           int    `json:"pageSize"`
}
