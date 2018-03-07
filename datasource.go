package giraffe

// Datasource is an abstraction of how to load template content
type Datasource interface {
	// LoadContent should load a template by its "convenience name"
	LoadContent(name string) (string, error)
}
