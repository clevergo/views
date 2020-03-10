// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import "html/template"

// Option is a function that applies on View.
type Option func(*Manager)

// Cache enables or disables cache.
func Cache(v bool) Option {
	return func(m *Manager) {
		m.cache = v
	}
}

// Delims sets the delimiters.
func Delims(left, right string) Option {
	return func(m *Manager) {
		m.delims[0] = left
		m.delims[1] = right
	}
}

// DefaultLayout sets the default layout.
func DefaultLayout(name string) Option {
	return func(m *Manager) {
		m.defaultLayout = name
	}
}

// LayoutsDir sets the layouts directory.
func LayoutsDir(dir string) Option {
	return func(m *Manager) {
		m.layoutsDir = dir
	}
}

// PartialsDir sets the partials directory.
func PartialsDir(dir string) Option {
	return func(m *Manager) {
		m.partialsDir = dir
	}
}

// Suffix sets the suffix.
func Suffix(suffix string) Option {
	return func(m *Manager) {
		m.suffix = suffix
	}
}

// FuncMap sets the global function map of all templates.
func FuncMap(funcMap template.FuncMap) Option {
	return func(m *Manager) {
		for name, f := range funcMap {
			m.AddFunc(name, f)
		}
	}
}
