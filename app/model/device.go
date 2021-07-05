package model


type NewDeviceReq struct {
	IP        string `json:"ip"`
	UserName     string
	Password        string
}

type DiscoveryRes struct {
	URL   string   `json:"url"`
	Name  string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	VideoSourceNumber string  `json:"videoSourceNumber"`
}