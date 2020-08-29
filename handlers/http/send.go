package http

import (
	"encoding/json"
	"github.com/carlosstrand/mailer/types"
	"net/http"
)

func (h *HTTPHandler) SendHandler(res http.ResponseWriter, req *http.Request) {
	var sendReq types.SendMailFromTemplateRequest
	err := json.NewDecoder(req.Body).Decode(&sendReq)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	err = h.mailer.SendFromReq(sendReq)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
}
