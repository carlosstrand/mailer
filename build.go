package mailer

import (
	"github.com/carlosstrand/mailer/service"
	"github.com/carlosstrand/mailer/types"
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
	mailer := service.NewMailerService(&types.MailerConfig{})
	if err := mailer.Init(); err != nil {
		panic(err)
	}
	for _, key := range templates {
		html, err := mailer.RenderToString(key, &sync.Map{}, true, false)
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
