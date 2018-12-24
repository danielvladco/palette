package palette

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func Parse(colorStr string) color.Color {
	colorStr = regexp.MustCompile(`^\s+/`).ReplaceAllString(colorStr, "")
	colorStr = regexp.MustCompile(`\s+$/`).ReplaceAllString(colorStr, "")

	const (
		cssInteger       = "[-\\+]?\\d+%?"
		cssNumber        = "[-\\+]?\\d*\\.\\d+%?"
		cssUnit          = "(?:" + cssNumber + ")|(?:" + cssInteger + ")"
		permissiveMatch3 = "[\\s|\\(]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")\\s*\\)?"
		permissiveMatch4 = "[\\s|\\(]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")\\s*\\)?"
		rgb              = "rgb" + permissiveMatch3
		rgba             = "RGBA" + permissiveMatch4
		hsl              = "hsl" + permissiveMatch3
		hsla             = "hsla" + permissiveMatch4
		hsv              = "hsv" + permissiveMatch3
		hsva             = "hsva" + permissiveMatch4
		hex3             = `^#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})$`
		hex6             = `^#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$`
		hex4             = `^#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})$`
		hex8             = `^#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$`
	)

	if match := regexp.MustCompile(rgb).FindAllStringSubmatch(colorStr, -1); match != nil {
		return color.RGBA{
			R: uint8(colorStringToUint8(match[0][1])),
			G: uint8(colorStringToUint8(match[0][2])),
			B: uint8(colorStringToUint8(match[0][3])),
			A: 255,
		}
	}
	if match := regexp.MustCompile(rgba).FindAllStringSubmatch(colorStr, -1); match != nil {
		return color.RGBA{
			R: uint8(colorStringToUint8(match[0][1])),
			G: uint8(colorStringToUint8(match[0][2])),
			B: uint8(colorStringToUint8(match[0][3])),
			A: uint8(alphaStringToUint8(match[0][4])),
		}
	}
	if match := regexp.MustCompile(hsl).FindAllStringSubmatch(colorStr, -1); match != nil {
		return nil // TODO parse from hsl
	}
	if match := regexp.MustCompile(hsla).FindAllStringSubmatch(colorStr, -1); match != nil {
		return nil // TODO parse from HSLA
	}
	if match := regexp.MustCompile(hsv).FindAllStringSubmatch(colorStr, -1); match != nil {
		return nil // TODO parse from hsv
	}
	if match := regexp.MustCompile(hsva).FindAllStringSubmatch(colorStr, -1); match != nil {
		return nil // TODO parse from hsva
	}
	if match := regexp.MustCompile(hex8).FindAllStringSubmatch(colorStr, -1); match != nil {
		return color.RGBA{
			R: uint8(parseFromHex(match[0][1])),
			G: uint8(parseFromHex(match[0][2])),
			B: uint8(parseFromHex(match[0][3])),
			A: uint8(parseFromHex(match[0][4])),
		}
	}
	if match := regexp.MustCompile(hex6).FindAllStringSubmatch(colorStr, -1); match != nil {
		return color.RGBA{
			R: uint8(parseFromHex(match[0][1])),
			G: uint8(parseFromHex(match[0][2])),
			B: uint8(parseFromHex(match[0][3])),
			A: 255,
		}
	}
	if match := regexp.MustCompile(hex4).FindAllStringSubmatch(colorStr, -1); match != nil {
		return color.RGBA{
			R: uint8(parseFromHex(match[0][1])),
			G: uint8(parseFromHex(match[0][2])),
			B: uint8(parseFromHex(match[0][3])),
			A: uint8(parseFromHex(match[0][4])),
		}
	}
	if match := regexp.MustCompile(hex3).FindAllStringSubmatch(colorStr, -1); match != nil {
		return color.RGBA{
			R: uint8(parseFromHex(match[0][1])),
			G: uint8(parseFromHex(match[0][2])),
			B: uint8(parseFromHex(match[0][3])),
			A: 255,
		}
	}

	return nil
}

func Random() color.Color {
	return color.RGBA64{
		R: uint16(rand.Intn(0xffff + 1)),
		G: uint16(rand.Intn(0xffff + 1)),
		B: uint16(rand.Intn(0xffff + 1)),
		A: 0xffff,
	}
}

// saturate by a percent
// if want to darken give negative amount
func Lighten(hsla HSLA, amount float64) HSLA {
	hsla.L = math.Min(1, math.Max(0, amount/100))
	return hsla
}

// saturate by a percent
// if want to desaturate give negative amount
func Saturate(hsla HSLA, amount float64) HSLA {
	hsla.S = math.Min(1, math.Max(0, hsla.S+(amount/100)))
	return hsla
}

// returns the grayscale of an grayscale
func Greyscale(hsla HSLA) HSLA {
	return Saturate(hsla, -100)
}

func Spin(hsla HSLA, amount float64) HSLA {
	hue := math.Mod(hsla.H+amount, 360)
	if hue < 0 {
		hsla.H = 360 + hue
	} else {
		hsla.H = hue
	}
	return hsla
}

func Tetrad(col HSLA) (color.Color, color.Color, color.Color, color.Color) {
	hsla1 := col
	hsla2 := col
	hsla2.H = math.Mod(hsla1.H+90, 360)
	hsla3 := col
	hsla3.H = math.Mod(hsla1.H+180, 360)
	hsla4 := col
	hsla4.H = math.Mod(hsla1.H+270, 360)
	return hsla1, hsla2, hsla3, hsla4
}

func Triad(col HSLA) (color.Color, color.Color, color.Color) {
	hsla1 := col
	hsla2 := hsla1
	hsla2.H = math.Mod(hsla1.H+120, 360)
	hsla3 := hsla1
	hsla3.H = math.Mod(hsla1.H+240, 360)
	return hsla1, hsla2, hsla3
}

func Brighten(col RGBA, amount float64) color.Color {
	col.R = math.Max(0, math.Min(1, float64(col.R)+(amount/100)))
	col.G = math.Max(0, math.Min(1, float64(col.G)+(amount/100)))
	col.B = math.Max(0, math.Min(1, float64(col.B)+(amount/100)))
	return col
}

func Mix(rgba1, rgba2 RGBA, amount float64) RGBA {
	rgba1.R = ((rgba2.R - rgba1.R) * (amount / 100)) + rgba1.R
	rgba1.G = ((rgba2.G - rgba1.G) * (amount / 100)) + rgba1.G
	rgba1.B = ((rgba2.B - rgba1.B) * (amount / 100)) + rgba1.B
	rgba1.A = ((rgba2.A - rgba1.A) * (amount / 100)) + rgba1.A
	return rgba1
}

func Multiply(rgba1, rgba2 RGBA) RGBA {
	rgba1.R = rgba1.R * rgba2.R
	rgba1.G = rgba1.G * rgba2.G
	rgba1.B = rgba1.B * rgba2.B
	rgba1.A = rgba1.A * rgba2.A
	return rgba1
}

func ToHex(col color.Color) string {
	r, g, b, a := col.RGBA()
	if a == 0xffff {
		return fmt.Sprintf("#%02X%02X%02X", r>>8, g>>8, b>>8)
	}

	return fmt.Sprintf("#%02X%02X%02X%02X", r>>8, g>>8, b>>8, a>>8)
}

func NewRGBA(col color.Color) RGBA {
	if col == nil {
		return RGBA{}
	}
	r, g, b, a := col.RGBA()
	rgba := RGBA{}
	rgba.R = float64(r) / 0xffff
	rgba.G = float64(g) / 0xffff
	rgba.B = float64(b) / 0xffff
	rgba.A = float64(a) / 0xffff
	return rgba
}

func NewHSLA(col color.Color) HSLA {
	if col == nil {
		return HSLA{}
	}
	switch c := col.(type) {
	case *HSLA:
		return *c
	case HSLA:
		return c
	}

	r, g, b, a := col.RGBA()
	rgba := new(RGBA)
	rgba.R = float64(r) / 0xffff
	rgba.G = float64(g) / 0xffff
	rgba.B = float64(b) / 0xffff
	rgba.A = float64(a) / 0xffff
	return rgba.ToHSLA()
}

// A color is stored internally using sRGB (standard RGB) values in the range 0-1
type RGBA struct {
	R, G, B, A float64
}

func (rgba RGBA) Mix(rgba2 RGBA, amount float64) RGBA { return Mix(rgba, rgba2, amount) }
func (rgba RGBA) Multiply(rgba2 RGBA) RGBA            { return Multiply(rgba, rgba2) }

func (rgba RGBA) ToHex() string {
	if rgba.A == 1 {
		return fmt.Sprintf("#%02X%02X%02X", uint8(rgba.R*0xff+0.5), uint8(rgba.G*0xff+0.5), uint8(rgba.B*0xff+0.5))
	}

	return fmt.Sprintf("#%02X%02X%02X%02X", uint8(rgba.R*0xff+0.5), uint8(rgba.G*0xff+0.5), uint8(rgba.B*0xff+0.5), uint8(rgba.A*0xff+0.5))
}

func (rgba RGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(rgba.R*0xffff + 0.5)
	g = uint32(rgba.G*0xffff + 0.5)
	b = uint32(rgba.B*0xffff + 0.5)
	a = uint32(rgba.A*0xffff + 0.5)
	return
}

// Hsl returns the Hue [0..360], Saturation [0..1], Luminance (lightness) [0..1] of the color and alpha [0..1]
func (rgba RGBA) ToHSLA() (hsla HSLA) {
	var h, s, l float64
	min := math.Min(math.Min(rgba.R, rgba.G), rgba.B)
	max := math.Max(math.Max(rgba.R, rgba.G), rgba.B)

	l = (max + min) / 2

	if min == max {
		s = 0
		h = 0
	} else {
		if l < 0.5 {
			s = (max - min) / (max + min)
		} else {
			s = (max - min) / (2.0 - max - min)
		}

		if max == rgba.R {
			h = (rgba.G - rgba.B) / (max - min)
		} else if max == rgba.G {
			h = 2.0 + (rgba.B-rgba.R)/(max-min)
		} else {
			h = 4.0 + (rgba.R-rgba.G)/(max-min)
		}

		h *= 60

		if h < 0 {
			h += 360
		}
	}

	return HSLA{H: h, S: s, L: l, A: rgba.A}
}

type HSLA struct {
	H, S, L, A float64
}

func (hsla HSLA) RGBA() (r, g, b, a uint32)                                    { return hsla.ToRGBA().RGBA() }
func (hsla HSLA) Triad() (color.Color, color.Color, color.Color)               { return Triad(hsla) }
func (hsla HSLA) Tetrad() (color.Color, color.Color, color.Color, color.Color) { return Tetrad(hsla) }
func (hsla HSLA) Saturate(amount float64) HSLA                                 { return Saturate(hsla, amount) }
func (hsla HSLA) Lighten(amount float64) HSLA                                  { return Lighten(hsla, amount) }
func (hsla HSLA) Greyscale() HSLA                                              { return Greyscale(hsla) }
func (hsla HSLA) Spin(amount float64) HSLA                                     { return Spin(hsla, amount) }

// Hsl creates a new RGBA color from hsl
func (hsla HSLA) ToRGBA() RGBA {
	h, s, l, a := hsla.H, hsla.S, hsla.L, hsla.A
	if s == 0 {
		return RGBA{l, l, l, a}
	}

	var r, g, b, t1, t2, tr, tg, tb float64

	if l < 0.5 {
		t1 = l * (1.0 + s)
	} else {
		t1 = l + s - l*s
	}

	t2 = 2*l - t1
	h = h / 360
	tr = h + 1.0/3.0
	tg = h
	tb = h - 1.0/3.0

	if tr < 0 {
		tr++
	}
	if tr > 1 {
		tr--
	}
	if tg < 0 {
		tg++
	}
	if tg > 1 {
		tg--
	}
	if tb < 0 {
		tb++
	}
	if tb > 1 {
		tb--
	}

	// Red
	if 6*tr < 1 {
		r = t2 + (t1-t2)*6*tr
	} else if 2*tr < 1 {
		r = t1
	} else if 3*tr < 2 {
		r = t2 + (t1-t2)*(2.0/3.0-tr)*6
	} else {
		r = t2
	}

	// Green
	if 6*tg < 1 {
		g = t2 + (t1-t2)*6*tg
	} else if 2*tg < 1 {
		g = t1
	} else if 3*tg < 2 {
		g = t2 + (t1-t2)*(2.0/3.0-tg)*6
	} else {
		g = t2
	}

	// Blue
	if 6*tb < 1 {
		b = t2 + (t1-t2)*6*tb
	} else if 2*tb < 1 {
		b = t1
	} else if 3*tb < 2 {
		b = t2 + (t1-t2)*(2.0/3.0-tb)*6
	} else {
		b = t2
	}

	return RGBA{r, g, b, a}
}

// parses string color to its uint8 representation
// allowed values may be:
//	0% <= n% <= 100%
//	0 <= n <= 1
//	0 <= n <= 255
func colorStringToUint8(n string) uint8 {
	if strings.HasSuffix(n, "%") {
		n := strings.TrimSuffix(n, "%")
		nr, _ := strconv.ParseFloat(n, 0)
		if nr < 0 {
			return 0
		}
		if nr > 100 {
			return 255
		}

		return uint8(math.Round((nr * 255) / 100))
	}
	nr, _ := strconv.ParseFloat(n, 0)
	if nr < 0 {
		return 0
	}
	if nr > 255 {
		return 255
	}
	return uint8(math.Round(nr))
}

// parses string alpha to its uint8 representation
// allowed values may be:
//	0% <= n% <= 100%
//	0 <= n <= 1
func alphaStringToUint8(n string) uint8 {
	if strings.HasSuffix(n, "%") {
		n := strings.TrimSuffix(n, "%")
		nr, _ := strconv.ParseFloat(n, 0)
		if nr < 0 {
			return 0
		}
		if nr > 100 {
			return 255
		}

		return uint8(math.Round((nr * 255) / 100))
	}
	nr, _ := strconv.ParseFloat(n, 0)
	if nr < 0 {
		return 0
	}
	if nr > 1 {
		return 1
	}
	return uint8(math.Round(nr * 255))
}

func parseFromHex(hex string) uint8 {
	i, _ := strconv.ParseUint(hex, 16, 0)
	return uint8(i)
}
