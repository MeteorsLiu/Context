package context

import (
	ctx "context"
	"sync"
	"time"

	"github.com/v2pro/plz/gls"
)

var (
	ctxIDMap sync.Map
)

type Context struct {
	sync.RWMutex
	parent ctx.Context
	c      ctx.Context
	cel    ctx.CancelFunc
}

func WithCancel(cntx ctx.Context) (ctx.Context, ctx.CancelFunc) {
	c := &Context{
		parent: cntx,
	}
	c.c, c.cel = ctx.WithCancel(cntx)
	return c, c.cancel
}

func (c *Context) cancel() {
	c.Lock()
	defer c.Unlock()
	// check
	ch, ok := ctxIDMap.LoadAndDelete(gls.GoID())
	if ok {
		select {
		case <-ch.(<-chan struct{}):
			return
		default:
		}
	}
	c.cel()
	c.c, c.cel = ctx.WithCancel(c.parent)
}
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	c.RLock()
	defer c.RUnlock()
	return c.c.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	c.RLock()
	defer c.RUnlock()
	ctxIDMap.Store(gls.GoID(), c.c.Done())
	return c.c.Done()
}

func (c *Context) Err() error {
	c.RLock()
	defer c.RUnlock()
	return c.c.Err()
}

func (c *Context) Value(key any) any {
	c.RLock()
	defer c.RUnlock()
	return c.c.Value(key)
}
