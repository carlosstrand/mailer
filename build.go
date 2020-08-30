package mailer

import (
	"github.com/carlosstrand/mailer/service"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/parse/v2/buffer"
	"os"
	"strings"
	"sync"
)

func (m *Mailer) Build(templates []string) error {
	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)
	mailer := service.NewMailerService(m.config)
	if err := mailer.Init(); err != nil {
		panic(err)
	}
	data := &sync.Map{}
	data.Store("mode", "html")
	for _, key := range templates {
		html, err := mailer.RenderToString(key, data, true, false)
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
