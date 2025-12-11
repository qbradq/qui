package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Theme struct {
	BackgroundColor  q2d.Color
	TextColor        q2d.Color
	ButtonColor      q2d.Color
	ButtonHoverColor q2d.Color
	BorderColor      q2d.Color
	PrimaryColor     q2d.Color
	SecondaryColor   q2d.Color
	Font             font.Face
	IconSheet        *q2d.Image
	Spacing          int
	Padding          Padding
}

var DefaultTheme *Theme

func InitTheme(f font.Face) {
	DefaultTheme = &Theme{
		BackgroundColor:  q2d.Color{30, 30, 30, 255},
		TextColor:        q2d.Color{220, 220, 220, 255},
		ButtonColor:      q2d.Color{60, 60, 60, 255},
		ButtonHoverColor: q2d.Color{80, 80, 80, 255},
		BorderColor:      q2d.Color{100, 100, 100, 255},
		PrimaryColor:     q2d.Color{0, 140, 255, 255},
		SecondaryColor:   q2d.Color{50, 50, 50, 255},
		Font:             f,
		IconSheet:        CreateDummyIconSheet(),
		Spacing:          5,
		Padding:          Padding{Top: 2, Right: 5, Bottom: 2, Left: 5},
	}
}

func GenerateTheme(base, text, complement q2d.Color, f font.Face) *Theme {
	return &Theme{
		BackgroundColor:  base,
		TextColor:        text,
		ButtonColor:      base.Lighten(0.1),
		ButtonHoverColor: base.Lighten(0.2),
		BorderColor:      base.Lighten(0.3),
		PrimaryColor:     complement,
		SecondaryColor:   base.Lighten(0.05),
		Font:             f,
		IconSheet:        CreateDummyIconSheet(),
		Spacing:          5,
		Padding:          Padding{Top: 2, Right: 5, Bottom: 2, Left: 5},
	}
}

// GenerateThemeFromColor generates a theme based on a single color.
// It uses HSL adjustments to create a dark background and appropriate accents.
func GenerateThemeFromColor(c q2d.Color, f font.Face) *Theme {
	// Base is a darkened version of the input color
	base := c.Darken(0.85)

	// Text is high contrast (white for dark backgrounds)
	text := q2d.Color{240, 240, 240, 255}

	// Complement/Primary is the input color (bright)
	// We ensure it has high saturation and lightness for "bright form"
	// But if input is already "full-intensity", we can just use it.
	// Or we can force it to be bright?
	// Let's just use c as requested "bright form of the base color for comp".
	comp := c

	return GenerateTheme(base, text, comp, f)
}
