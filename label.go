package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Label struct {
	BaseWidget
	Text string
	Icon Icon
}

func NewLabel(text string) *Label {
	return &Label{Text: text}
}

func (l *Label) MinSize() Size {
	theme := l.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	width := font.MeasureString(theme.Font, l.Text).Ceil()
	metrics := theme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()

	if l.Icon != IconNone {
		width += IconSize + theme.Spacing
		if IconSize > height {
			height = IconSize
		}
	}

	return Size{width, height}
}

func (l *Label) Draw(img *q2d.Image) {
	theme := l.GetTheme()
	if theme == nil {
		return
	}
	img.PushSubImage(l.Rect)
	defer img.PopSubImage()

	x := 0

	metrics := theme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()

	// Center vertically
	y := (l.Rect.Height() - textHeight) / 2
	if y < 0 {
		y = 0
	}

	if l.Icon != IconNone {
		iconY := (l.Rect.Height() - IconSize) / 2
		DrawIcon(img, l.Icon, q2d.Point{x, iconY}, theme.TextColor)
		x += IconSize + theme.Spacing
	}

	img.Text(q2d.Point{x, y}, theme.TextColor, theme.Font, true, "%s", l.Text)
}
