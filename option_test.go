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
		v := New("", Cache(cache))
		if v.cache != cache {
			t.Errorf("expected cache %t, got %t", cache, v.cache)
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
		v := New("", Delims(test[0], test[1]))
		if !reflect.DeepEqual(v.delims, test) {
			t.Errorf("expected delims %v, got %v", test, v.delims)
		}
	}
}

func TestLayout(t *testing.T) {
	tests := [][]string{
		{},
		{"foo"},
		{"foo", "bar"},
	}
	for _, layouts := range tests {
		v := New("", Layouts(layouts...))
		if !reflect.DeepEqual(v.layouts, layouts) {
			t.Errorf("expected layouts %v, got %v", layouts, v.layouts)
		}
	}
}

func TestSuffix(t *testing.T) {
	tests := []string{".tmpl", ".tpl", ".html", "htm"}
	for _, suffix := range tests {
		v := New("", Suffix(suffix))
		if v.suffix != suffix {
			t.Errorf("expected suffix %q, got %q", suffix, v.suffix)
		}
	}
}

func TestTheme(t *testing.T) {
	tests := []string{"foo", "bar"}
	for _, theme := range tests {
		v := New("", Theme(theme))
		if v.theme != theme {
			t.Errorf("expected theme %q, got %q", theme, v.theme)
		}
	}
}

func TestFuncMap(t *testing.T) {
	tests := []template.FuncMap{
		{},
		{
			"foo": func() string { return "foo" },
		},
	}
	for _, funcMap := range tests {
		v := New("", FuncMap(funcMap))
		if !reflect.DeepEqual(v.funcMap, funcMap) {
			t.Errorf("expected funcMap %v, got %v", funcMap, v.funcMap)
		}
	}
}
