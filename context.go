package context

import (
	ctx "context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/v2pro/plz/gls"
)

var (
	ctxIDMap sync.Map
	timercnt atomic.Int64
)

type Context struct {
	sync.RWMutex
	dur     time.Duration
	timerID int64
	parent  ctx.Context
	timer   *time.Timer
	c       ctx.Context
	cel     ctx.CancelFunc
}

func WithTimeout(cntx ctx.Context, timeout time.Duration) (ctx.Context, ctx.CancelFunc) {
	timerid := timercnt.Add(-1)
	c := &Context{
		parent:  cntx,
		dur:     timeout,
		timerID: timerid,
	}
	c.c, c.cel = ctx.WithCancel(cntx)
	ctxIDMap.Store(timerid, c.c.Done())
	c.timer = time.AfterFunc(timeout, func() {
		c.cancel(c.timerID)
	})
	return c, func() { c.cancel() }
}

func WithCancel(cntx ctx.Context) (ctx.Context, ctx.CancelFunc) {
	c := &Context{
		parent: cntx,
	}
	c.c, c.cel = ctx.WithCancel(cntx)

	return c, func() { c.cancel() }
}

func (c *Context) cancel(id ...int64) {
	c.Lock()
	defer c.Unlock()
	// check
	var gid int64
	if len(id) > 0 {
		gid = id[0]
	} else {
		gid = gls.GoID()
	}
	ch, ok := ctxIDMap.LoadAndDelete(gid)
	if ok {
		select {
		case <-ch.(<-chan struct{}):
			if c.timer != nil {
				c.timer.Stop()
			}
			return
		default:
		}
	}
	c.cel()
	c.c, c.cel = ctx.WithCancel(c.parent)

	if c.timer != nil {
		c.timer.Stop()
		ctxIDMap.Store(c.timerID, c.c.Done())
		c.timer = time.AfterFunc(c.dur, func() {
			c.cancel(c.timerID)
		})
	}
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
