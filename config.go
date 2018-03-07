package giraffe

import (
	"html/template"
)

// Config is the config of the template engine
type Config struct {
	Datasource   Datasource       // The datasource which loads the template content
	Layout       string           // Filename / path to layout (without extension)
	Funcs        template.FuncMap // Template functions
	DisableCache bool             // Disables caching
}
