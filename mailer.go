package mailer

import (
    "fmt"
    "github.com/carlosstrand/mailer/handlers"
    "github.com/carlosstrand/mailer/service"
    "github.com/carlosstrand/mailer/utils"
    "github.com/fatih/color"
    "github.com/go-zepto/zepto"
    "github.com/gorilla/mux"
    "github.com/tdewolff/minify/v2"
    "github.com/tdewolff/minify/v2/html"
    "github.com/tdewolff/parse/v2/buffer"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

type MailerConfig struct {
    DefaultFrom string
    PublicPath string
}

type Mailer struct {
    config MailerConfig
}

func NewMailer(config MailerConfig)*Mailer {
    return &Mailer{config: config}
}

func (m *Mailer) Build(templates []string) error {
    minifier := minify.New()
    minifier.AddFunc("text/html", html.Minify)
    mailer := service.NewMailerService()
    if err := mailer.Init(); err != nil {
        panic(err)
    }
    for _, key := range templates {
        html, err := mailer.RenderToString(key, &sync.Map{},true, false)
        if err != nil {
            return err
        }
        _ = os.Mkdir("templates/build/", 0700)
        f, err := os.Create("templates/build/" + key + ".html")
        if err != nil {
            return err
        }
        in := strings.NewReader(html)
        var res buffer.Writer
        if err := minifier.Minify("text/html", &res, in); err != nil {
            return err
        }
        _, err = f.WriteString(string(res.Bytes()))
        if err != nil {
            return err
        }
    }
    return nil
}

func (m *Mailer) FindAvailableTemplates() []string {
    var paths []string
    var files []string

    root := "templates/mails"
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if strings.HasSuffix(path, ".html") {
            paths = append(paths, path)
        }
        return nil
    })
    if err != nil {
        panic(err)
    }
    for _, path := range paths {
        file := strings.Replace(filepath.Base(path), ".html", "", 1)
        files = append(files, file)
    }
    return files
}

func (m *Mailer) Start() {

    fmt.Println("Mailer is Starting... 📨")

    host := utils.GetEnv("APP_HOST", "0.0.0.0")
    port := utils.GetEnv("APP_PORT", "8000")


    templates := m.FindAvailableTemplates()

    err := m.Build(templates)
    if err != nil {
        panic(err)
    }

    if len(templates) > 0 {
        fmt.Println("-----------------------------------------------------")
        fmt.Printf("There are %d templates available:\n", len(templates))
        for _, t := range templates {
            color.Green("- %s", t)
        }
        fmt.Println("-----------------------------------------------------")
    }

    // Create Zepto
    z := zepto.NewZepto(
        zepto.Name("mailer"),
        zepto.Version("latest"),
    )

    mailer := service.NewMailerService()
    router := mux.NewRouter()

    p := handlers.NewPreviewer(mailer)
    router.HandleFunc("/preview/{tmplName}", p.PreviewerHandler).Methods("GET")

    publicDir := "/public/"

    // Create the route
    router.
        PathPrefix(publicDir).
        Handler(http.StripPrefix(publicDir, http.FileServer(http.Dir("."+publicDir))))

    // Setup HTTP Server
    z.SetupHTTP(host + ":" + port, router)


    if err := mailer.Init(); err != nil {
        panic(err)
    }
    z.Start()

}
