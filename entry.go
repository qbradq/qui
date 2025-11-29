package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Entry struct {
	BaseWidget
	Input *TextInput
	Width int
}

func NewEntry(initialText string, t EntryType) *Entry {
	return &Entry{
		Input: NewTextInput(initialText, t),
	}
}

func (e *Entry) MinSize() Size {
	theme := e.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}

	inputSize := e.Input.MinSize()
	width := inputSize.Width

	// Minimum width override
	if e.Width > 0 {
		width = e.Width
	} else {
		// Default min width if not set
		minW := font.MeasureString(theme.Font, "MMMMMMMMMM").Ceil()
		if width < minW {
			width = minW
		}
	}

	height := inputSize.Height

	return Size{width + theme.Padding*2, height + theme.Padding*2}
}

func (e *Entry) Event(evt Event) bool {
	// Forward events to Input
	// But first, handle focus click on the container
	switch event := evt.(type) {
	case MouseEvent:
		if event.TypeVal == EventMouseDown {
			if e.Rect.Contains(event.Pos) {
				// We don't need to manually focus Input here if Master handles it?
				// But Master calls Focus() on the widget returned by FindWidgetAt.
				// If FindWidgetAt returns e.Input, then e.Input gets focus.
				// If FindWidgetAt returns e, then e gets focus.
				// e.Input is internal.
				return true
			}
		}
	}

	// Delegate to Input
	return e.Input.Event(evt)
}

func (e *Entry) Focus() {
	e.Input.Focus()
}

func (e *Entry) Unfocus() {
	e.Input.Unfocus()
}

func (e *Entry) Draw(img *q2d.Image) {
	theme := e.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(e.Rect)

	bgColor := theme.BackgroundColor.Darken(0.1)
	if e.Input.focused {
		bgColor = theme.BackgroundColor.Lighten(0.1)
	}

	img.Fill(bgColor)
	img.Border(theme.BorderColor)
	img.PopSubImage()

	// Update Input rect to be inside padding
	// We need to set Input.Rect relative to what?
	// Widgets usually store absolute Rects (in window coordinates).
	// So we should calculate absolute position.
	// e.Rect is absolute.

	inputRect := q2d.Rectangle{
		e.Rect.X() + theme.Padding,
		e.Rect.Y() + theme.Padding,
		e.Rect.Width() - theme.Padding*2,
		e.Rect.Height() - theme.Padding*2,
	}
	e.Input.Rect = inputRect

	// Draw Input
	// Input.Draw expects to draw at its Rect.
	// Since we are in a SubImage (e.Rect), drawing at inputRect (absolute) might be wrong if SubImage offsets origin.
	// q2d.Image.PushSubImage sets the clip rect. It does NOT translate the coordinate system.
	// So drawing at absolute coordinates is correct.
	e.Input.Draw(img)
}

func (e *Entry) FindWidgetAt(pos q2d.Point) Widget {
	if e.Rect.Contains(pos) {
		// Return e, so e gets focus?
		// If we return e, Master calls e.Focus(), which calls e.Input.Focus(). Correct.
		// If we return e.Input, Master calls e.Input.Focus(). Also correct.
		// But e.Input.Rect is inside e.Rect.
		// If we click on the border (padding), e.Input.Rect might not contain pos.
		// So we should return e.
		return e
	}
	return nil
}

// Proxy methods for convenience/compatibility if needed
func (e *Entry) SetText(t string) {
	e.Input.Text = t
}

func (e *Entry) GetText() string {
	return e.Input.Text
}
