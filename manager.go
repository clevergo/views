// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import (
	"fmt"
	"html/template"
	"io"
	"path"
	"sync"
)

type layout struct {
	name     string
	partials []string
}

// BeforeRenderEvent is an evnet that trigger before rendering an view.
type BeforeRenderEvent struct {
	w      io.Writer
	layout string
	view   string
	ctx    Context
}

// AfterRenderEvent is an evnet that trigger after rendering an view.
type AfterRenderEvent struct {
	w      io.Writer
	layout string
	view   string
	ctx    Context
}

// Manager is the views manager.
type Manager struct {
	path           string
	defaultLayout  string
	layouts        map[string]*layout
	layoutsDir     string
	partialsDir    string
	suffix         string
	delims         []string
	funcMap        template.FuncMap
	cache          bool
	mutex          *sync.Mutex
	views          map[string]map[string]*View
	onBeforeRender []func(*BeforeRenderEvent) error
	onAfterRender  []func(*AfterRenderEvent) error
}

// New returns a manager with the given path and options.
func New(path string, opts ...Option) *Manager {
	m := &Manager{
		path:          path,
		suffix:        ".tmpl",
		delims:        []string{"{{", "}}"},
		mutex:         &sync.Mutex{},
		defaultLayout: "main",
		layoutsDir:    "layouts",
		partialsDir:   "layouts/partials",
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
func (m *Manager) Render(w io.Writer, view string, ctx Context) error {
	return m.RenderLayout(w, m.defaultLayout, view, ctx)
}

// RenderLayout renders a view with particular layout.
func (m *Manager) RenderLayout(w io.Writer, layout, view string, ctx Context) error {
	return m.render(w, layout, view, ctx)
}

// RenderPartial renders a view without layout.
func (m *Manager) RenderPartial(w io.Writer, view string, ctx Context) error {
	return m.render(w, "", view, ctx)
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
	tmpl, err := template.New(path.Base(files[0])).
		Funcs(m.funcMap).
		Delims(m.delims[0], m.delims[1]).
		ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return &View{tmpl}, nil
}

func (m *Manager) findViewFile(view string) string {
	return path.Join(m.path, m.getFileName(view))
}

func (m *Manager) findLayoutFile(layout string) string {
	return path.Join(m.path, m.layoutsDir, m.getFileName(layout))
}

func (m *Manager) findPartialFile(partial string) string {
	return path.Join(m.path, m.partialsDir, m.getFileName(partial))
}

func (m *Manager) getFileName(view string) string {
	return view + m.suffix
}

func (m *Manager) render(w io.Writer, layout, view string, ctx Context) error {
	v, err := m.getView(layout, view)
	if err != nil {
		return err
	}

	if err = m.beforeRender(w, layout, view, ctx); err != nil {
		return err
	}

	if err = v.Execute(w, ctx); err != nil {
		return err
	}

	return m.afterRender(w, layout, view, ctx)
}

// RegisterOnBeforeRender registers a BeforeRenderEvent listener.
func (m *Manager) RegisterOnBeforeRender(f func(*BeforeRenderEvent) error) {
	m.onBeforeRender = append(m.onBeforeRender, f)
}

func (m *Manager) beforeRender(w io.Writer, layout, view string, ctx Context) (err error) {
	if m.onBeforeRender != nil {
		event := &BeforeRenderEvent{w, layout, view, ctx}
		for _, f := range m.onBeforeRender {
			if err = f(event); err != nil {
				return
			}
		}
	}
	return
}

// RegisterOnAfterRender registers a BeforeRenderEvent listener.
func (m *Manager) RegisterOnAfterRender(f func(*AfterRenderEvent) error) {
	m.onAfterRender = append(m.onAfterRender, f)
}

func (m *Manager) afterRender(w io.Writer, layout, view string, ctx Context) (err error) {
	if m.onAfterRender != nil {
		event := &AfterRenderEvent{w, layout, view, ctx}
		for _, f := range m.onAfterRender {
			if err = f(event); err != nil {
				return
			}
		}
	}

	return
}
