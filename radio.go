package qui

import (
	"github.com/qbradq/tremor/lib/q2d"
	"golang.org/x/image/font"
)

type RadioGroup struct {
	SelectedValue string
	OnChange      func(string)
	buttons       []*RadioButton
}

func NewRadioGroup(onChange func(string)) *RadioGroup {
	return &RadioGroup{
		OnChange: onChange,
		buttons:  make([]*RadioButton, 0),
	}
}

func (g *RadioGroup) Add(r *RadioButton) {
	g.buttons = append(g.buttons, r)
	r.Group = g
}

func (g *RadioGroup) Select(value string) {
	g.SelectedValue = value
	if g.OnChange != nil {
		g.OnChange(value)
	}
}

type RadioButton struct {
	BaseWidget
	Label string
	Value string
	Group *RadioGroup

	hovered bool
	pressed bool
	focused bool
}

func NewRadioButton(label, value string, group *RadioGroup) *RadioButton {
	rb := &RadioButton{
		Label: label,
		Value: value,
		Group: group,
	}
	if group != nil {
		group.Add(rb)
	}
	return rb
}

func (r *RadioButton) MinSize() Size {
	if DefaultTheme == nil || DefaultTheme.Font == nil {
		return Size{0, 0}
	}
	width := IconSize + DefaultTheme.Spacing
	width += font.MeasureString(DefaultTheme.Font, r.Label).Ceil()

	metrics := DefaultTheme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()
	if IconSize > height {
		height = IconSize
	}

	return Size{width + DefaultTheme.Padding*2, height + DefaultTheme.Padding*2}
}

func (r *RadioButton) Event(e Event) bool {
	switch evt := e.(type) {
	case MouseEvent:
		inRect := r.Rect.Contains(evt.Pos)

		if evt.TypeVal == EventMouseMove {
			wasHovered := r.hovered
			r.hovered = inRect
			return wasHovered || r.hovered
		}

		if evt.TypeVal == EventMouseDown && inRect {
			r.pressed = true
			return true
		}

		if evt.TypeVal == EventMouseUp {
			if r.pressed && inRect {
				r.Select()
			}
			r.pressed = false
			return inRect
		}
	case KeyEvent:
		if r.focused && evt.TypeVal == EventKeyDown {
			if evt.Key == KeyEnter || evt.Key == 32 { // Space
				r.Select()
				return true
			}
		}
	}
	return false
}

func (r *RadioButton) Select() {
	if r.Group != nil {
		r.Group.Select(r.Value)
	}
}

func (r *RadioButton) Focus() {
	r.focused = true
}

func (r *RadioButton) Unfocus() {
	r.focused = false
}

func (r *RadioButton) FindWidgetAt(pos q2d.Point) Widget {
	if r.Rect.Contains(pos) {
		return r
	}
	return nil
}

func (r *RadioButton) Draw(img *q2d.Image) {
	if DefaultTheme == nil {
		return
	}

	img.PushSubImage(r.Rect)
	defer img.PopSubImage()

	if r.focused {
		img.Fill(DefaultTheme.BackgroundColor.Lighten(0.1))
	}

	// Draw Icon
	icon := IconRadioOff
	if r.Group != nil && r.Group.SelectedValue == r.Value {
		icon = IconRadioOn
	}

	metrics := DefaultTheme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	contentHeight := textHeight
	if IconSize > contentHeight {
		contentHeight = IconSize
	}

	y := (r.Rect.Height() - contentHeight) / 2
	if y < 0 {
		y = 0
	}

	iconY := y + (contentHeight-IconSize)/2
	DrawIcon(img, icon, q2d.Point{DefaultTheme.Padding, iconY}, DefaultTheme.TextColor)

	textX := DefaultTheme.Padding + IconSize + DefaultTheme.Spacing
	textY := y + (contentHeight-textHeight)/2
	img.Text(q2d.Point{textX, textY}, DefaultTheme.TextColor, DefaultTheme.Font, false, "%s", r.Label)
}
