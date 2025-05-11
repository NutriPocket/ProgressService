package model

// ErrorRfc9457 is a struct that will be used to return errors in the RFC 9457 format
type ErrorRfc9457 struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}
