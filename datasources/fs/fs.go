package fs

import (
	"io/ioutil"
	"path/filepath"
)

// Datasource is a simple "load from disk" datasource
type Datasource struct {
	Root      string
	Extension string
}

// Default creates a Datasource with sensible defaults
func Default() *Datasource {
	return &Datasource{
		Root:      "views",
		Extension: ".html",
	}
}

// LoadContent is needed to fulfill giraffe.Datasource interface type
func (d *Datasource) LoadContent(name string) (string, error) {
	path := filepath.Join(d.Root, name+d.Extension)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
