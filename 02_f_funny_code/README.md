# 3-D SVG Surface Plotter (02_f_funny_code)

This project is a 3-D wireframe surface plotter that renders a two-variable function $z = f(x, y)$ as an isometric SVG drawing. The original implementation was adapted from Chapter 3 of *The Go Programming Language* book, but has been fixed and enhanced with premium aesthetic improvements and error handling.

---

## Bugs Fixed

### 1. Constant Declaration Syntax Error
In Go, semicolons are automatically inserted at the end of a line if it ends with an identifier (like `cells`, `xyrange`, etc.). The original code split the constant declarations across lines in a way that resulted in syntax errors (e.g. `cells` on one line and `= 100` on the next, producing `cells; = 100`).
- **Fix:** Joined the constants onto single, readable declaration lines.

### 2. Singularity at the Origin ($r = 0$)
The function plotted is $f(x, y) = \frac{\sin r}{r}$, where $r = \sqrt{x^2 + y^2}$. At the origin $(0, 0)$, $r = 0$, leading to $\frac{\sin 0}{0}$ which evaluates to `NaN` (Not a Number) in floating-point math.
- **Fix:** Handled $r = 0$ as a special case in `f(x, y)` to return `1.0` (since $\lim_{r\to 0}\frac{\sin r}{r} = 1$).

### 3. Safety Checks for Coordinates
If any calculation returns `NaN` or `±Inf`, printing those coordinate values into the `<polygon points="..." />` tag generates invalid SVG markup.
- **Fix:** Updated the `corner()` function to validate coordinates and return a boolean indicating success. Polygons containing invalid values are safely skipped.

---

## Aesthetic Enhancements

To create a visually stunning output, we replaced the plain white-fill/grey-stroke grid with a **height-based color gradient**:
- The program calculates the average height ($z_{avg}$) of each grid cell.
- The height is normalized against the function's range ($[-0.2172, 1.0]$).
- A custom interpolation function maps this normalized value to a premium HSL color gradient:
  - **Valleys (Low):** Deep Slate Blue `hsl(220, 80%, 40%)`
  - **Slopes (Mid):** Magenta `hsl(300, 70%, 45%)`
  - **Peaks (High):** Vibrant Orange-Red `hsl(15, 90%, 55%)`

Additionally, the program now automatically writes the generated SVG output to a file named `surface.svg` in the current directory, in addition to writing it to standard output.

---

## Running the Program

Navigate to this directory (`02_f_funny_code`) and run:

```bash
go run main.go
```

This will:
1. Print the raw SVG XML to your terminal.
2. Generate/overwrite the file `surface.svg` in the directory, which you can open with any web browser or vector image viewer to see the beautiful 3-D wireframe model.
