package injection

import "reflect"

// NewContext return an instance of Context
func NewContext() *Context {
	return &Context{container: make(map[string]interface{})}
}

// Context of dependencies to be injected as context
type Context struct {
	container map[string]interface{}
}

// AddDependency adds a dependency to the container
func (c *Context) AddDependency(key string, d interface{}) {
	c.container[key] = d
}

// Dependency returns an object given a key
func (c *Context) Dependency(key string) interface{} {
	return c.container[key]
}

// Inject add the Context instance as dependency
func (c *Context) Inject(o interface{}) {
	rv := reflect.ValueOf(c)
	reflectedContext := reflect.ValueOf(o).Elem()
	typeOfT := rv.Type()

	for i := 0; i < reflectedContext.NumField(); i++ {
		f := reflectedContext.Field(i)
		if f.Type().String() == typeOfT.String() {
			reflectedContext.Field(i).Set(rv)
		}
	}
}
