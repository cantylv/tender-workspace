package entity

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseDetail struct {
	Detail string `json:"detail"`
}

type ResponseReason struct {
	Reason string `json:"reason"`
}
