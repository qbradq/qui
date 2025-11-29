package qui

import (
	"strconv"
	"strings"

	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type EntryType int

const (
	EntryText EntryType = iota
	EntryPassword
	EntryInteger
	EntryFloat
)

type TextInput struct {
	BaseWidget
	Text string
	Type EntryType

	focused   bool
	cursorPos int
}

func NewTextInput(initialText string, t EntryType) *TextInput {
	return &TextInput{
		Text:      initialText,
		Type:      t,
		cursorPos: len([]rune(initialText)),
	}
}

func (t *TextInput) MinSize() Size {
	theme := t.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	// Minimum width for some text
	width := font.MeasureString(theme.Font, t.Text+"|").Ceil()
	if width == 0 {
		width = font.MeasureString(theme.Font, "M").Ceil()
	}
	metrics := theme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()

	return Size{width, height}
}

func (t *TextInput) Event(evt Event) bool {
	switch event := evt.(type) {
	case MouseEvent:
		if event.TypeVal == EventMouseDown {
			if t.Rect.Contains(event.Pos) {
				return true
			}
		}
	case TextInputEvent:
		if t.focused {
			t.insertText(event.Text)
			return true
		}
	case KeyEvent:
		if t.focused && event.TypeVal == EventKeyDown {
			runes := []rune(t.Text)
			switch event.Key {
			case KeyLeft:
				if t.cursorPos > 0 {
					t.cursorPos--
				}
				return true
			case KeyRight:
				if t.cursorPos < len(runes) {
					t.cursorPos++
				}
				return true
			case KeyHome:
				t.cursorPos = 0
				return true
			case KeyEnd:
				t.cursorPos = len(runes)
				return true
			case KeyBackspace:
				if t.cursorPos > 0 {
					t.Text = string(append(runes[:t.cursorPos-1], runes[t.cursorPos:]...))
					t.cursorPos--
				}
				return true
			case KeyDelete:
				if t.cursorPos < len(runes) {
					t.Text = string(append(runes[:t.cursorPos], runes[t.cursorPos+1:]...))
				}
				return true
			}
		}
	}
	return false
}

func (t *TextInput) Focus() {
	t.focused = true
}

func (t *TextInput) Unfocus() {
	t.focused = false
}

func (t *TextInput) insertText(text string) {
	runes := []rune(t.Text)
	insert := []rune(text)

	// Insert at cursor
	newRunes := append(runes[:t.cursorPos], append(insert, runes[t.cursorPos:]...)...)
	newText := string(newRunes)

	if t.Type == EntryInteger {
		if _, err := strconv.Atoi(newText); err != nil && newText != "-" && newText != "" {
			return
		}
	} else if t.Type == EntryFloat {
		if _, err := strconv.ParseFloat(newText, 64); err != nil && newText != "-" && newText != "" && newText != "." && newText != "-." {
			return
		}
	}
	t.Text = newText
	t.cursorPos += len(insert)
}

func (t *TextInput) Draw(img *q2d.Image) {
	theme := t.GetTheme()
	if theme == nil {
		return
	}

	displayText := t.Text
	if t.Type == EntryPassword {
		displayText = strings.Repeat("*", len([]rune(t.Text)))
	}

	img.PushSubImage(t.Rect)
	defer img.PopSubImage()

	img.Text(q2d.Point{0, 0}, theme.TextColor, theme.Font, false, "%s", displayText)

	if t.focused {
		runes := []rune(displayText)
		cursorX := 0
		if t.cursorPos > 0 {
			if t.cursorPos > len(runes) {
				t.cursorPos = len(runes)
			}
			cursorX = font.MeasureString(theme.Font, string(runes[:t.cursorPos])).Ceil()
		}

		metrics := theme.Font.Metrics()
		height := (metrics.Ascent + metrics.Descent).Ceil()

		img.VLine(cursorX, 0, height, 1, theme.TextColor)
	}
}

func (t *TextInput) FindWidgetAt(pos q2d.Point) Widget {
	if t.Rect.Contains(pos) {
		return t
	}
	return nil
}

func (t *TextInput) GetText() string {
	return t.Text
}

func (t *TextInput) SetText(text string) {
	t.Text = text
}
