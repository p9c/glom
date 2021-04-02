package gel

import (
	"github.com/p9c/gel/fonts/p9fonts"
	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/meta"
	"math"
	"time"
	
	"gioui.org/io/event"
	
	"github.com/p9c/qu"
	
	"gioui.org/app"
	"gioui.org/io/system"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

type CallbackQueue chan func() error

func NewCallbackQueue(bufSize int) CallbackQueue {
	return make(CallbackQueue, bufSize)
}

type scaledConfig struct {
	Scale float32
}

func (s *scaledConfig) Now() time.Time {
	return time.Now()
}

func (s *scaledConfig) Px(v unit.Value) int {
	scale := s.Scale
	if v.U == unit.UnitPx {
		scale = 1
	}
	return int(math.Round(float64(scale * v.V)))
}

type Window struct {
	*Theme
	*app.Window
	opts    []app.Option
	scale   *scaledConfig
	Width   int // stores the width at the beginning of render
	Height  int
	ops     op.Ops
	evQ     system.FrameEvent
	Runner  CallbackQueue
	overlay []func(gtx l.Context)
}

func (w *Window) PushOverlay(overlay func(gtx l.Context)) {
	w.overlay = append(w.overlay, overlay)
}

func (w *Window) PopOverlay() {
	if len(w.overlay) == 0 {
		return
	}
	w.overlay = w.overlay[:len(w.overlay)-1]
}

func (w *Window) Overlay(gtx l.Context) {
	for _, overlay := range w.overlay {
		overlay(gtx)
	}
}

// NewWindowP9 creates a new window
func NewWindowP9(quit chan struct{}) (out *Window) {
	out = &Window{
		scale:  &scaledConfig{1},
		Runner: NewCallbackQueue(32),
	}
	out.Theme = NewTheme(
		binary.New(meta.Data{}, false, out.Theme.SetDarkTheme),
		p9fonts.Collection(), quit,
	)
	out.Theme.WidgetPool = out.NewPool()
	return
}

// NewWindow creates a new window
func NewWindow(th *Theme) (out *Window) {
	out = &Window{
		Theme: th,
		scale: &scaledConfig{1},
	}
	return
}

// Title sets the title of the window
func (w *Window) Title(title string) (out *Window) {
	w.opts = append(w.opts, app.Title(title))
	return w
}

// Size sets the dimensions of the window
func (w *Window) Size(width, height float32) (out *Window) {
	w.opts = append(
		w.opts,
		app.Size(w.Theme.TextSize.Scale(width), w.Theme.TextSize.Scale(height)),
	)
	return w
}

// Scale sets the scale factor for rendering
func (w *Window) Scale(s float32) *Window {
	w.scale = &scaledConfig{s}
	return w
}

// Open sets the window options and initialise the app.window
func (w *Window) Open() (out *Window) {
	if w.scale == nil {
		w.Scale(1)
	}
	if w.opts != nil {
		w.Window = app.NewWindow(w.opts...)
		w.opts = nil
	}
	return w
}

func (w *Window) Run(
	frame func(ctx l.Context) l.Dimensions,
	overlay func(ctx l.Context), destroy func(), quit qu.C,
) (e error) {
	for {
		select {
		case fn := <-w.Runner:
			if e = fn(); E.Chk(e) {
				return
			}
		case <-quit.Wait():
			return
			// by repeating selectors we decrease the chance of a runner delaying
			// a frame event hitting the physical frame deadline
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		case ev := <-w.Window.Events():
			if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
				return
			}
		}
	}
}

func (w *Window) processEvents(e event.Event, frame func(ctx l.Context) l.Dimensions, destroy func()) error {
	switch e := e.(type) {
	case system.DestroyEvent:
		D.Ln("received destroy event", e.Err)
		// if e.Err != nil {
		// 	if strings.Contains(e.Err.Error(), "eglCreateWindowSurface failed") {
		// 		return nil
		// 	}
		// }
		destroy()
		return e.Err
	case system.FrameEvent:
		ops := op.Ops{}
		c := l.NewContext(&ops, e)
		// update dimensions for responsive sizing widgets
		w.Width = c.Constraints.Max.X
		w.Height = c.Constraints.Max.Y
		frame(c)
		w.Overlay(c)
		e.Frame(c.Ops)
	}
	return nil
}
