package theme

// Unified Pastelle Theme
//
// This theme adapts automatically to Fyne's light/dark variants.
// Fyne will call Color(name, variant). We use the variant to select from the
// Pastelle Light (teal-focused) or Pastelle Dark (muted) palettes.
//
// To use dynamic light/dark support:
//   a.Settings().SetTheme(theme.NewPastelleTheme())
//
// If the user / OS switches preference (or you toggle it in settings), Fyne
// will pass the appropriate fyne.ThemeVariant to Color() so the palette changes
// seamlessly.
//
// Light Palette (Teal):
//   Background: #F6F9F8
//   Foreground: #1F2A2E
//   Primary:    #2BB5A8
//   Secondary:  #7CCFCC
//   Tertiary:   #A8D5BA
//   Accent:     #2BB5A8
//   Success:    mix(Tertiary, Primary, 0.25)
//   Warning:    #FFC773
//   Error:      #FF8A8A
//
// Dark Palette (Pastelle):
//   Background: #333333   (from original "Text")
//   Foreground: #DDDDDD   (from original "Background")
//   Primary:    #BBBBBB
//   Secondary:  #727485
//   Tertiary:   #A5A88D
//   Accent:     #FF6F61
//   Success:    mix(Tertiary, Foreground, 0.15)
//   Warning:    fallback to dark theme
//   Error:      fallback to dark theme
//
// Any unhandled color names fall back to the underlying Fyne Light/Dark theme.
//
// NOTE:
// Icons, fonts, and sizes are delegated to the corresponding built-in theme
// to ensure consistency and full glyph coverage.

import (
	"image/color"

	"fyne.io/fyne/v2"
	fyneTheme "fyne.io/fyne/v2/theme"
)

// Ensure PastelleTheme implements fyne.Theme.
var _ fyne.Theme = (*PastelleTheme)(nil)

// PastelleTheme provides a unified theme for both variants.
type PastelleTheme struct {
	variant fyne.ThemeVariant
}

// NewPastelleTheme constructs the unified Pastelle theme.
func NewPastelleTheme() fyne.Theme {
	return &PastelleTheme{}
}

// NewPastelleLight constructs the Pastelle Light theme.
func NewPastelleLight() fyne.Theme {
	return &PastelleTheme{
		variant: fyneTheme.VariantLight,
	}
}

// NewPastelleDark constructs the Pastelle Dark theme.
func NewPastelleDark() fyne.Theme {
	return &PastelleTheme{
		variant: fyneTheme.VariantDark,
	}
}

// Color returns a color for the given name and variant.
func (t *PastelleTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if t.variant == fyneTheme.VariantDark {
		return t.darkColor(n)
	}
	return t.lightColor(n)
}

// lightColor maps color names for the light variant.
func (t *PastelleTheme) lightColor(n fyne.ThemeColorName) color.Color {
	// Light palette
	colorMap := NewColorMap(
		"#FFFFE6", //"#F6F9F8",
		"#1F2A2E",
		"#1C7685",
		"#7CCFCC",
		"#A8D5BA",
		"#3ABFB2",
		"#FFC773",
		"#FF8A8A",
		"#000000",
	)

	return colorMap.convertToColor(n, fyneTheme.VariantLight)
}

// darkColor maps color names for the dark variant.
func (t *PastelleTheme) darkColor(n fyne.ThemeColorName) color.Color {
	// Dark palette (re-mapped from original Pastelle Dark data)

	colorMap := &ColorMap{
		bg:        hex("#333333"),
		fg:        hex("#DDDDDD"),
		primary:   hex("#BBBBBB"),
		secondary: hex("#727485"),
		tertiary:  hex("#A5A88D"),
		accent:    hex("#FF6F61"),
		success:   mix(hex("#A5A88D"), hex("#DDDDDD"), 0.15),
		warning:   hex("#FFC98B"),
		errCol:    hex("#E07474"),
		separator: withAlpha(mix(hex("#DDDDDD"), hex("#333333"), 0.85), 96),
		shadow:    withAlpha(hex("#2BB5A8"), 60),
	}

	return colorMap.convertToColor(n, fyneTheme.VariantDark)
}

// Icon delegates to the base theme (icons are neutral).
func (t *PastelleTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	// Icons do not differ dramatically; use light theme icons for both.
	return fyneTheme.LightTheme().Icon(n)
}

// Font delegates to the variant-appropriate base theme.
func (t *PastelleTheme) Font(s fyne.TextStyle) fyne.Resource {
	// Use dark/light fonts from existing theme depending on style weight.
	// (Fyne's built-ins ensure symbol coverage.)
	return fyneTheme.DefaultTheme().Font(s)
}

// Size delegates to the default theme for consistency.
func (t *PastelleTheme) Size(n fyne.ThemeSizeName) float32 {
	return fyneTheme.DefaultTheme().Size(n)
}

// --- Helpers ---

func withAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = a
	return c
}

func clamp01(p float32) float32 {
	if p < 0 {
		return 0
	}
	if p > 1 {
		return 1
	}
	return p
}

func mix(a, b color.NRGBA, p float32) color.NRGBA {
	p = clamp01(p)
	q := 1 - p
	return color.NRGBA{
		R: uint8(float32(a.R)*q + float32(b.R)*p + 0.5),
		G: uint8(float32(a.G)*q + float32(b.G)*p + 0.5),
		B: uint8(float32(a.B)*q + float32(b.B)*p + 0.5),
		A: uint8(float32(a.A)*q + float32(b.A)*p + 0.5),
	}
}

func darken(c color.NRGBA, p float32) color.NRGBA {
	return mix(c, hex("#000000"), p)
}

func hex(s string) color.NRGBA {
	// Expect "#RRGGBB"
	if len(s) != 7 || s[0] != '#' {
		return color.NRGBA{A: 0xFF} // fallback to transparent black if malformed
	}
	r := (hexNibble(s[1]) << 4) | hexNibble(s[2])
	g := (hexNibble(s[3]) << 4) | hexNibble(s[4])
	b := (hexNibble(s[5]) << 4) | hexNibble(s[6])
	return color.NRGBA{R: r, G: g, B: b, A: 0xFF}
}

func hexNibble(b byte) uint8 {
	switch {
	case b >= '0' && b <= '9':
		return uint8(b - '0')
	case b >= 'a' && b <= 'f':
		return 10 + uint8(b-'a')
	case b >= 'A' && b <= 'F':
		return 10 + uint8(b-'A')
	default:
		return 0
	}
}

// helpers are defined in pastelle-dark.go to avoid duplication
type ColorMap struct {
	bg        color.NRGBA
	fg        color.NRGBA
	primary   color.NRGBA
	secondary color.NRGBA
	tertiary  color.NRGBA
	accent    color.NRGBA
	success   color.NRGBA
	warning   color.NRGBA
	errCol    color.NRGBA
	separator color.NRGBA
	shadow    color.NRGBA
}

func NewColorMap(bg, fg, primary, secondary, tertiary, accent, warning, errCol, shadow string) *ColorMap {
	return &ColorMap{
		bg:        hex(bg),
		fg:        hex(fg),
		primary:   hex(primary),
		secondary: hex(secondary),
		tertiary:  hex(tertiary),
		accent:    hex(accent),
		success:   mix(hex(tertiary), hex(primary), 0.25),
		warning:   hex(warning),
		errCol:    hex(errCol),
		separator: withAlpha(mix(hex(fg), hex(bg), 0.85), 96),
		shadow:    withAlpha(hex(shadow), 60),
	}
}

func (cm *ColorMap) convertToColor(n fyne.ThemeColorName, tv fyne.ThemeVariant) color.Color {
	switch n {
	case fyneTheme.ColorNameBackground:
		return cm.bg
	case fyneTheme.ColorNameForeground:
		return cm.fg
	case fyneTheme.ColorNamePrimary:
		return cm.primary
	case fyneTheme.ColorNameButton:
		return mix(cm.bg, cm.primary, 0.16)
	case fyneTheme.ColorNamePressed:
		return darken(cm.accent, 0.10)
	case fyneTheme.ColorNameHover:
		return withAlpha(cm.accent, 64)
	case fyneTheme.ColorNameFocus:
		return cm.accent
	case fyneTheme.ColorNamePlaceHolder:
		return mix(cm.fg, cm.bg, 0.55)
	case fyneTheme.ColorNameDisabled:
		return mix(cm.fg, cm.bg, 0.70)
	case fyneTheme.ColorNameHyperlink:
		return cm.accent
	case fyneTheme.ColorNameSelection:
		return mix(cm.bg, cm.accent, 0.28)
	case fyneTheme.ColorNameSeparator:
		return cm.separator
	case fyneTheme.ColorNameShadow:
		return cm.shadow
	case fyneTheme.ColorNameInputBackground:
		return mix(cm.bg, cm.secondary, 0.10)
	case fyneTheme.ColorNameScrollBar:
		return withAlpha(cm.secondary, 180)
	case fyneTheme.ColorNameSuccess:
		return cm.success
	case fyneTheme.ColorNameWarning:
		return cm.warning
	case fyneTheme.ColorNameError:
		return cm.errCol
	default:
		return fyneTheme.DefaultTheme().Color(n, tv)
	}
}
