package qui

import (
	"github.com/qbradq/q2d"
)

type TextArea struct {
	BaseWidget
	Text string

	focused bool
	Width   int
	Height  int
}

func NewTextArea(text string) *TextArea {
	return &TextArea{Text: text}
}

func (t *TextArea) MinSize() Size {
	theme := t.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	// Arbitrary minimum size
	metrics := theme.Font.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	w := 100
	if t.Width > 0 {
		w = t.Width
	}
	h := lineHeight * 3
	if t.Height > 0 {
		h = t.Height
	}
	return Size{w, h}
}

func (t *TextArea) Event(evt Event) bool {
	switch event := evt.(type) {
	case MouseEvent:
		if event.TypeVal == EventMouseDown {
			// Focus handled by Master
			if t.Rect.Contains(event.Pos) {
				return true
			}
		}
	case TextInputEvent:
		if t.focused {
			t.Text += event.Text
			return true
		}
	case KeyEvent:
		if t.focused && event.TypeVal == EventKeyDown {
			if event.Key == KeyBackspace {
				if len(t.Text) > 0 {
					t.Text = t.Text[:len(t.Text)-1]
				}
				return true
			} else if event.Key == KeyEnter {
				t.Text += "\n"
				return true
			}
		}
	}
	return false
}

func (t *TextArea) Focus() {
	t.focused = true
}

func (t *TextArea) Unfocus() {
	t.focused = false
}

func (t *TextArea) FindWidgetAt(pos q2d.Point) Widget {
	if t.Rect.Contains(pos) {
		return t
	}
	return nil
}

func (t *TextArea) Draw(img *q2d.Image) {
	theme := t.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(t.Rect)
	defer img.PopSubImage()

	bgColor := theme.BackgroundColor.Darken(0.1)
	if t.focused {
		bgColor = theme.BackgroundColor.Lighten(0.1)
	}

	img.Fill(bgColor)
	img.Border(theme.BorderColor)

	displayText := t.Text
	if t.focused {
		displayText += "|"
	}

	img.Text(q2d.Point{theme.Padding, theme.Padding}, theme.TextColor, theme.Font, true, "%s", displayText)
}
