package handlers

import (
    "github.com/carlosstrand/mailer/service"
    "github.com/gorilla/mux"
    "net/http"
    "sync"
)

type Previewer struct {
    mailer *service.MailerService
}

func NewPreviewer(mailer *service.MailerService)*Previewer {
    return &Previewer{mailer: mailer}
}

func (p *Previewer) PreviewerHandler (res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    data := sync.Map{}
    output, err := p.mailer.RenderToString(vars["tmplName"], &data,true, true)

    if err != nil {
        res.WriteHeader(404)
        res.Write([]byte(err.Error()))
    }

    res.Write([]byte(output))
}
