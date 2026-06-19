// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	svg := fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az, aOK := corner(i+1, j)
			bx, by, bz, bOK := corner(i, j)
			cx, cy, cz, cOK := corner(i, j+1)
			dx, dy, dz, dOK := corner(i+1, j+1)

			// Only render the polygon if all four corners are valid, finite numbers
			if aOK && bOK && cOK && dOK {
				zavg := (az + bz + cz + dz) / 4.0
				// Map average height to a premium HSL color gradient
				fillColor := color(zavg, -0.2172, 1.0)
				svg += fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g' style='fill:%s; stroke:grey; stroke-opacity:0.3; stroke-width:0.5'/>\n",
					ax, ay, bx, by, cx, cy, dx, dy, fillColor)
			}
		}
	}
	svg += "</svg>\n"

	// Print to standard output
	fmt.Print(svg)

	// Save to surface.svg
	err := os.WriteFile("surface.svg", []byte(svg), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing surface.svg: %v\n", err)
	}
}

func corner(i, j int) (float64, float64, float64, bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)
	if math.IsNaN(z) || math.IsInf(z, 0) {
		return 0, 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	// Ensure the projected coordinates are valid finite numbers
	if math.IsNaN(sx) || math.IsInf(sx, 0) || math.IsNaN(sy) || math.IsInf(sy, 0) {
		return 0, 0, 0, false
	}
	return sx, sy, z, true
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	// Handle the division-by-zero singularity at r = 0.
	// Since lim(r->0) sin(r)/r = 1.0, we return 1.0.
	if r == 0 {
		return 1.0
	}
	return math.Sin(r) / r
}

// color maps a z value to a beautiful HSL color string using a premium color palette.
// Range: Deep Slate Blue (valleys) -> Magenta (mid-slopes) -> Vibrant Orange-Red (peaks)
func color(z, zmin, zmax float64) string {
	// Normalize z to [0, 1] range
	t := 0.0
	if zmax > zmin {
		t = (z - zmin) / (zmax - zmin)
	}
	if t < 0.0 {
		t = 0.0
	}
	if t > 1.0 {
		t = 1.0
	}

	var h, s, l float64
	if t < 0.5 {
		// Interpolate between Blue: HSL(220, 80%, 40%) and Magenta: HSL(300, 70%, 45%)
		f := t / 0.5
		h = 220.0 + (300.0-220.0)*f
		s = 80.0 + (70.0-80.0)*f
		l = 40.0 + (45.0-40.0)*f
	} else {
		// Interpolate between Magenta: HSL(300, 70%, 45%) and Orange/Red: HSL(375, 90%, 55%)
		// (Hue 375 is equivalent to hue 15 in modulo 360)
		f := (t - 0.5) / 0.5
		h = 300.0 + (375.0-300.0)*f
		if h >= 360.0 {
			h -= 360.0
		}
		s = 70.0 + (90.0-70.0)*f
		l = 45.0 + (55.0-45.0)*f
	}

	return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", h, s, l)
}