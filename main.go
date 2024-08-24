package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/gen2brain/go-mpv"
	"github.com/zSnails/peruere/xlib"
)

var (
	videoFile string
	width     uint
	height    uint
	xOffset   int
	yOffset   int
)

func init() {
	flag.StringVar(&videoFile, "file", "video.mp4", "the file to play as a wallpaper")
	flag.UintVar(&height, "height", 1080, "the height of the window")
	flag.UintVar(&width, "width", 1920, "the width of the window")
	flag.IntVar(&xOffset, "x-offset", 0, "the x axis offset")
	flag.IntVar(&yOffset, "y-offset", 0, "the y axis offset")
	flag.Parse()
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	attrs := xlib.SetWindowAttributes{
		BackgroundPixmap: xlib.ParentRelative,
		BackingStore:     xlib.Always,
		SaveUnder:        xlib.False,
		OverrideRedirect: xlib.True,
	}
	display := xlib.XOpenDisplay(nil)
	defer xlib.XCloseDisplay(display)

	root := xlib.XDefaultRootWindow(display)
	window := xlib.XCreateWindow(display, root, xOffset, yOffset, width, height, 0, 0, xlib.InputOutput, nil, xlib.CWOverrideRedirect|xlib.CWBackingStore, &attrs)
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
	defer m.TerminateDestroy()

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

	if err := m.Command([]string{"loadfile", videoFile}); err != nil {
		panic(err)
	}

	xlib.XStoreName(display, window, "peruere")
	xlib.XMapWindow(display, window)
	xlib.XFlush(display)

	sig := make(chan os.Signal, 0)
	signal.Notify(sig, os.Kill, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	<-sig
}
