package service

import (
	"bytes"
	"errors"
	"github.com/wellington/go-libsass"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const STYLES_EXT = ".scss"
const STYLES_DIR = "styles"

type Styler struct {
	loadedAt time.Time         // loaded at (last loading time)
	styleMap map[string]string // Map of key => CSS Style Value
}

func NewStyler() *Styler {
	return &Styler{
		styleMap: make(map[string]string),
	}
}

func (s *Styler) Init() error {
	return s.Load()
}

func CompileStyleToString(path string) (string, error) {
	scss, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString(string(scss))
	var css bytes.Buffer
	comp, err := libsass.New(&css, buf)
	if err != nil {
		return "", err
	}
	if err := comp.Run(); err != nil {
		return "", err
	}
	return css.String(), nil
}

// Load or reload templates
func (s *Styler) Load() (err error) {

	// time point
	s.loadedAt = time.Now()

	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {

		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != STYLES_EXT {
			return
		}

		// get relative path
		var rel string
		if rel, err = filepath.Rel(STYLES_DIR, path); err != nil {
			return err
		}

		// name of a template is its relative path
		// without extension
		rel = strings.TrimSuffix(rel, STYLES_EXT)
		css, err := CompileStyleToString(path)
		if err != nil {
			return err
		}
		s.styleMap[rel] = css
		return err
	}

	if err = filepath.Walk(STYLES_DIR, walkFunc); err != nil {
		return
	}

	return
}

// IsModified lookups directory for changes to
// reload (or not to reload) templates if autoReloadopment
// pin is true.
func (s *Styler) IsModified() (yep bool, err error) {
	var errStop = errors.New("stop")
	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {

		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != STYLES_EXT {
			return
		}

		if yep = info.ModTime().After(s.loadedAt); yep == true {
			return errStop
		}

		return
	}
	// clear the errStop
	if err = filepath.Walk(STYLES_DIR, walkFunc); err == errStop {
		err = nil
	}
	return
}
