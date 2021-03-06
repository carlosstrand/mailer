package service

import (
	"bufio"
	"bytes"
	"github.com/carlosstrand/mailer/types"
	pongo2lib "github.com/flosch/pongo2"
	"github.com/go-zepto/zepto/web/renderer/pongo2"
	"github.com/vanng822/go-premailer/premailer"
	"strings"
	"sync"
)

type MailTmpl struct {
	Template string
	Styles   []string
}

type MailerService struct {
	config   *types.MailerConfig
	renderer *pongo2.Pongo2Engine
	styler   *Styler
}

func NewMailerService(config *types.MailerConfig) *MailerService {
	renderer := pongo2.NewPongo2Engine(pongo2.AutoReload(true))
	styler := NewStyler()
	return &MailerService{
		config:   config,
		renderer: renderer,
		styler:   styler,
	}
}

func (m *MailerService) Init() error {
	pongo2lib.RegisterTag("public", tagPublicParser)
	pongo2lib.RegisterTag("var", tagVarParser)
	if err := m.styler.Init(); err != nil {
		return err
	}
	return m.renderer.Init()
}

func (m *MailerService) RenderToString(mailTmpl string, data *sync.Map, withPremailer bool, builded bool) (string, error) {

	var buffer bytes.Buffer
	wr := MailerWriter{Writer: bufio.NewWriter(&buffer)}

	if builded {
		err := m.renderer.Render(&wr, 200, "build/"+mailTmpl, data)
		if err != nil {
			return "", err
		}
		return string(wr.Value()), nil
	}

	// Auto-reload styles
	if modified, err := m.styler.IsModified(); err == nil && modified {
		m.styler.Load()
	}

	stylesCSS := ""

	// Include styled related to template (same name)
	if style, exists := m.styler.styleMap[mailTmpl]; exists {
		stylesCSS += style + "\n"
	}

	// Include styles from global scope
	for key, value := range m.styler.styleMap {
		if strings.HasPrefix(key, "global/") {
			stylesCSS += value + "\n"
		}
	}

	data.Store("styles", stylesCSS)
	data.Store("publicPath", m.config.PublicPath)
	err := m.renderer.Render(&wr, 200, "mails/"+mailTmpl, data)
	if err != nil {
		return "", err
	}

	if !withPremailer {
		return string(wr.Value()), nil
	}

	prem, err := premailer.NewPremailerFromBytes(wr.Value(), &premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   true,
		KeepBangImportant: false,
	})

	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
