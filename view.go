package views

import "html/template"

// View wraps template.Template.
type View struct {
	*template.Template
}
