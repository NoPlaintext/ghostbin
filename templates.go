package main

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"html/template"
	"io"
	"net/http"
)

type RenderContext struct {
	Obj      interface{}
	Request  *http.Request
	Page     string
	template *template.Template
}

var templateFunctions template.FuncMap = template.FuncMap{}
var tmpl func() *template.Template

func RegisterTemplateFunction(name string, function interface{}) {
	templateFunctions[name] = function
}

func InitTemplates() {
	tmpl = func() *template.Template {
		return template.Must(template.New("base").Funcs(templateFunctions).ParseGlob("templates/*.html"))
	}
	if !arguments.rebuild {
		glog.Info("Caching templates.")
		t := tmpl()
		tmpl = func() *template.Template {
			return t
		}
	}
	glog.Info("Loaded templates.")
}

func ExecuteTemplate(w io.Writer, name string, ctx *RenderContext) error {
	if ctx.template == nil {
		ctx.template = tmpl()
	}
	if ctx.template.Lookup(name) == nil {
		return fmt.Errorf("Template %v not found", name)
	}
	return ctx.template.ExecuteTemplate(w, name, ctx)
}

func init() {
	RegisterTemplateFunction("equal", func(t1, t2 string) bool { return t1 == t2 })
	RegisterTemplateFunction("subtemplate", func(ctx *RenderContext, name string) template.HTML {
		buf := &bytes.Buffer{}
		err := ctx.template.ExecuteTemplate(buf, ctx.Page+"_"+name, ctx)
		if err != nil {
			return template.HTML("")
		}
		return template.HTML(buf.String())
	})

	RegisterReloadFunction(InitTemplates)
}
