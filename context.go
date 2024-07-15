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
	return e.issuer.formatMsg(e.msg)
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
	}
}

func (c *Context) Name() string {
	return c.name
}

func (l *Context) prefix() string {
	if l.parent != nil {
		return l.parent.prefix() + "::" + l.name
	}
	return l.name
}

func (l *Context) formatMsg(msg string) string {
	return fmt.Sprintf("[%s] %s", l.prefix(), msg)
}

func (l *Context) Println(args ...interface{}) {
	fmt.Fprintln(l.out, l.formatMsg(fmt.Sprint(args...)))
}

func (l *Context) Printf(format string, args ...interface{}) {
	fmt.Fprintf(l.out, l.formatMsg(format), args...)
}

func (l *Context) Errorf(format string, args ...interface{}) error {
	return &CtxError{issuer: l, msg: fmt.Sprintf(format, args...)}
}

func (l *Context) Error(msg string) error {
	return &CtxError{issuer: l, msg: msg}
}

func (l *Context) wrap(err error) error {
	if ce, ok := err.(*CtxError); ok {
		return ce
	}
	return &CtxError{issuer: l, msg: err.Error()}
}

func (l *Context) WrapAndTag(err error, msg string) error {
	if msg == "" {
		return l.wrap(err)
	}
	ce := l.wrap(err).(*CtxError)
	return &CtxError{issuer: ce.issuer, msg: fmt.Sprintf("%s: %s", msg, ce.msg)}
}

func (l *Context) WrapAndTagf(err error, format string, args ...interface{}) error {
	return l.WrapAndTag(err, fmt.Sprintf(format, args...))
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
	Ctx *Context
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
	return f(&ContextExecution[T]{Ctx: ctx})
}
