package glycogen

import (
	"fmt"
	"io"
	"os"
)

type CtxError struct {
	issuer *Context
	msg    string
}

func (e *CtxError) Error() string {
	return e.issuer.formatMsg("ERROR " + e.msg)
}

type Context struct {
	name   string
	parent *Context
	out    io.Writer
}

func (c *Context) Child(name string) *Context {
	return &Context{
		name:   name,
		parent: c,
		out:    c.out,
	}
}

func (c *Context) Name() string {
	return c.name
}

func (c *Context) prefix() string {
	if c.parent != nil {
		return c.parent.prefix() + "::" + c.name
	}
	return c.name
}

func (c *Context) formatMsg(msg string) string {
	return fmt.Sprintf("[%s] %s", c.prefix(), msg)
}

func (c *Context) Println(args ...interface{}) {
	fmt.Fprintln(c.out, c.formatMsg(fmt.Sprint(args...)))
}

func (c *Context) Printf(format string, args ...interface{}) {
	fmt.Fprintf(c.out, c.formatMsg(format), args...)
}

func (c *Context) Writeln(s string) {
	fmt.Fprintln(c.out, s)
}

func (c *Context) Errorf(format string, args ...interface{}) error {
	return &CtxError{issuer: c, msg: fmt.Sprintf(format, args...)}
}

func (c *Context) Error(msg string) error {
	return &CtxError{issuer: c, msg: msg}
}

func (c *Context) wrap(err error) error {
	if ce, ok := err.(*CtxError); ok {
		return ce
	}
	return &CtxError{issuer: c, msg: err.Error()}
}

func (c *Context) WrapAndTag(err error, msg string) error {
	if msg == "" {
		return c.wrap(err)
	}
	ce := c.wrap(err).(*CtxError)
	return &CtxError{issuer: ce.issuer, msg: fmt.Sprintf("%s: %s", msg, ce.msg)}
}

func (c *Context) WrapAndTagf(err error, format string, args ...interface{}) error {
	return c.WrapAndTag(err, fmt.Sprintf(format, args...))
}

func NewContext(name string, writer io.Writer) *Context {
	if writer == nil {
		writer = os.Stdout
	}
	return &Context{name: name, out: writer}
}

var DefaultContext = NewContext("default", nil)

type contextExecutionBreak[T any] struct {
	ret T
}

type ContextExecution[T any] struct {
	*Context
}

func (c *ContextExecution[T]) Return(r T) {
	panic(&contextExecutionBreak[T]{ret: r})
}

func ExecuteWithContext[T any](ctx *Context, f func(*ContextExecution[T]) T) (ret T) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if ceb, ok := r.(*contextExecutionBreak[T]); ok {
			ret = ceb.ret
		} else {
			panic(r)
		}
	}()
	return f(&ContextExecution[T]{Context: ctx})
}

func ExecuteWithContextGR[T any](ctx *Context, f func(*ContextExecution[T]) T, cb func(T)) {
	go func() {
		r := ExecuteWithContext(ctx, f)
		cb(r)
	}()
}
