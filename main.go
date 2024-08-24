package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/gen2brain/go-mpv"
	"github.com/zSnails/peruere/xlib"
)

func main() {

	attrs := xlib.SetWindowAttributes{
		BackgroundPixmap: xlib.ParentRelative,
		BackingStore:     xlib.Always,
		SaveUnder:        xlib.False,
		OverrideRedirect: xlib.True,
	}
	display := xlib.XOpenDisplay(nil)
	defer xlib.XCloseDisplay(display)

	root := xlib.XDefaultRootWindow(display)
	window := xlib.XCreateWindow(display, root, 0, 0, 1920, 1080, 0, 0, xlib.InputOutput, nil, xlib.CWOverrideRedirect|xlib.CWBackingStore, &attrs)
	fmt.Printf("window: %v\n", window)
	defer xlib.XDestroyWindow(display, window)

	{
		prop := xlib.XInternAtom(display, "_NET_WM_WINDOW_TYPE_DESKTOP", xlib.False)
		xa := xlib.XInternAtom(display, "_NET_WM_WINDOW_TYPE", xlib.False)
		xlib.XChangeProperty(display, window, xa, xlib.XA_ATOM, 32, xlib.PropModeReplace, unsafe.Pointer(&prop), 1)
	}

	{
		xa := xlib.XInternAtom(display, "_MOTIF_WM_HINTS", xlib.False)
		if xa != xlib.None {
			prop := [5]int64{2, 0, 0, 0, 0}
			xlib.XChangeProperty(display, window, xa, xa, 32, xlib.PropModeReplace, unsafe.Pointer(&prop), 5)
		}
	}

	{
		xa := xlib.XInternAtom(display, "_WIN_LAYER", xlib.False)
		if xa != xlib.None {
			prop := int64(0)
			xlib.XChangeProperty(display, window, xa, xlib.XA_CARDINAL, 32, xlib.PropModeAppend, unsafe.Pointer(&prop), 1)
		}

		xa = xlib.XInternAtom(display, "_NET_WM_STATE", xlib.False)
		if xa != xlib.None {
			xa_prop := xlib.XInternAtom(display, "_NET_WM_STATE_BELOW", xlib.False)
			xlib.XChangeProperty(display, window, xa, xlib.XA_ATOM, 32, xlib.PropModeAppend, unsafe.Pointer(&xa_prop), 1)
		}
	}

	{
		hints := xlib.WMHints{
			Input: xlib.False,
		}
		xlib.XSetWMProperties(display, window, nil, nil, os.Args, len(os.Args), nil, &hints, nil)
	}

	{
		xa := xlib.XInternAtom(display, "_NET_WM_DESKTOP", xlib.False)
		xa_xprop := xlib.XInternAtom(display, "_NET_WM_STATE_STICKY", xlib.False)
		xlib.XChangeProperty(display, window, xa, xlib.XA_CARDINAL, 32, xlib.PropModeAppend, unsafe.Pointer(&xa_xprop), 1)
	}

	{
		region := xlib.XCreateRegion()
		if region != nil {
			xlib.XShapeCombineRegion(display, window, xlib.ShapeInput, 0, 0, region, xlib.ShapeSet)
			xlib.XDestroyRegion(region)
		}
	}

	xlib.XLowerWindow(display, window)

	m := mpv.New()

	if err := m.SetProperty("wid", mpv.FormatInt64, int(window)); err != nil {
		panic(err)
	}

	if err := m.SetPropertyString("loop", "yes"); err != nil {
		panic(err)
	}

	if err := m.SetPropertyString("x11-bypass-compositor", "yes"); err != nil {
		panic(err)
	}

	if err := m.SetPropertyString("vo", "gpu"); err != nil {
		panic(err)
	}

	if err := m.Initialize(); err != nil {
		panic(err)
	}

	if err := m.Command([]string{"loadfile", "/home/ayaka/Genshin Impact - Shenhe [ Live Wallpaper ] [rOEyER6lXWo].webm"}); err != nil {
		panic(err)
	}

	xlib.XMapWindow(display, window)
	xlib.XStoreName(display, window, "peruere")
	xlib.XFlush(display)

loop:
	for {
		event := m.WaitEvent(10000)
		switch event.EventID {
		case mpv.EventShutdown:
			break loop
		}
	}
}
