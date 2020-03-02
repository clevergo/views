// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import (
	"bytes"
	"html"
	"html/template"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var (
	testView      *View
	testCacheView *View
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	themesPath := path.Join(path.Dir(filename), "example", "themes")
	testView = New(
		themesPath,
		Theme("default"),
		Layouts("layouts/main", "layouts/header", "layouts/footer"),
		FuncMap(template.FuncMap{
			"title": strings.Title,
		}),
	)
	testCacheView = &View{}
	*testCacheView = *testView
	testCacheView.cache = true
	m.Run()
}

func TestView_Render(t *testing.T) {
	w := bytes.NewBuffer(nil)
	err := testView.Render(w, "site/index", map[string]interface{}{
		"title": "home",
	})
	if err != nil {
		t.Errorf("failed to render: %s", err)
	}
	strs := []string{
		"<h1>Hello World</h1>",
		"<title>Home</title>",
		"<head>",
		"<footer>I am footer</footer>",
	}
	for _, s := range strs {
		if !bytes.Contains(w.Bytes(), []byte(s)) {
			t.Errorf("render result doesn't contains %q", s)
		}
	}
}

func TestView_Partial(t *testing.T) {
	w := bytes.NewBuffer(nil)
	err := testView.RenderPartial(w, "site/partial", map[string]interface{}{
		"title": "standalone",
	})
	if err != nil {
		t.Errorf("failed to render: %s", err)
	}
	strs := []string{
		"<h1>Standalone</h1>",
		"<title>Standalone</title>",
		"<head>",
	}
	for _, s := range strs {
		if !bytes.Contains(w.Bytes(), []byte(s)) {
			t.Errorf("render result doesn't contains %q", s)
		}
	}
}

func TestView_GetTemplate(t *testing.T) {
	v := &View{}
	*v = *testView
	v.cache = true

	tests := []struct {
		layout bool
		view   string
	}{
		{true, "site/index"},
		{false, "site/partial"},
	}
	for _, test := range tests {
		if test.layout {
			v.Render(bytes.NewBuffer(nil), test.view, nil)
		} else {
			v.RenderPartial(bytes.NewBuffer(nil), test.view, nil)
		}
		cachedTemaplte, ok := v.templates[test.layout][test.view]
		if !ok {
			t.Errorf("failed to cache template: %s", test.view)
		}
		tmpl, err := v.getTemplate(test.view, test.layout)
		if err != nil || !reflect.DeepEqual(cachedTemaplte, tmpl) {
			t.Errorf("failed to retrieve cached template: %s", test.view)
		}
	}

	err := v.Render(bytes.NewBuffer(nil), "nonexistent", nil)
	if err == nil {
		t.Error("expected an error about view file not found, got nil")
	}
	err = v.RenderPartial(bytes.NewBuffer(nil), "nonexistent", nil)
	if err == nil {
		t.Error("expected an error about view file not found, got nil")
	}
}

func BenchmarkView_Render(b *testing.B) {
	data := map[string]interface{}{
		"title": "home",
	}
	w := bytes.NewBuffer(nil)
	for n := 0; n < b.N; n++ {
		testView.Render(w, "site/index", data)
		w.Reset()
	}
}

func BenchmarkView_RenderPartial(b *testing.B) {
	data := map[string]interface{}{
		"title": "standalone",
	}
	w := bytes.NewBuffer(nil)
	for n := 0; n < b.N; n++ {
		testView.RenderPartial(w, "site/partial", data)
		w.Reset()
	}
}

func BenchmarkCacheView_Render(b *testing.B) {
	data := map[string]interface{}{
		"title": "home",
	}
	w := bytes.NewBuffer(nil)
	for n := 0; n < b.N; n++ {
		testCacheView.Render(w, "site/index", data)
		w.Reset()
	}
}

func BenchmarkCacheView_RenderPartial(b *testing.B) {
	data := map[string]interface{}{
		"title": "standalone",
	}
	w := bytes.NewBuffer(nil)
	for n := 0; n < b.N; n++ {
		testCacheView.RenderPartial(w, "site/partial", data)
		w.Reset()
	}
}

func TestAddFuncMap(t *testing.T) {
	v := &View{}
	v.AddFunc("title", strings.Title)
	v.AddFunc("escapeString", html.EscapeString)
	for _, name := range []string{"title", "escapeString"} {
		_, ok := v.funcMap[name]
		if !ok {
			t.Errorf("failed to add func: %s", name)
		}
	}
}
