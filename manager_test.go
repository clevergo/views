// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package views

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"html/template"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var (
	testManager      *Manager
	testCacheManager *Manager
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	viewPath := path.Join(path.Dir(filename), "example", "views")
	testManager = New(
		viewPath,
		FuncMap(template.FuncMap{
			"title": strings.Title,
		}),
	)
	testManager.AddLayout("main", "head", "header", "footer")
	testCacheManager = &Manager{}
	*testCacheManager = *testManager
	testCacheManager.cache = true
	m.Run()
}

func TestManagerRender(t *testing.T) {
	m := &Manager{}
	*m = *testManager
	w := bytes.NewBuffer(nil)
	err := m.Render(w, "site/index", map[string]interface{}{
		"title": "home",
	})
	if err != nil {
		t.Fatalf("failed to render: %s", err)
	}
	strs := []string{
		"<h1>Hello World</h1>",
		"<title>Home</title>",
		"<head>",
		"<header>Header</header>",
		"<footer>Footer</footer>",
	}
	for _, s := range strs {
		if !bytes.Contains(w.Bytes(), []byte(s)) {
			t.Errorf("render result doesn't contains %q", s)
		}
	}
	errBeforeRender := errors.New("before render")
	m.RegisterOnBeforeRender(func(_ *BeforeRenderEvent) error {
		return errBeforeRender
	})
	err = m.Render(w, "site/index", map[string]interface{}{
		"title": "home",
	})
	if err == nil || err != errBeforeRender {
		t.Errorf("expected error %q, got %q", errBeforeRender, err)
	}

	m.onBeforeRender = m.onBeforeRender[:0]
	errAfterRender := errors.New("after render")
	m.RegisterOnAfterRender(func(_ *AfterRenderEvent) error {
		return errAfterRender
	})
	err = m.Render(w, "site/index", map[string]interface{}{
		"title": "home",
	})
	if err == nil || err != errAfterRender {
		t.Errorf("expected error %q, got %q", errAfterRender, err)
	}
}

func TestManagerPartial(t *testing.T) {
	w := bytes.NewBuffer(nil)
	err := testManager.RenderPartial(w, "site/partial", Context{
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

func TestManagerGetView(t *testing.T) {
	m := &Manager{}
	_, err := m.getView("invalid", "view")
	expcetedErr := fmt.Sprintf("no such layout %q", "invalid")
	if err == nil || err.Error() != expcetedErr {
		t.Errorf("expected error %s, got %s", expcetedErr, err)
	}

	*m = *testManager
	m.cache = true

	tests := []struct {
		layout string
		view   string
	}{
		{"main", "site/index"},
		{"", "site/partial"},
	}
	for _, test := range tests {
		if test.layout != "" {
			m.Render(bytes.NewBuffer(nil), test.view, nil)
		} else {
			m.RenderPartial(bytes.NewBuffer(nil), test.view, nil)
		}
		cachedV, ok := m.views[test.layout][test.view]
		if !ok {
			t.Fatalf("failed to cache view: %s", test.view)
		}
		v, err := m.getView(test.layout, test.view)
		if err != nil || !reflect.DeepEqual(cachedV, v) {
			t.Errorf("failed to retrieve cached view: %s", test.view)
		}
	}

	err = m.Render(bytes.NewBuffer(nil), "nonexistent", nil)
	if err == nil {
		t.Error("expected an error about view file not found, got nil")
	}
	err = m.RenderPartial(bytes.NewBuffer(nil), "nonexistent", nil)
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
		testManager.Render(w, "site/index", data)
		w.Reset()
	}
}

func BenchmarkView_RenderPartial(b *testing.B) {
	data := map[string]interface{}{
		"title": "standalone",
	}
	w := bytes.NewBuffer(nil)
	for n := 0; n < b.N; n++ {
		testManager.RenderPartial(w, "site/partial", data)
		w.Reset()
	}
}

func BenchmarkCacheView_Render(b *testing.B) {
	data := map[string]interface{}{
		"title": "home",
	}
	w := bytes.NewBuffer(nil)
	for n := 0; n < b.N; n++ {
		testCacheManager.Render(w, "site/index", data)
		w.Reset()
	}
}

func BenchmarkCacheView_RenderPartial(b *testing.B) {
	data := map[string]interface{}{
		"title": "standalone",
	}
	w := bytes.NewBuffer(nil)
	for n := 0; n < b.N; n++ {
		testCacheManager.RenderPartial(w, "site/partial", data)
		w.Reset()
	}
}

func TestManagerAddFuncMap(t *testing.T) {
	m := &Manager{}
	m.AddFunc("title", strings.Title)
	m.AddFunc("escapeString", html.EscapeString)
	for _, name := range []string{"title", "escapeString"} {
		_, ok := m.funcMap[name]
		if !ok {
			t.Errorf("failed to add func: %s", name)
		}
	}
}

func TestManagerOnBeforeRender(t *testing.T) {
	res := ""
	f := func(event *BeforeRenderEvent) error {
		res += "foo"
		return nil
	}
	expectedErr := errors.New("BeforeRenderEvent err")
	f2 := func(event *BeforeRenderEvent) error {
		res += "bar"
		return expectedErr
	}
	m := New("")
	m.RegisterOnBeforeRender(f)
	m.RegisterOnBeforeRender(f2)
	err := m.beforeRender(nil, "", "", Context{})
	if res != "foobar" {
		t.Errorf("expected result %q, got %q", "foobar", res)
	}
	if err == nil || err != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err)
	}
}

func TestManagerOnAfterRender(t *testing.T) {
	res := ""
	f := func(event *AfterRenderEvent) error {
		res += "foo"
		return nil
	}
	expectedErr := errors.New("AfterRenderEvent err")
	f2 := func(event *AfterRenderEvent) error {
		res += "bar"
		return expectedErr
	}
	m := New("")
	m.RegisterOnAfterRender(f)
	m.RegisterOnAfterRender(f2)
	err := m.afterRender(nil, "", "", Context{})
	if res != "foobar" {
		t.Errorf("expected result %q, got %q", "foobar", res)
	}
	if err == nil || err != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err)
	}
}
