// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
)

type layout struct {
	name     string
	partials []string
}

// Manager is the views manager.
type Manager struct {
	fs            http.FileSystem
	path          string
	defaultLayout string
	layouts       map[string]*layout
	layoutsDir    string
	partialsDir   string
	suffix        string
	delims        []string
	funcMap       template.FuncMap
	cache         bool
	mutex         *sync.Mutex
	views         map[string]map[string]*View
}

// New returns a manager with the given filesystem and options.
func New(fs http.FileSystem, opts ...Option) *Manager {
	m := &Manager{
		fs:            fs,
		suffix:        ".tmpl",
		delims:        []string{"{{", "}}"},
		mutex:         &sync.Mutex{},
		defaultLayout: "main",
		layoutsDir:    "layouts",
		partialsDir:   "partials",
		cache:         true,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// AddLayout adds a layout with the given name and partials.
func (m *Manager) AddLayout(name string, partials ...string) {
	if m.layouts == nil {
		m.layouts = make(map[string]*layout)
	}
	m.layouts[name] = &layout{name, partials}
}

// AddFunc add function to funcMap.
func (m *Manager) AddFunc(name string, f interface{}) {
	if m.funcMap == nil {
		m.funcMap = template.FuncMap{}
	}
	m.funcMap[name] = f
}

// Render renders a view with default layout.
func (m *Manager) Render(w io.Writer, view string, data interface{}) error {
	return m.RenderLayout(w, m.defaultLayout, view, data)
}

// RenderLayout renders a view with particular layout.
func (m *Manager) RenderLayout(w io.Writer, layout, view string, data interface{}) error {
	return m.render(w, layout, view, data)
}

// RenderPartial renders a view without layout.
func (m *Manager) RenderPartial(w io.Writer, view string, data interface{}) error {
	return m.render(w, "", view, data)
}

func (m *Manager) getView(layout, view string) (*View, error) {
	if v, ok := m.views[layout][view]; ok {
		return v, nil
	}

	files := []string{}
	if layout != "" {
		l, ok := m.layouts[layout]
		if !ok {
			return nil, fmt.Errorf("no such layout %q", layout)
		}
		files = append(files, m.findLayoutFile(l.name))
		for _, partial := range l.partials {
			files = append(files, m.findPartialFile(partial))
		}
	}
	files = append(files, m.findViewFile(view))

	v, err := m.newView(files)
	if err != nil {
		return nil, err
	}

	if m.cache {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		if m.views == nil {
			m.views = make(map[string]map[string]*View)
		}
		if _, ok := m.views[layout]; !ok {
			m.views[layout] = make(map[string]*View)
		}
		m.views[layout][view] = v
	}

	return v, nil
}

func (m *Manager) newView(files []string) (*View, error) {
	tmpl := template.New(path.Base(files[0])).
		Funcs(m.funcMap).
		Delims(m.delims[0], m.delims[1])

	for _, filename := range files {
		file, err := m.fs.Open(filename)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		tmpl, err = tmpl.Parse(string(content))
		if err != nil {
			return nil, err
		}
	}

	return &View{tmpl}, nil
}

func (m *Manager) findViewFile(view string) string {
	return path.Join(m.getFileName(view))
}

func (m *Manager) findLayoutFile(layout string) string {
	return path.Join(m.layoutsDir, m.getFileName(layout))
}

func (m *Manager) findPartialFile(partial string) string {
	return path.Join(m.layoutsDir, m.partialsDir, m.getFileName(partial))
}

func (m *Manager) getFileName(view string) string {
	return view + m.suffix
}

func (m *Manager) render(w io.Writer, layout, view string, data interface{}) error {
	v, err := m.getView(layout, view)
	if err != nil {
		return err
	}

	return v.Execute(w, data)
}
