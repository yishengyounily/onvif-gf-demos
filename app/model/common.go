package model


type CommonRes struct {
	Code        int64 `json:"ErrorCode"`
	Message     string
	Data        interface{}
	RefreshTime int64
}