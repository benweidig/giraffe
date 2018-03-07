package giraffe

import (
	"bytes"
	"html/template"
	"io"
	"path"
	"sync"

	"github.com/benweidig/giraffe/datasources/fs"
	"github.com/gin-gonic/gin/render"
)

// Giraffe is a gin html render
type Giraffe struct {
	config    Config
	templates map[string]*template.Template
	mutex     sync.RWMutex
}

// New creates a new Giraffe but bring your own config
func New(config Config) *Giraffe {
	g := &Giraffe{
		config:    config,
		templates: make(map[string]*template.Template),
		mutex:     sync.RWMutex{},
	}

	g.PrepareTemplateFuncs()

	return g
}

// Default creates a new Giraffe with sensible defaults
func Default() *Giraffe {
	config := Config{
		Datasource:   fs.Default(),
		Layout:       "layout",
		Funcs:        make(template.FuncMap),
		DisableCache: false,
	}

	return New(config)
}

// Debug creates a new Giraffe for debugging
func Debug() *Giraffe {
	g := Default()
	g.config.DisableCache = true
	return g
}

// Instance fulfills the gin.render.HTMLRender interface type
func (g *Giraffe) Instance(name string, data interface{}) render.Render {
	return Render{
		Giraffe: g,
		Name:    name,
		Data:    data,
	}
}

// PrepareTemplateFuncs adds the "partial" func and the user provided funcs to Giraffe
func (g *Giraffe) PrepareTemplateFuncs() {
	g.config.Funcs["partial"] = func(partial string, partialData interface{}) (template.HTML, error) {
		buf := new(bytes.Buffer)
		name := path.Join("partials", partial)
		err := g.render(buf, name, partialData, false)
		return template.HTML(buf.String()), err
	}
}

func (g *Giraffe) render(out io.Writer, name string, data interface{}, useLayout bool) error {
	// Try getting the template from cache
	g.mutex.RLock()
	tpl, ok := g.templates[name]
	g.mutex.RUnlock()

	// Check if found or if we shouldn't cache
	if !ok || g.config.DisableCache {

		// We need to add "include" here because it uses "data"
		g.config.Funcs["include"] = func(layout string) (template.HTML, error) {
			buf := new(bytes.Buffer)
			err := g.render(buf, layout, data, false)
			return template.HTML(buf.String()), err
		}

		// Load Layout template
		layoutStr, err := g.config.Datasource.LoadContent(g.config.Layout)
		if err != nil {
			return err
		}
		layoutTpl, err := template.New(g.config.Layout).Funcs(g.config.Funcs).Parse(layoutStr)
		if err != nil {
			return err
		}
		g.templates[g.config.Layout] = layoutTpl

		// Load template
		templateStr, err := g.config.Datasource.LoadContent(name)
		if err != nil {
			return err
		}

		// Combine both
		tpl, err = layoutTpl.New(name).Funcs(g.config.Funcs).Parse(templateStr)
		if err != nil {
			return err
		}

		// Cache template
		g.mutex.Lock()
		g.templates[name] = tpl
		g.mutex.Unlock()
	}

	// Check if we need to wrap the template in the layout
	execName := name
	if useLayout {
		execName = g.config.Layout
	}

	// Render
	err := tpl.ExecuteTemplate(out, execName, data)

	return err
}
