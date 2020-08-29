package mailer

import (
	"fmt"
	httphandler "github.com/carlosstrand/mailer/handlers/http"
	"github.com/carlosstrand/mailer/service"
	"github.com/carlosstrand/mailer/utils"
	"github.com/go-zepto/zepto"
	"github.com/gorilla/mux"
	"net/http"
)

func (m *Mailer) Start() {

	fmt.Println("Mailer is Starting... 📨")

	host := utils.GetEnv("APP_HOST", "0.0.0.0")
	port := utils.GetEnv("APP_PORT", "8000")

	templates := utils.FindAvailableTemplates()
	err := m.Build(templates)
	if err != nil {
		panic(err)
	}
	utils.PrintAvailableTemplates(templates)

	// Create Zepto
	z := zepto.NewZepto(
		zepto.Name("mailer"),
		zepto.Version("latest"),
	)

	mailer := service.NewMailerService(&m.config)
	router := mux.NewRouter()

	p := httphandler.NewHTTPHandler(mailer)
	router.HandleFunc("/preview/{tmplName}", p.PreviewerHandler).Methods("GET")
	router.HandleFunc("/send", p.SendHandler).Methods("POST")

	publicDir := "/public/"

	// Create the route
	router.
		PathPrefix(publicDir).
		Handler(http.StripPrefix(publicDir, http.FileServer(http.Dir("."+publicDir))))

	// Setup HTTP Server
	z.SetupHTTP(host+":"+port, router)

	if err := mailer.Init(); err != nil {
		panic(err)
	}
	z.Start()
}
