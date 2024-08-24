package xlib

// #cgo LDFLAGS: -lX11
// #include <stdlib.h>
// #include <X11/Xlib.h>
// #include <X11/Xlib.h>
// #include <X11/Xatom.h>
// #include <X11/extensions/shape.h>
// #include "xlib.h"
import "C"

const (
	None           = C.None
	ParentRelative = C.ParentRelative
	CopyFromParent = C.CopyFromParent

	NotUseful  = C.NotUseful
	WhenMapped = C.WhenMapped
	Always     = C.Always

	False = C.False
	True  = C.True

	InputOutput = C.InputOutput

	CWOverrideRedirect = C.CWOverrideRedirect
	CWBackingStore     = C.CWBackingStore

	XA_ATOM         = C.XA_ATOM
	XA_CARDINAL     = C.XA_CARDINAL
	XA_STRING       = C.XA_STRING
	PropModeReplace = C.PropModeReplace
	PropModeAppend  = C.PropModeAppend

	ShapeInput = C.ShapeInput
	ShapeSet   = C.ShapeSet

	Success           = C.Success
	BadRequest        = C.BadRequest
	BadValue          = C.BadValue
	BadWindow         = C.BadWindow
	BadPixmap         = C.BadPixmap
	BadAtom           = C.BadAtom
	BadCursor         = C.BadCursor
	BadFont           = C.BadFont
	BadMatch          = C.BadMatch
	BadDrawable       = C.BadDrawable
	BadAccess         = C.BadAccess
	BadAlloc          = C.BadAlloc
	BadColor          = C.BadColor
	BadGC             = C.BadGC
	BadIDChoice       = C.BadIDChoice
	BadName           = C.BadName
	BadLength         = C.BadLength
	BadImplementation = C.BadImplementation
)

var errmap = map[int]string{
	Success:           "Success",
	BadRequest:        "BadRequest",
	BadValue:          "BadValue",
	BadWindow:         "BadWindow",
	BadPixmap:         "BadPixmap",
	BadAtom:           "BadAtom",
	BadCursor:         "BadCursor",
	BadFont:           "BadFont",
	BadMatch:          "BadMatch",
	BadDrawable:       "BadDrawable",
	BadAccess:         "BadAccess",
	BadAlloc:          "BadAlloc",
	BadColor:          "BadColor",
	BadGC:             "BadGC",
	BadIDChoice:       "BadIDChoice",
	BadName:           "BadName",
	BadLength:         "BadLength",
	BadImplementation: "BadImplementation",
}

func ErrorString(err int) string {
	str, ok := errmap[err]
	if !ok {
		return "Unknown Error"
	}
	return str
}
