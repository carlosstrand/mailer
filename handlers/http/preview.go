package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

func (h *HTTPHandler) PreviewerHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data := sync.Map{}
	data.Store("mode", "html")
	output, err := h.mailer.RenderToString(vars["tmplName"], &data, true, false)

	if err != nil {
		res.WriteHeader(404)
		res.Write([]byte(err.Error()))
	}

	res.Write([]byte(output))
}
