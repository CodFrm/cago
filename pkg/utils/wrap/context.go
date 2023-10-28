package wrap

import stdContext "context"

type Context struct {
	stdContext.Context
	name    string
	args    []interface{}
	handler []Handler
	pos     int
	abort   error
}

type Handler func(ctx *Context)

func (c *Context) Name() string {
	return c.name
}

func (c *Context) Args(i int) interface{} {
	return c.args[i]
}

func (c *Context) WithContext(ctx stdContext.Context) *Context {
	c.Context = ctx
	return c
}

func (c *Context) Next() {
	for {
		c.pos++
		if c.pos >= len(c.handler) {
			break
		}
		c.handler[c.pos](c)
	}
}

func (c *Context) Abort(err error) {
	c.abort = err
	if err == nil {
		return
	}
	c.pos = len(c.handler)
}

func (c *Context) IsAbort() error {
	return c.abort
}
