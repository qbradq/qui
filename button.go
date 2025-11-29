package qui

import (
	"github.com/qbradq/tremor/lib/q2d"
	"golang.org/x/image/font"
)

type Button struct {
	BaseWidget
	Text    string
	Icon    Icon
	OnClick func()

	hovered bool
	pressed bool
}

func NewButton(text string, onClick func()) *Button {
	return &Button{
		Text:    text,
		OnClick: onClick,
	}
}

func (b *Button) MinSize() Size {
	theme := b.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	width := font.MeasureString(theme.Font, b.Text).Ceil()
	metrics := theme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()

	if b.Icon != IconNone {
		width += IconSize
		if b.Text != "" {
			width += theme.Spacing
		}
		if IconSize > height {
			height = IconSize
		}
	}

	return Size{width + theme.Padding*4, height + theme.Padding*4}
}

func (b *Button) Event(e Event) bool {
	switch evt := e.(type) {
	case MouseEvent:
		inRect := b.Rect.Contains(evt.Pos)

		if evt.TypeVal == EventMouseMove {
			wasHovered := b.hovered
			b.hovered = inRect
			return wasHovered || b.hovered
		}

		if evt.TypeVal == EventMouseDown && inRect {
			b.pressed = true
			return true
		}

		if evt.TypeVal == EventMouseUp {
			if b.pressed && inRect {
				if b.OnClick != nil {
					b.OnClick()
				}
			}
			b.pressed = false
			return inRect
		}
	}
	return false
}

func (b *Button) Draw(img *q2d.Image) {
	theme := b.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(b.Rect)
	defer img.PopSubImage()

	bgColor := theme.ButtonColor
	if b.pressed {
		bgColor = bgColor.Darken(0.2)
	} else if b.hovered {
		bgColor = theme.ButtonHoverColor
	}

	img.Fill(bgColor)
	img.Border(theme.BorderColor)

	textWidth := font.MeasureString(theme.Font, b.Text).Ceil()
	if b.Icon != IconNone {
		textWidth += IconSize
		if b.Text != "" {
			textWidth += theme.Spacing
		}
	}

	metrics := theme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	// Use max height of text or icon
	contentHeight := textHeight
	if b.Icon != IconNone && IconSize > contentHeight {
		contentHeight = IconSize
	}

	x := (b.Rect.Width() - textWidth) / 2
	y := (b.Rect.Height() - contentHeight) / 2

	// Ensure x and y are not negative
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	if b.Icon != IconNone {
		iconY := y + (contentHeight-IconSize)/2
		DrawIcon(img, b.Icon, q2d.Point{x, iconY}, theme.TextColor)
		x += IconSize + theme.Spacing
	}

	textY := y + (contentHeight-textHeight)/2
	img.Text(q2d.Point{x, textY}, theme.TextColor, theme.Font, false, "%s", b.Text)
}
