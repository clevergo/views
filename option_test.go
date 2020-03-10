// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import (
	"html/template"
	"reflect"
	"testing"
)

func TestCache(t *testing.T) {
	tests := []bool{false, true, false}
	for _, cache := range tests {
		m := New("", Cache(cache))
		if m.cache != cache {
			t.Errorf("expected cache %t, got %t", cache, m.cache)
		}
	}
}

func TestDelimis(t *testing.T) {
	tests := [][]string{
		{"{{", "}}"},
		{"{{{", "}}}"},
		{"{{#", "#}}"},
	}
	for _, test := range tests {
		m := New("", Delims(test[0], test[1]))
		if !reflect.DeepEqual(m.delims, test) {
			t.Errorf("expected delims %v, got %v", test, m.delims)
		}
	}
}

/*
func TestLayout(t *testing.T) {
	tests := [][]string{
		{},
		{"foo"},
		{"foo", "bar"},
	}
	for _, layouts := range tests {
		m := New("", Layouts(layouts...))
		if !reflect.DeepEqual(v.layouts, layouts) {
			t.Errorf("expected layouts %v, got %v", layouts, v.layouts)
		}
	}
}
*/

func TestSuffix(t *testing.T) {
	tests := []string{".tmpl", ".tpl", ".html", "htm"}
	for _, suffix := range tests {
		m := New("", Suffix(suffix))
		if m.suffix != suffix {
			t.Errorf("expected suffix %q, got %q", suffix, m.suffix)
		}
	}
}

func TestDefaultLayout(t *testing.T) {
	tests := []struct {
		layout string
	}{
		{"main"},
		{"page"},
	}
	for _, test := range tests {
		m := New("", DefaultLayout(test.layout))
		if m.defaultLayout != test.layout {
			t.Errorf("expected default layout %q, got %q", test.layout, m.defaultLayout)
		}
	}
}

func TestLayoutsDir(t *testing.T) {
	tests := []string{"layouts1", "layouts2"}
	for _, dir := range tests {
		m := New("", LayoutsDir(dir))
		if m.layoutsDir != dir {
			t.Errorf("expected layouts directory %q, got %q", dir, m.partialsDir)
		}
	}
}

func TestPartialsDir(t *testing.T) {
	tests := []string{"partials1", "partials2"}
	for _, dir := range tests {
		m := New("", PartialsDir(dir))
		if m.partialsDir != dir {
			t.Errorf("expected partials directory %q, got %q", dir, m.partialsDir)
		}
	}
}

func TestFuncMap(t *testing.T) {
	tests := []template.FuncMap{
		{
			"foo": func() string { return "foo" },
		},
		{
			"bar": func() string { return "bar" },
		},
	}
	for _, funcMap := range tests {
		m := New("", FuncMap(funcMap))
		for name := range funcMap {
			if _, ok := m.funcMap[name]; !ok {
				t.Errorf("failed to add function %s", name)
			}
		}
	}
}
