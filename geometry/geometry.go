package geometry

import (
	"strconv"
	"strings"
)

func ParseGeometry(geometry string) (uint, uint, int, int, error) {
	split := strings.FieldsFunc(geometry, func(r rune) bool {
		return r == 'x' || r == '+'
	})

	width, err := strconv.ParseUint(split[0], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	height, err := strconv.ParseUint(split[1], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	xoff, err := strconv.ParseInt(split[2], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	yoff, err := strconv.ParseInt(split[3], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return uint(width), uint(height), int(xoff), int(yoff), nil

}
