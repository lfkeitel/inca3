package models

type JobApiRequest struct {
	Type    string `json:"type"`
	Devices []int  `json:"devices"`
	Start   int64  `json:"start"`
}
