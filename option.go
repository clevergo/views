// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import "html/template"

// Option is a function that applies on View.
type Option func(*View)

// Cache enables or disables cache.
func Cache(cache bool) Option {
	return func(v *View) {
		v.cache = cache
	}
}

// Delims sets the delimiters.
func Delims(left, right string) Option {
	return func(v *View) {
		v.delims[0] = left
		v.delims[1] = right
	}
}

// Suffix sets the suffix.
func Suffix(suffix string) Option {
	return func(v *View) {
		v.suffix = suffix
	}
}

// Theme sets the theme.
func Theme(theme string) Option {
	return func(v *View) {
		v.theme = theme
	}
}

// Layouts sets the layouts.
func Layouts(layouts ...string) Option {
	return func(v *View) {
		v.layouts = layouts
	}
}

// FuncMap sets the global function map of all templates.
func FuncMap(funcMap template.FuncMap) Option {
	return func(v *View) {
		v.funcMap = funcMap
	}
}
