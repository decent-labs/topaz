package models

type DeployResponse struct {
	Tx   string `json:"tx"`
	Addr string `json:"addr"`
}

type StoreRequest struct {
	Address string
	Hash    string
}

type StoreResponse struct {
	Hash string `json:"hash"`
}
