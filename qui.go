package qui

import (
	"github.com/qbradq/tremor/lib/q2d"
)

type Size struct {
	Width, Height int
}

type ListItem struct {
	Text string
	Icon Icon
}

type OverlayManager interface {
	PushOverlay(w Widget)
	PopOverlay()
}

type Dismissable interface {
	OnDismiss()
}

type ManagedOverlay interface {
	SetOverlayManager(m OverlayManager)
}

type Focusable interface {
	Focus()
	Unfocus()
}

type Widget interface {
	// Layout calculates the size and position of the widget and its children.
	// It receives the available space constraints.
	Layout(available Size) Size

	// Draw renders the widget onto the q2d.Image.
	Draw(img *q2d.Image)

	// Event handles input events (mouse, keyboard).
	// Returns true if the event was consumed.
	Event(e Event) bool

	// MinSize returns the minimum size required by the widget.
	MinSize() Size

	// SetRect sets the absolute position and size of the widget.
	SetRect(r q2d.Rectangle)

	// GetRect returns the absolute position and size of the widget.
	GetRect() q2d.Rectangle

	// GetTooltip returns the tooltip text.
	GetTooltip() string

	// GetTheme returns the widget's theme or the default theme.
	GetTheme() *Theme
	// SetTheme sets the widget's theme.
	SetTheme(t *Theme)

	// FindWidgetAt returns the widget at the given position, or nil.
	FindWidgetAt(pos q2d.Point) Widget

	// IsFill returns true if the widget should fill available space in its parent container.
	IsFill() bool
}

// BaseWidget can be embedded to provide common functionality
type BaseWidget struct {
	Rect    q2d.Rectangle
	Tooltip string
	Theme   *Theme
	Fill    bool
}

func (b *BaseWidget) SetRect(r q2d.Rectangle) {
	b.Rect = r
}

func (b *BaseWidget) GetRect() q2d.Rectangle {
	return b.Rect
}

func (b *BaseWidget) Layout(available Size) Size {
	return available
}

func (b *BaseWidget) Draw(img *q2d.Image) {
	// No-op
}

func (b *BaseWidget) Event(e Event) bool {
	return false
}

func (b *BaseWidget) MinSize() Size {
	return Size{0, 0}
}

func (b *BaseWidget) GetTooltip() string {
	return b.Tooltip
}

func (b *BaseWidget) GetTheme() *Theme {
	if b.Theme != nil {
		return b.Theme
	}
	return DefaultTheme
}

func (b *BaseWidget) SetTheme(t *Theme) {
	b.Theme = t
}

func (b *BaseWidget) FindWidgetAt(pos q2d.Point) Widget {
	if b.Rect.Contains(pos) {
		return b
	}
	return nil
}

func (b *BaseWidget) IsFill() bool {
	return b.Fill
}
