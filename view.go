// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import (
	"html/template"
	"io"
	"path"
	"sync"
)

// View is the templates manager.
type View struct {
	directory string
	theme     string
	layouts   []string
	suffix    string
	delims    []string
	funcMap   template.FuncMap
	cache     bool
	mutex     *sync.Mutex
	templates map[bool]map[string]*template.Template
}

// New returns a view with the given directory and options.
func New(dir string, opts ...Option) *View {
	v := &View{
		directory: dir,
		suffix:    ".tmpl",
		delims:    []string{"{{", "}}"},
		mutex:     &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

// Cache enables or disables cache.
func (v *View) Cache(cache bool) {
	v.cache = cache
}

// SetDelims sets the delimiters.
func (v *View) SetDelims(left, right string) {
	v.delims[0] = left
	v.delims[1] = right
}

// SetSuffix sets the suffix.
func (v *View) SetSuffix(suffix string) {
	v.suffix = suffix
}

// SetTheme sets the theme.
func (v *View) SetTheme(theme string) {
	v.theme = theme
}

// SetLayouts sets the layouts.
func (v *View) SetLayouts(layouts ...string) {
	v.layouts = layouts
}

// SetFuncMap sets the global function map of all templates.
func (v *View) SetFuncMap(funcMap template.FuncMap) {
	v.funcMap = funcMap
}

// Render executes a template with layouts.
func (v *View) Render(w io.Writer, view string, data interface{}) error {
	tmpl, err := v.getTemplate(view, true)
	if err != nil {
		return err
	}

	return v.execute(tmpl, w, data)
}

// RenderPartial executes a template without layouts.
func (v *View) RenderPartial(w io.Writer, view string, data interface{}) error {
	tmpl, err := v.getTemplate(view, false)
	if err != nil {
		return err
	}

	return v.execute(tmpl, w, data)
}

func (v *View) getTemplate(view string, layout bool) (*template.Template, error) {
	if tmpl, ok := v.templates[layout][view]; ok {
		return tmpl, nil
	}

	files := []string{}
	if layout {
		for _, f := range v.layouts {
			files = append(files, v.findViewFile(f))
		}
		files = append(files, v.findViewFile(view))
	}
	files = append(files, v.findViewFile(view))

	tmpl, err := template.New(path.Base(files[0])).
		Funcs(v.funcMap).
		Delims(v.delims[0], v.delims[1]).
		ParseFiles(files...)
	if err != nil {
		return tmpl, err
	}

	if v.cache {
		v.mutex.Lock()
		defer v.mutex.Unlock()
		if v.templates == nil {
			v.templates = map[bool]map[string]*template.Template{
				false: make(map[string]*template.Template),
				true:  make(map[string]*template.Template),
			}
		}
		v.templates[layout][view] = tmpl
	}

	return tmpl, nil
}

func (v *View) findViewFile(view string) string {
	return path.Join(v.directory, v.theme, view+v.suffix)
}

func (v *View) execute(tmpl *template.Template, w io.Writer, data interface{}) error {
	return tmpl.Execute(w, data)
}
