package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/gen2brain/go-mpv"
	"github.com/zSnails/peruere/geometry"
	"github.com/zSnails/peruere/xlib"
)

var (
	videoFile string
	geom      string
)

func init() {
	flag.StringVar(&videoFile, "file", "video.mp4", "the file to play as a wallpaper")
	flag.StringVar(&geom, "geometry", "1920x1080+0+0", "the geometry for the background window")
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
	width, height, xOffset, yOffset, err := geometry.ParseGeometry(geom)
	if err != nil {
		log.Fatalln(err)
	}
	window := xlib.XCreateWindow(display, root, xOffset, yOffset, width, height, 0, 0, xlib.InputOutput, nil, xlib.CWOverrideRedirect|xlib.CWBackingStore, &attrs)
	defer xlib.XDestroyWindow(display, window)

	xlib.XSetClassHint(
		display,
		window,
		&xlib.ClassHint{
			ResName:  "peruere",
			ResClass: "peruere",
		},
	)

	windowTypeDesktop := xlib.XInternAtom(display, "_NET_WM_WINDOW_TYPE_DESKTOP", xlib.False)
	windowType := xlib.XInternAtom(display, "_NET_WM_WINDOW_TYPE", xlib.False)
	xlib.XChangeProperty(display, window, windowType, xlib.XA_ATOM, 32, xlib.PropModeReplace, unsafe.Pointer(&windowTypeDesktop), 1)

	motifWmHints := xlib.XInternAtom(display, "_MOTIF_WM_HINTS", xlib.False)
	if motifWmHints != xlib.None {
		prop := [5]int64{2, 0, 0, 0, 0}
		xlib.XChangeProperty(display, window, motifWmHints, motifWmHints, 32, xlib.PropModeReplace, unsafe.Pointer(&prop), 5)
	}

	winLayer := xlib.XInternAtom(display, "_WIN_LAYER", xlib.False)
	if winLayer != xlib.None {
		layerZero := int64(0)
		xlib.XChangeProperty(display, window, winLayer, xlib.XA_CARDINAL, 32, xlib.PropModeAppend, unsafe.Pointer(&layerZero), 1)
	}

	wmState := xlib.XInternAtom(display, "_NET_WM_STATE", xlib.False)
	if wmState != xlib.None {
		stateBelow := xlib.XInternAtom(display, "_NET_WM_STATE_BELOW", xlib.False)
		xlib.XChangeProperty(display, window, wmState, xlib.XA_ATOM, 32, xlib.PropModeAppend, unsafe.Pointer(&stateBelow), 1)
	}

	hints := xlib.WMHints{
		Input: xlib.False,
	}
	xlib.XSetWMProperties(display, window, nil, nil, os.Args, len(os.Args), nil, &hints, nil)

	wmDesktop := xlib.XInternAtom(display, "_NET_WM_DESKTOP", xlib.False)
	stateSticky := xlib.XInternAtom(display, "_NET_WM_STATE_STICKY", xlib.False)
	xlib.XChangeProperty(display, window, wmDesktop, xlib.XA_CARDINAL, 32, xlib.PropModeAppend, unsafe.Pointer(&stateSticky), 1)

	region := xlib.XCreateRegion()
	if region != nil {
		xlib.XShapeCombineRegion(display, window, xlib.ShapeInput, 0, 0, region, xlib.ShapeSet)
		xlib.XDestroyRegion(region)
	}

	xlib.XLowerWindow(display, window)

	m := mpv.New()
	defer m.TerminateDestroy()

	if err := m.SetProperty("wid", mpv.FormatInt64, int(window)); err != nil {
		log.Fatalln(err)
	}

	if err := m.SetPropertyString("loop", "yes"); err != nil {
		log.Fatalln(err)
	}

	if err := m.SetPropertyString("x11-bypass-compositor", "yes"); err != nil {
		log.Fatalln(err)
	}

	if err := m.SetPropertyString("vo", "gpu"); err != nil {
		log.Fatalln(err)
	}

	if err := m.Initialize(); err != nil {
		log.Fatalln(err)
	}

	err = m.RequestLogMessages("trace")
	if err := m.Command([]string{"loadfile", videoFile}); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			event := m.WaitEvent(10000)
			switch event.EventID {
			case mpv.EventLogMsg:
				log.Printf("%v\n", event.LogMessage())
			}
		}
	}()

	xlib.XStoreName(display, window, "peruere")
	xlib.XMapWindow(display, window)
	xlib.XFlush(display)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}
