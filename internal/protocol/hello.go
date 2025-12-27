package protocol

type Hello struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}
