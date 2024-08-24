//          Copyright 2016 Vitali Baumtrok
// Distributed under the Boost Software License, Version 1.0.
//      (See accompanying file LICENSE or copy at
//        http://www.boost.org/LICENSE_1_0.txt)

// Binding of Xlib (version 11, release 6.7).
package xlib

// #cgo LDFLAGS: -lX11 -lXext
// #include <stdlib.h>
// #include <X11/Xlib.h>
// #include <X11/Xutil.h>
// #include <X11/extensions/shape.h>
// #include "xlib.h"
import "C"
import (
	"unsafe"
)

type Display C.Display
type Screen C.Screen
type Window C.Window
type Atom C.Atom
type XWMHints C.XWMHints
type Bool C.Bool
type WindowAttributes C.XWindowAttributes
type Cursor C.Cursor
type XTextProperty C.XTextProperty
type Visual C.Visual
type XSetWindowAttributes C.XSetWindowAttributes
type Region C.Region
type XSizeHints C.XSizeHints
type XClassHint C.XClassHint

type SetWindowAttributes struct {
	BackgroundPixmap   uint64
	BackgroundPixel    uint64
	BorderPixmap       uint64
	BorderPixel        uint64
	BitGravity         int32
	WinGravity         int32
	BackingStore       int32
	BackingPlanes      uint64
	BackingPixel       uint64
	SaveUnder          uint8
	EventMask          int64
	DoNotPropagateMask int64
	OverrideRedirect   uint8
	Colormap           uint64
	Cursor             uint64
}

type WMHints struct {
	Flags        int64
	Input        int
	InitialState int
	IconPixmap   uint64
	IconWindow   uint64
	IconX        int
	IconY        int
	IconMask     uint64
	WindowGroup  uint64
}

func XInternAtom(display *Display, atom string, state int) Atom {
	displayC := (*C.Display)(display)
	atomC := C.CString(atom)
	defer C.free(unsafe.Pointer(atomC))
	stateC := C.int(state)
	_atom := C.XInternAtom(displayC, atomC, stateC)
	return Atom(_atom)

}

func stringSliceToCArray(strs []string) **C.char {
	// Allocate an array of char pointers
	cArray := C.malloc(C.size_t(len(strs)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	// Convert the allocated memory to a pointer to a pointer of C.char
	cStrings := (*[1<<30 - 1]*C.char)(cArray)

	for i, s := range strs {
		cStr := C.CString(s) // Convert Go string to C string
		cStrings[i] = cStr   // Assign the C string to the array
	}

	return (**C.char)(cArray)
}

func XCreateRegion() Region {
	return (Region)(C.XCreateRegion())
}

func XShapeCombineRegion(display *Display, window Window, destKind, xOff, yOff int, region Region, op int) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	destKindC := C.int(destKind)
	xOffC := C.int(xOff)
	yOffC := C.int(yOff)
	regionC := (C.Region)(region)
	opC := C.int(op)
	C.XShapeCombineRegion(displayC, windowC, destKindC, xOffC, yOffC, regionC, opC)
}

func XDestroyRegion(region Region) {
	C.XDestroyRegion(region)
}

func XSetWMProperties(display *Display, window Window, windowName, iconName *XTextProperty, argv []string, argc int, normalHints *XSizeHints, hints *WMHints, classHint *XClassHint) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	windowNameC := (*C.XTextProperty)(windowName)
	iconNameC := (*C.XTextProperty)(iconName)
	argvC := stringSliceToCArray(argv)
	argcC := C.int(argc)
	normalHintsC := (*C.XSizeHints)(normalHints)
	hintsC := (*C.XWMHints)(makeWMHints(hints))
	classHintC := (*C.XClassHint)(classHint)
	C.XSetWMProperties(displayC, windowC, windowNameC, iconNameC, argvC, argcC, normalHintsC, hintsC, classHintC)
}

func XChangeProperty(display *Display, window Window, property, _type Atom, format, mode int, data unsafe.Pointer, nElements int) int {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	propertyC := (C.Atom)(property)
	_typeC := (C.Atom)(_type)
	formatC := C.int(format)
	modeC := C.int(mode)
	dataC := (*C.uchar)(data)
	nElementsC := C.int(nElements)

	return (int)(C.XChangeProperty(displayC, windowC, propertyC, _typeC, formatC, modeC, dataC, nElementsC))
}

func strConcat(a []interface{}) string {
	str := ""
	for _, strPart := range a {
		switch s := strPart.(type) {
		case string:
			str += s
		}
	}
	return str
}

func XOpenDisplay(displayNameParts ...interface{}) *Display {
	if len(displayNameParts) == 0 {
		display := C.XOpenDisplay(nil)
		return (*Display)(display)

	} else {
		displayNameComplete := strConcat(displayNameParts)
		if len(displayNameComplete) > 0 {
			displayNameCompleteC := C.CString(displayNameComplete)
			display := C.XOpenDisplay(displayNameCompleteC)
			C.free(unsafe.Pointer(displayNameCompleteC))
			return (*Display)(display)

		} else {
			display := C.XOpenDisplay(nil)
			return (*Display)(display)
		}
	}
}

func XCloseDisplay(display *Display) {
	displayC := (*C.Display)(display)
	C.XCloseDisplay(displayC)
}

func XDisplayString(display *Display) string {
	displayC := (*C.Display)(display)
	displayNameC := C.XDisplayString(displayC)
	displayName := C.GoString(displayNameC)
	C.free(unsafe.Pointer(displayNameC))
	return displayName
}

func XScreenCount(display *Display) int {
	displayC := (*C.Display)(display)
	screenCount := C.XScreenCount(displayC)
	return int(screenCount)
}

func XScreenOfDisplay(display *Display, screenNumber int) *Screen {
	displayC := (*C.Display)(display)
	screen := C.XScreenOfDisplay(displayC, C.int(screenNumber))
	return (*Screen)(screen)
}

func XWidthOfScreen(screen *Screen) int {
	screenC := (*C.Screen)(screen)
	width := C.XWidthOfScreen(screenC)
	return int(width)
}

func XHeightOfScreen(screen *Screen) int {
	screenC := (*C.Screen)(screen)
	height := C.XHeightOfScreen(screenC)
	return int(height)
}

func XDefaultScreenOfDisplay(display *Display) *Screen {
	displayC := (*C.Display)(display)
	defaultScreen := C.XDefaultScreenOfDisplay(displayC)
	return (*Screen)(defaultScreen)
}

func XRootWindowOfScreen(screen *Screen) Window {
	screenC := (*C.Screen)(screen)
	rootWindow := C.XRootWindowOfScreen(screenC)
	return Window(rootWindow)
}

func XCreateWindow(display *Display, parent Window, x, y int, width, height uint, borderWidth, depth int, class uint, visual *Visual, valueMask uint64, attributes *SetWindowAttributes) Window {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(parent)
	xC := C.int(x)
	yC := C.int(y)
	widthC := C.uint(width)
	heightC := C.uint(height)
	borderWidthC := C.uint(borderWidth)
	depthC := C.int(depth)
	classC := C.uint(class)
	visualC := (*C.Visual)(visual)
	valueMaskC := C.ulong(valueMask)
	attributesC := (*C.XSetWindowAttributes)(makeAttrs(attributes))
	window := C.XCreateWindow(displayC, windowC, xC, yC, widthC, heightC, borderWidthC, depthC, classC, visualC, valueMaskC, attributesC)
	return Window(window)
}

func makeWMHints(hints *WMHints) *XWMHints {
	return &XWMHints{
		flags:         C.long(hints.Flags),
		input:         C.int(hints.Input),
		initial_state: C.int(hints.InitialState),
		icon_pixmap:   C.ulong(hints.IconPixmap),
		icon_window:   C.ulong(hints.IconWindow),
		icon_x:        C.int(hints.IconX),
		icon_y:        C.int(hints.IconY),
		icon_mask:     C.ulong(hints.IconMask),
		window_group:  C.ulong(hints.WindowGroup),
	}
}

func makeAttrs(attrs *SetWindowAttributes) *XSetWindowAttributes {
	return &XSetWindowAttributes{
		background_pixmap:     C.ulong(attrs.BackgroundPixmap),
		background_pixel:      C.ulong(attrs.BackgroundPixel),
		border_pixmap:         C.ulong(attrs.BorderPixmap),
		border_pixel:          C.ulong(attrs.BorderPixel),
		bit_gravity:           C.int(attrs.BitGravity),
		win_gravity:           C.int(attrs.WinGravity),
		backing_store:         C.int(attrs.BackingStore),
		backing_planes:        C.ulong(attrs.BackingPlanes),
		backing_pixel:         C.ulong(attrs.BackingPixel),
		save_under:            C.int(attrs.SaveUnder),
		event_mask:            C.long(attrs.EventMask),
		do_not_propagate_mask: C.long(attrs.DoNotPropagateMask),
		override_redirect:     C.int(attrs.OverrideRedirect),
		colormap:              C.ulong(attrs.Colormap),
		cursor:                C.ulong(attrs.Cursor),
	}
}

func XLowerWindow(display *Display, window Window) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	C.XLowerWindow(displayC, windowC)
}

func XCreateSimpleWindow(display *Display, parent Window, x, y int, width, height, borderWidth uint, border, background uint64) Window {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(parent)
	xC := C.int(x)
	yC := C.int(y)
	widthC := C.uint(width)
	heightC := C.uint(height)
	borderWidthC := C.uint(borderWidth)
	borderC := C.ulong(border)
	backgroundC := C.ulong(background)
	window := C.XCreateSimpleWindow(displayC, windowC, xC, yC, widthC, heightC, borderWidthC, borderC, backgroundC)
	return Window(window)
}

func XMapWindow(display *Display, window Window) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	C.XMapWindow(displayC, windowC)
}

func XSelectInput(display *Display, window Window, eventMask int64) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	eventMaskC := (C.long)(eventMask)
	C.XSelectInput(displayC, windowC, eventMaskC)
}

func XGrabKey(display *Display, keycode int, modifiers int, grab_window Window, owner_events Bool, pointer_mode int, keyboard_mode int) int {
	displayC := (*C.Display)(display)
	keycodeC := (C.int)(keycode)
	modifiersC := (C.uint)(modifiers)
	grab_windowC := (C.Window)(grab_window)
	owner_eventsC := (C.Bool)(owner_events)
	pointer_modeC := (C.int)(pointer_mode)
	keyboard_modeC := (C.int)(keyboard_mode)
	status := C.XGrabKey(displayC, keycodeC, modifiersC, grab_windowC, owner_eventsC, pointer_modeC, keyboard_modeC)
	return int(status)

}

func XRaiseWindow(display *Display, window Window) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	C.XRaiseWindow(displayC, windowC)
}

func XGetWindowAttributes(display *Display, window Window) *C.XWindowAttributes {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	windowAttributes := new(C.XWindowAttributes)
	C.XGetWindowAttributes(displayC, windowC, windowAttributes)
	return windowAttributes
}

func XMoveResizeWindow(display *Display, window Window, x, y int, width, height uint) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	xC := (C.int)(x)
	yC := (C.int)(y)
	widthC := (C.uint)(width)
	heightC := (C.uint)(height)
	C.XMoveResizeWindow(displayC, windowC, xC, yC, widthC, heightC)
}

func XKeysymToKeycode(display *Display, keysym uint64) int {
	displayC := (*C.Display)(display)
	keysymC := (C.KeySym)(keysym)
	keycode := C.XKeysymToKeycode(displayC, keysymC)
	return int(keycode)
}

func XStringToKeysym(string string) uint64 {
	stringC := C.CString(string)
	keysym := C.XStringToKeysym(stringC)
	C.free(unsafe.Pointer(stringC))
	return uint64(keysym)
}

func XGrabButton(display *Display, button int, modifiers int, grab_window Window, owner_events Bool, event_mask int64, pointer_mode int, keyboard_mode int, confine_to Window, cursor Cursor) int {
	displayC := (*C.Display)(display)
	buttonC := (C.uint)(button)
	modifiersC := (C.uint)(modifiers)
	grab_windowC := (C.Window)(grab_window)
	owner_eventsC := (C.Bool)(owner_events)
	event_maskC := (C.uint)(event_mask)
	pointer_modeC := (C.int)(pointer_mode)
	keyboard_modeC := (C.int)(keyboard_mode)
	confine_toC := (C.Window)(confine_to)
	cursorC := (C.Cursor)(cursor)
	status := C.XGrabButton(displayC, buttonC, modifiersC, grab_windowC, owner_eventsC, event_maskC, pointer_modeC, keyboard_modeC, confine_toC, cursorC)
	return int(status)
}

func XDefaultRootWindow(display *Display) Window {
	displayC := (*C.Display)(display)
	rootWindow := C.XDefaultRootWindow(displayC)
	return Window(rootWindow)
}

func XUngrabKey(display *Display, keycode int, modifiers int, grab_window Window) {
	displayC := (*C.Display)(display)
	keycodeC := (C.int)(keycode)
	modifiersC := (C.uint)(modifiers)
	grab_windowC := (C.Window)(grab_window)
	C.XUngrabKey(displayC, keycodeC, modifiersC, grab_windowC)
}

func XUngrabButton(display *Display, button int, modifiers int, grab_window Window) {
	displayC := (*C.Display)(display)
	buttonC := (C.uint)(button)
	modifiersC := (C.uint)(modifiers)
	grab_windowC := (C.Window)(grab_window)
	C.XUngrabButton(displayC, buttonC, modifiersC, grab_windowC)
}

func XGrabPointer(display *Display, grab_window Window, owner_events Bool, event_mask int64, pointer_mode int, keyboard_mode int, confine_to Window, cursor Cursor, time uint64) int {
	displayC := (*C.Display)(display)
	grab_windowC := (C.Window)(grab_window)
	owner_eventsC := (C.Bool)(owner_events)
	event_maskC := (C.uint)(event_mask)
	pointer_modeC := (C.int)(pointer_mode)
	keyboard_modeC := (C.int)(keyboard_mode)
	confine_toC := (C.Window)(confine_to)
	cursorC := (C.Cursor)(cursor)
	timeC := (C.Time)(time)
	status := C.XGrabPointer(displayC, grab_windowC, owner_eventsC, event_maskC, pointer_modeC, keyboard_modeC, confine_toC, cursorC, timeC)
	return int(status)
}

func XUngrabPointer(display *Display, time uint64) {
	displayC := (*C.Display)(display)
	timeC := (C.Time)(time)
	C.XUngrabPointer(displayC, timeC)
}

func XWarpPointer(display *Display, src_window Window, dest_window Window, src_x, src_y int, src_width, src_height uint, dest_x, dest_y int) {
	displayC := (*C.Display)(display)
	src_windowC := (C.Window)(src_window)
	dest_windowC := (C.Window)(dest_window)
	src_xC := (C.int)(src_x)
	src_yC := (C.int)(src_y)
	src_widthC := (C.uint)(src_width)
	src_heightC := (C.uint)(src_height)
	dest_xC := (C.int)(dest_x)
	dest_yC := (C.int)(dest_y)
	C.XWarpPointer(displayC, src_windowC, dest_windowC, src_xC, src_yC, src_widthC, src_heightC, dest_xC, dest_yC)
}

func XQueryPointer(display *Display, window Window) (root_return Window, child_return Window, root_x_return, root_y_return, win_x_return, win_y_return int, mask_return uint) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	var root_returnC C.Window
	var child_returnC C.Window
	var root_x_returnC C.int
	var root_y_returnC C.int
	var win_x_returnC C.int
	var win_y_returnC C.int
	var mask_returnC C.uint
	C.XQueryPointer(displayC, windowC, &root_returnC, &child_returnC, &root_x_returnC, &root_y_returnC, &win_x_returnC, &win_y_returnC, &mask_returnC)
	root_return = Window(root_returnC)
	child_return = Window(child_returnC)
	root_x_return = int(root_x_returnC)
	root_y_return = int(root_y_returnC)
	win_x_return = int(win_x_returnC)
	win_y_return = int(win_y_returnC)
	mask_return = uint(mask_returnC)
	return
}

func XUnmapWindow(display *Display, window Window) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	C.XUnmapWindow(displayC, windowC)
}

func XDestroyWindow(display *Display, window Window) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	C.XDestroyWindow(displayC, windowC)
}

func XFree(data unsafe.Pointer) {
	C.XFree(data)
}

func XFlush(display *Display) {
	displayC := (*C.Display)(display)
	C.XFlush(displayC)
}

func XStoreName(display *Display, window Window, name string) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	nameC := C.CString(name)
	C.XStoreName(displayC, windowC, nameC)
	C.free(unsafe.Pointer(nameC))
}

func XCreateFontCursor(display *Display, shape uint) Cursor {
	displayC := (*C.Display)(display)
	shapeC := (C.uint)(shape)
	cursor := C.XCreateFontCursor(displayC, shapeC)
	return Cursor(cursor)
}

func XDefineCursor(display *Display, window Window, cursor Cursor) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	cursorC := (C.Cursor)(cursor)
	C.XDefineCursor(displayC, windowC, cursorC)
}

func XUndefineCursor(display *Display, window Window) {
	displayC := (*C.Display)(display)
	windowC := (C.Window)(window)
	C.XUndefineCursor(displayC, windowC)
}

func XCreateGC(display *Display, drawable Window, mask uint64, values *C.XGCValues) C.GC {
	displayC := (*C.Display)(display)
	drawableC := (C.Drawable)(drawable)
	gc := C.XCreateGC(displayC, drawableC, C.ulong(mask), values)
	return gc
}

func XFreeGC(display *Display, gc C.GC) {
	displayC := (*C.Display)(display)
	C.XFreeGC(displayC, gc)
}

func XSetForeground(display *Display, gc C.GC, foreground uint64) {
	displayC := (*C.Display)(display)
	foregroundC := (C.ulong)(foreground)
	C.XSetForeground(displayC, gc, foregroundC)
}

func XSetBackground(display *Display, gc C.GC, background uint64) {
	displayC := (*C.Display)(display)
	backgroundC := (C.ulong)(background)
	C.XSetBackground(displayC, gc, backgroundC)
}

func XSetLineAttributes(display *Display, gc C.GC, line_width uint, line_style int, cap_style int, join_style int) {
	displayC := (*C.Display)(display)
	line_widthC := (C.uint)(line_width)
	line_styleC := (C.int)(line_style)
	cap_styleC := (C.int)(cap_style)
	join_styleC := (C.int)(join_style)
	C.XSetLineAttributes(displayC, gc, line_widthC, line_styleC, cap_styleC, join_styleC)
}

func XDrawLine(display *Display, drawable Window, gc C.GC, x1, y1, x2, y2 int) {
	displayC := (*C.Display)(display)
	drawableC := (C.Drawable)(drawable)
	x1C := (C.int)(x1)
	y1C := (C.int)(y1)
	x2C := (C.int)(x2)
	y2C := (C.int)(y2)
	C.XDrawLine(displayC, drawableC, gc, x1C, y1C, x2C, y2C)
}

func XDrawRectangle(display *Display, drawable Window, gc C.GC, x, y int, width, height uint) {
	displayC := (*C.Display)(display)
	drawableC := (C.Drawable)(drawable)
	xC := (C.int)(x)
	yC := (C.int)(y)
	widthC := (C.uint)(width)
	heightC := (C.uint)(height)
	C.XDrawRectangle(displayC, drawableC, gc, xC, yC, widthC, heightC)
}

func XFillRectangle(display *Display, drawable Window, gc C.GC, x, y int, width, height uint) {
	displayC := (*C.Display)(display)
	drawableC := (C.Drawable)(drawable)
	xC := (C.int)(x)
	yC := (C.int)(y)
	widthC := (C.uint)(width)
	heightC := (C.uint)(height)
	C.XFillRectangle(displayC, drawableC, gc, xC, yC, widthC, heightC)
}

func XQlength(display *Display) int {
	displayC := (*C.Display)(display)
	length := C.XQLength(displayC)
	return int(length)
}

func XRootWindow(display *Display, screenNumber int) Window {
	displayC := (*C.Display)(display)
	screenNumberC := (C.int)(screenNumber)
	rootWindow := C.XRootWindow(displayC, screenNumberC)
	return Window(rootWindow)
}

func XServerVendor(display *Display) string {
	displayC := (*C.Display)(display)
	serverVendorC := C.XServerVendor(displayC)
	serverVendor := C.GoString(serverVendorC)
	C.free(unsafe.Pointer(serverVendorC))
	return serverVendor
}

func XVendorRelease(display *Display) int {
	displayC := (*C.Display)(display)
	vendorRelease := C.XVendorRelease(displayC)
	return int(vendorRelease)
}

func XDefaultScreen(display *Display) int {
	displayC := (*C.Display)(display)
	screenNumber := C.XDefaultScreen(displayC)
	return int(screenNumber)
}

func XDefaultDepth(display *Display, screenNumber int) int {
	displayC := (*C.Display)(display)
	screenNumberC := (C.int)(screenNumber)
	depth := C.XDefaultDepth(displayC, screenNumberC)
	return int(depth)
}
