package entity

type ResponseError struct {
	Error string `json:"error" valid:"-"`
}

type ResponseDetail struct {
	Detail string `json:"detail" valid:"-"`
}

type ResponseReason struct {
	Reason string `json:"reason" valid:"-"`
}

