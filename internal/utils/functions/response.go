package functions

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ResponseProps struct {
	W           http.ResponseWriter
	Payload     any
	CodeStatus  int
	ContentType string
}

func NewResponseProps(w http.ResponseWriter, payload any, codeStatus int, contentType string) ResponseProps {
	return ResponseProps{
		W:           w,
		Payload:     payload,
		CodeStatus:  codeStatus,
		ContentType: contentType,
	}
}

// Response
// Forms server response. For marshalling uses easyjson (mailru).
func Response(props ResponseProps) {
	props.W.Header().Add("Content-Type", props.ContentType)
	body, err := json.Marshal(props.Payload)
	if err != nil {
		props.W.Header().Add("Content-Length", "0")
		props.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	props.W.WriteHeader(props.CodeStatus)
	contentLength, err := props.W.Write(body)
	if err != nil {
		props.W.WriteHeader(http.StatusInternalServerError)
	}
	props.W.Header().Add("Content-Length", strconv.Itoa(contentLength))
}
