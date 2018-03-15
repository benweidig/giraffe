package giraffe

import (
	"bytes"
	"html/template"
	"io"
	"path"
	"sync"

	"github.com/benweidig/giraffe/datasources/fs"
)

// Giraffe is the main struct holding it all together
type Giraffe struct {
	config    Config
	templates map[string]*template.Template
	mutex     sync.RWMutex
}

// New creates a new Giraffe but you have to bring your own config
func New(config Config) *Giraffe {
	g := &Giraffe{
		config:    config,
		templates: make(map[string]*template.Template),
		mutex:     sync.RWMutex{},
	}

	// We add our usual template funcs in New, because they are useful for everyone
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

// PrepareTemplateFuncs adds the "partial" func and the user provided funcs to Giraffe
func (g *Giraffe) PrepareTemplateFuncs() {
	g.config.Funcs["partial"] = func(partial string, partialData interface{}) (template.HTML, error) {
		buf := new(bytes.Buffer)
		name := path.Join("partials", partial)
		err := g.Render(buf, name, partialData, false)
		return template.HTML(buf.String()), err
	}
}

// Render does what it's called, it renders a template with the provided data.
// It supports caching (if not disabled) and the provided funcs.
func (g *Giraffe) Render(out io.Writer, name string, data interface{}, useLayout bool) error {
	// Try getting the template from cache. We do this with a mutex to not do work we don't need
	g.mutex.RLock()
	tpl, ok := g.templates[name]
	g.mutex.RUnlock()

	// Check if found a template or if we shouldn't cache
	if !ok || g.config.DisableCache {

		// We need to add "include" here and not in New(...) because it uses "data"
		g.config.Funcs["include"] = func(layout string) (template.HTML, error) {
			buf := new(bytes.Buffer)
			err := g.Render(buf, layout, data, false)
			return template.HTML(buf.String()), err
		}

		// Load Layout template. At this point we don't support "non-layout"-based rendering.
		layoutStr, err := g.config.Datasource.LoadContent(g.config.Layout)
		if err != nil {
			return err
		}
		layoutTpl, err := template.New(g.config.Layout).Funcs(g.config.Funcs).Parse(layoutStr)
		if err != nil {
			return err
		}
		g.templates[g.config.Layout] = layoutTpl

		// Load the requested template
		templateStr, err := g.config.Datasource.LoadContent(name)
		if err != nil {
			return err
		}

		// Combine both
		tpl, err = layoutTpl.New(name).Funcs(g.config.Funcs).Parse(templateStr)
		if err != nil {
			return err
		}

		// Cache template if necessary
		if g.config.DisableCache == false {
			g.mutex.Lock()
			g.templates[name] = tpl
			g.mutex.Unlock()
		}
	}

	// Check if we need to wrap the template in the layout
	execName := name
	if useLayout {
		execName = g.config.Layout
	}

	// Finally render the template
	err := tpl.ExecuteTemplate(out, execName, data)

	return err
}
