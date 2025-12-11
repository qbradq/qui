package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Checkbox struct {
	BaseWidget
	Label    string
	Checked  bool
	OnChange func(bool)

	hovered bool
	pressed bool
	focused bool
}

func NewCheckbox(label string, checked bool, onChange func(bool)) *Checkbox {
	return &Checkbox{
		Label:    label,
		Checked:  checked,
		OnChange: onChange,
	}
}

func (c *Checkbox) MinSize() Size {
	if DefaultTheme == nil || DefaultTheme.Font == nil {
		return Size{0, 0}
	}
	width := IconSize + DefaultTheme.Spacing
	width += font.MeasureString(DefaultTheme.Font, c.Label).Ceil()

	metrics := DefaultTheme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()
	if IconSize > height {
		height = IconSize
	}

	return Size{width + (DefaultTheme.Padding.Left + DefaultTheme.Padding.Right), height + (DefaultTheme.Padding.Top + DefaultTheme.Padding.Bottom)}
}

func (c *Checkbox) Event(e Event) bool {
	switch evt := e.(type) {
	case MouseEvent:
		inRect := c.Rect.Contains(evt.Pos)

		if evt.TypeVal == EventMouseMove {
			wasHovered := c.hovered
			c.hovered = inRect
			return wasHovered || c.hovered
		}

		if evt.TypeVal == EventMouseDown && inRect {
			c.pressed = true
			return true
		}

		if evt.TypeVal == EventMouseUp {
			if c.pressed && inRect {
				c.Toggle()
			}
			c.pressed = false
			return inRect
		}
	case KeyEvent:
		if c.focused && evt.TypeVal == EventKeyDown {
			if evt.Key == KeyEnter || evt.Key == 32 { // Space
				c.Toggle()
				return true
			}
		}
	}
	return false
}

func (c *Checkbox) Toggle() {
	c.Checked = !c.Checked
	if c.OnChange != nil {
		c.OnChange(c.Checked)
	}
}

func (c *Checkbox) Focus() {
	c.focused = true
}

func (c *Checkbox) Unfocus() {
	c.focused = false
}

func (c *Checkbox) FindWidgetAt(pos q2d.Point) Widget {
	if c.Rect.Contains(pos) {
		return c
	}
	return nil
}

func (c *Checkbox) Draw(img *q2d.Image) {
	if DefaultTheme == nil {
		return
	}

	img.PushSubImage(c.Rect)
	defer img.PopSubImage()

	if c.focused {
		img.Fill(DefaultTheme.BackgroundColor.Lighten(0.1))
	}

	// Draw Icon
	icon := IconUncheck
	if c.Checked {
		icon = IconCheck
	}

	metrics := DefaultTheme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	contentHeight := textHeight
	if IconSize > contentHeight {
		contentHeight = IconSize
	}

	y := (c.Rect.Height() - contentHeight) / 2
	if y < 0 {
		y = 0
	}

	iconY := y + (contentHeight-IconSize)/2
	DrawIcon(img, icon, q2d.Point{DefaultTheme.Padding.Left, iconY}, DefaultTheme.TextColor)

	textX := DefaultTheme.Padding.Left + IconSize + DefaultTheme.Spacing
	textY := y + (contentHeight-textHeight)/2
	img.Text(q2d.Point{textX, textY}, DefaultTheme.TextColor, DefaultTheme.Font, false, "%s", c.Label)
}
