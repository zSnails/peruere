package geometry

import (
	"fmt"
	"testing"
)

func TestParseGeometry(t *testing.T) {
	x, y, xoff, yoff, err := ParseGeometry("1920x1080+10+20")
    if err != nil {
        t.Fatal(err)
    }
	fmt.Printf("x: %v\n", x)
	fmt.Printf("y: %v\n", y)
	fmt.Printf("xoff: %v\n", xoff)
	fmt.Printf("yoff: %v\n", yoff)

	if x != 1920 {
		t.Fatal("x != 1920")
	}

	if y != 1080 {
		t.Fatal("y != 1080")
	}

	if xoff != 10 {
		t.Fatal("xoff != 0")
	}

	if yoff != 20 {
		t.Fatal("yoff != 0")
	}
}
